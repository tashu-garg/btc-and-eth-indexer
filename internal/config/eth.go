package config

import (
	"os"
	"time"
)

type ETHConfig struct {
	RPCURL            string
	RateLimitInterval time.Duration
}

func LoadETHConfig() ETHConfig {
	return ETHConfig{
		RPCURL:            os.Getenv("ETH_RPC_URL"),
		RateLimitInterval: 100 * time.Millisecond,
	}
}
