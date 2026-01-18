package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPass            string
	DBName            string
	BTCRPC            string
	BTCKey            string
	BTCStartHeight    int
	BTCSyncIntervalMS int
	ETHRPC            string
	ETHStartHeight    int
	ETHSyncIntervalMS int
	ServerPort        string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBUser:            os.Getenv("DB_USER"),
		DBPass:            os.Getenv("DB_PASSWORD"),
		DBName:            os.Getenv("DB_NAME"),
		BTCRPC:            os.Getenv("BTC_RPC_URL"),
		BTCKey:            os.Getenv("BTC_RPC_PASS"),
		BTCStartHeight:    getEnvInt("BTC_START_HEIGHT", 0),
		BTCSyncIntervalMS: getEnvInt("BTC_SYNC_INTERVAL_MS", 2000),
		ETHRPC:            os.Getenv("ETH_RPC_URL"),
		ETHStartHeight:    getEnvInt("ETH_START_HEIGHT", 0),
		ETHSyncIntervalMS: getEnvInt("ETH_SYNC_INTERVAL_MS", 2000),
		ServerPort:        os.Getenv("PORT"),
	}
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	var res int
	fmt.Sscanf(val, "%d", &res)
	return res
}
