package config

type Config struct {
	ETH  ETHConfig
	BTC  BTCConfig
	DB   DBConfig
	HTTP HTTPConfig
}

func NewConfig() *Config {
	return &Config{
		ETH:  LoadETHConfig(),
		BTC:  LoadBTCConfig(),
		DB:   LoadDBConfig(),
		HTTP: LoadHTTPConfig(),
	}
}
