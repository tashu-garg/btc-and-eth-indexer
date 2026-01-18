package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"indexer/internal/model"
	"indexer/internal/repository"
	"log"
	"net/http"
	"time"
)

type BTCWorker struct {
	repo         repository.Repository
	rpcURL       string
	user         string
	pass         string
	startHeight  int
	syncInterval time.Duration
	client       *http.Client
}

func NewBTCWorker(repo repository.Repository, rpcURL, apiKey string, startHeight, syncIntervalMS int) *BTCWorker {
	return &BTCWorker{
		repo:         repo,
		rpcURL:       rpcURL,
		user:         "", // Could be configured if needed
		pass:         apiKey,
		startHeight:  startHeight,
		syncInterval: time.Duration(syncIntervalMS) * time.Millisecond,
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

type jsonRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type jsonRPCResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	ID int `json:"id"`
}

func (w *BTCWorker) callRPC(method string, params []interface{}) (json.RawMessage, error) {
	reqBody, _ := json.Marshal(jsonRPCRequest{
		JSONRPC: "1.0",
		Method:  method,
		Params:  params,
		ID:      1,
	})

	req, err := http.NewRequest("POST", w.rpcURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if w.user != "" || w.pass != "" {
		req.SetBasicAuth(w.user, w.pass)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return nil, fmt.Errorf("RPC returned status %d: %s", resp.StatusCode, buf.String())
	}

	var rpcResp jsonRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, fmt.Errorf("failed to decode RPC response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error: %d %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}

func (w *BTCWorker) Start(ctx context.Context) {
	log.Println("[BTC] Worker starting...")

	// 1. Validate RPC connection
	tip, err := w.getTip()
	if err != nil {
		log.Printf("[BTC] RPC connection failed: %v. Indexer will remain idle.", err)
		// We don't return here so the loop starts and can retry later,
		// or we can just return if we want this worker specifically to stop but keep the app alive.
		return
	}
	log.Printf("[BTC] RPC connection validated. Current tip: %d", tip)

	// 2. Ensure state exists
	_, err = w.repo.GetOrCreateState(model.ChainBTC, tip, w.startHeight)
	if err != nil {
		log.Fatalf("[BTC] Failed to initialize state: %v", err)
	}

	ticker := time.NewTicker(w.syncInterval)
	defer ticker.Stop()

	log.Println("[BTC] Worker sync loop started")
	for {
		select {
		case <-ctx.Done():
			log.Println("[BTC] Worker stopping...")
			return
		case <-ticker.C:
			if err := w.sync(ctx); err != nil {
				log.Printf("[BTC] Sync error: %v", err)
			}
		}
	}
}

func (w *BTCWorker) sync(ctx context.Context) error {
	tip, err := w.getTip()
	if err != nil {
		return fmt.Errorf("failed to get tip: %w", err)
	}

	lastIndexed, err := w.repo.GetOrCreateState(model.ChainBTC, tip, w.startHeight)
	if err != nil {
		return fmt.Errorf("failed to get state: %w", err)
	}

	if lastIndexed >= tip {
		return nil
	}

	nextHeight := lastIndexed + 1

	log.Printf("[BTC] Syncing block %d / %d", nextHeight, tip)

	block, txs, err := w.fetchBlock(nextHeight)
	if err != nil {
		return fmt.Errorf("failed to fetch block %d: %w", nextHeight, err)
	}

	err = w.repo.SaveBlockWithTransactions(block, txs)
	if err != nil {
		return fmt.Errorf("failed to save block %d: %w", nextHeight, err)
	}

	return nil
}

func (w *BTCWorker) getTip() (uint64, error) {
	res, err := w.callRPC("getblockcount", nil)
	if err != nil {
		return 0, err
	}
	var count uint64
	if err := json.Unmarshal(res, &count); err != nil {
		return 0, err
	}
	return count, nil
}

func (w *BTCWorker) getTx(txid string) ([]byte, error) {
	return w.callRPC("getrawtransaction", []interface{}{txid, true})
}

func (w *BTCWorker) fetchBlock(height uint64) (*model.Block, []*model.Transaction, error) {
	// 1. Get block hash
	resHash, err := w.callRPC("getblockhash", []interface{}{height})
	if err != nil {
		return nil, nil, err
	}
	var hash string
	if err := json.Unmarshal(resHash, &hash); err != nil {
		return nil, nil, err
	}

	// 2. Get block details (verbosity 2 for full tx details)
	resBlock, err := w.callRPC("getblock", []interface{}{hash, 2})
	if err != nil {
		return nil, nil, err
	}

	var rpcBlock struct {
		Hash              string `json:"hash"`
		Height            uint64 `json:"height"`
		Time              int64  `json:"time"`
		PreviousBlockHash string `json:"previousblockhash"`
		Tx                []struct {
			Txid string `json:"txid"`
			Vin  []struct {
				Txid string `json:"txid"`
				Vout int    `json:"vout"`
			} `json:"vin"`
			Vout []struct {
				Value        float64 `json:"value"`
				ScriptPubKey struct {
					Address   string   `json:"address"`
					Addresses []string `json:"addresses"`
				} `json:"scriptPubKey"`
			} `json:"vout"`
		} `json:"tx"`
	}

	if err := json.Unmarshal(resBlock, &rpcBlock); err != nil {
		return nil, nil, err
	}

	block := &model.Block{
		Chain:     model.ChainBTC,
		Height:    rpcBlock.Height,
		Hash:      rpcBlock.Hash,
		BlockHash: rpcBlock.PreviousBlockHash,
		TXCount:   uint64(len(rpcBlock.Tx)),
		Timestamp: time.Unix(rpcBlock.Time, 0),
	}

	var txs []*model.Transaction
	for _, rt := range rpcBlock.Tx {
		// ---------------- FROM ADDRESSES ----------------
		fromSet := map[string]bool{}

		isCoinbase := false
		if len(rt.Vin) > 0 && rt.Vin[0].Txid == "" {
			isCoinbase = true
			fromSet["coinbase"] = true
		}

		if !isCoinbase {
			for _, vin := range rt.Vin {
				resPrev, err := w.getTx(vin.Txid)
				if err != nil {
					// This often fails if txindex=1 is not set on the node
					continue
				}

				var prevTx struct {
					Vout []struct {
						ScriptPubKey struct {
							Address   string   `json:"address"`
							Addresses []string `json:"addresses"`
						} `json:"scriptPubKey"`
					} `json:"vout"`
				}

				if err := json.Unmarshal(resPrev, &prevTx); err != nil {
					continue
				}

				if vin.Vout < len(prevTx.Vout) {
					spk := prevTx.Vout[vin.Vout].ScriptPubKey
					if spk.Address != "" {
						fromSet[spk.Address] = true
					} else {
						for _, addr := range spk.Addresses {
							fromSet[addr] = true
						}
					}
				}
			}
		}

		var fromList []string
		for a := range fromSet {
			fromList = append(fromList, a)
		}
		from := "unknown"
		if len(fromList) > 0 {
			// Use the first address as primary and indicate if there are more
			from = fromList[0]
			if len(fromList) > 1 {
				from = fmt.Sprintf("%s,+%d others", fromList[0], len(fromList)-1)
			}
		} else if isCoinbase {
			from = "coinbase"
		}

		// ---------------- TO ADDRESSES ----------------
		toSet := map[string]bool{}
		value := 0.0

		for _, v := range rt.Vout {
			value += v.Value
			if v.ScriptPubKey.Address != "" {
				toSet[v.ScriptPubKey.Address] = true
			} else {
				for _, addr := range v.ScriptPubKey.Addresses {
					toSet[addr] = true
				}
			}
		}

		var toList []string
		for a := range toSet {
			toList = append(toList, a)
		}
		to := "unknown"
		if len(toList) > 0 {
			to = toList[0]
			if len(toList) > 1 {
				to = fmt.Sprintf("%s,+%d others", toList[0], len(toList)-1)
			}
		} else if len(rt.Vout) > 0 {
			to = "non-standard"
		}

		txs = append(txs, &model.Transaction{
			Chain:     model.ChainBTC,
			Hash:      rt.Txid,
			BlockHash: rpcBlock.Hash,
			Height:    rpcBlock.Height,
			From:      from,
			To:        to,
			Value:     fmt.Sprintf("%.8f", value),
			Status:    "success",
			Timestamp: block.Timestamp,
		})
	}

	return block, txs, nil
}
