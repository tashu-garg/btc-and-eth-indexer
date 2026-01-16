package config

import "os"

type HTTPConfig struct {
	Host string
	Port string
	Mode string
}

func LoadHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
		Mode: os.Getenv("GIN_MODE"),
	}
}
