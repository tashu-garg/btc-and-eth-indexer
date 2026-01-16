package config

import (
	"log"
	"os"
	"strconv"
)

type DBConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Name           string
	Driver         string
	SSLMode        string
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBConnMaxLife  int
	AppEnv         string
}

func getenv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:           getenv("DB_HOST", "localhost"),
		Port:           getenv("DB_PORT", "5435"),
		User:           getenv("DB_USER", "admin"),
		Password:       getenv("DB_PASSWORD", "admin"),
		Name:           getenv("DB_NAME", "indexer_db"),
		Driver:         getenv("DB_DRIVER", "postgres"),
		SSLMode:        getenv("DB_SSL", "disable"),
		DBMaxOpenConns: getenvInt("DB_MAX_OPEN_CONNS", 10),
		DBMaxIdleConns: getenvInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLife:  getenvInt("DB_CONN_MAX_LIFE", 360),
		AppEnv:         getenv("GIN_MODE", "debug"),
	}
}

func getenvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Fatalf("Invalid %s: %v", key, err)
	}
	return val
}
