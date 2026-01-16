package main

import (
	application "indexer"
	"indexer/internal/config"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	cfg := config.NewConfig()
	application.Start(cfg)
}
