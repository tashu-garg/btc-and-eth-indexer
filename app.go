package application

import (
	"context"
	"indexer/internal/btc"
	"indexer/internal/config"
	"indexer/internal/eth"
	"indexer/internal/infrastructure/db"
	"indexer/internal/middleware"
	"indexer/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"indexer/internal/cron"
)

func Start(cfg *config.Config) {
	middleware.InitLogger()
	log.Println("Starting Native Sync Service...")

	// 1. Init DB
	database := db.InitDB(cfg.DB)

	// 2. Init ETH Service
	ethService, err := eth.NewService(database, cfg.ETH)
	if err != nil {
		log.Fatalf("Failed to create ETH service: %v", err)
	}

	// 3. Init BTC Service
	btcService, err := btc.NewService(database, cfg.BTC)
	if err != nil {
		log.Fatalf("Failed to create BTC service: %v", err)
	}

	// 5. Init Handlers
	ethHandler := eth.NewHandler(ethService)
	btcHandler := btc.NewHandler(btcService)

	// 6. Init Server
	srv := server.NewServer(cfg.HTTP)
	server.RegisterRoutes(srv.GetRouter(), ethHandler, btcHandler)

	// 7. Context for services
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 8. Start ETH Service
	go ethService.Run(ctx)

	// 9. Start BTC Service
	// go btcService.Start(ctx)

	// 10. Setup Custom Scheduler
	scheduler := cron.NewScheduler()
	defer scheduler.Stop()

	// Schedule ETH sync every 200ms
	scheduler.AddJob(200*time.Millisecond, func() {
		if err := ethService.ProcessNextBlock(); err != nil {
			log.Printf("ETH Process Error: %v", err)
		}
	})

	// Schedule BTC sync every 200ms
	scheduler.AddJob(200*time.Millisecond, func() {
		if err := btcService.ProcessNextBlock(ctx); err != nil {
			log.Printf("[DEBUG_APP] BTC Process Error: %v", err)
		}
	})

	// 11. Start HTTP Server
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 12. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
