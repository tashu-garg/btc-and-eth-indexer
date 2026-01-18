package main

import (
	"context"
	"indexer/internal/config"
	"indexer/internal/db"
	"indexer/internal/handlers"
	"indexer/internal/repository"
	"indexer/internal/routes"
	"indexer/internal/workers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := config.LoadConfig()
	database := db.InitDB(cfg)
	repo := repository.NewRepository(database)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Initializing Workers
	btcWorker := workers.NewBTCWorker(repo, cfg.BTCRPC, cfg.BTCKey, cfg.BTCStartHeight, cfg.BTCSyncIntervalMS)
	ethWorker, err := workers.NewETHWorker(repo, cfg.ETHRPC, cfg.ETHStartHeight, cfg.ETHSyncIntervalMS)
	if err != nil {
		log.Printf("[MAIN] ETH Worker initialization warning: %v", err)
	}

	// 2. Starting Workers in Goroutines
	if btcWorker != nil {
		go btcWorker.Start(ctx)
		log.Println("[MAIN] BTC sync worker spawned")
	}
	if ethWorker != nil {
		go ethWorker.Start(ctx)
		log.Println("[MAIN] ETH sync worker spawned")
	}

	// 3. API Handlers Layer
	apiHandler := handlers.NewAPIHandler(repo)

	// 4. Router Setup
	r := routes.SetupRouter(apiHandler)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Printf("[MAIN] API Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[MAIN] Server failed: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[MAIN] Shutting down gracefully...")
	cancel() // Stop workers

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("[MAIN] Server forced shutdown:", err)
	}

	log.Println("[MAIN] Good bye!")
}
