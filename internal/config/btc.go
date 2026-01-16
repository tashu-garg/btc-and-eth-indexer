package config

import (
	"os"
	"time"
)

type BTCConfig struct {
	RPCURL            string
	RPCUser           string
	RPCPass           string
	Network           string // mainnet | testnet
	RateLimitInterval time.Duration
}

func LoadBTCConfig() BTCConfig {
	return BTCConfig{
		RPCURL:            os.Getenv("BTC_RPC_URL"),
		RPCUser:           os.Getenv("BTC_RPC_USER"),
		RPCPass:           os.Getenv("BTC_RPC_PASS"),
		Network:           getEnv("BTC_NETWORK", "mainnet"),
		RateLimitInterval: 100 * time.Millisecond,
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
