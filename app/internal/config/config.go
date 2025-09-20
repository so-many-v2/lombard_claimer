package config

type GeneralConfig struct {
	ResidentialProxy string
	WorkerPool       int64
}

type EvmClientConfig struct {
	GasMultiplier int64
	SendTxTimeout int64
}

type ClaimLombardConfig struct {
	ChainName string
}

type Config struct {
	General      GeneralConfig
	ClaimLombard ClaimLombardConfig
	EvmClient    EvmClientConfig
}

func NewConfig() *Config {
	return &Config{
		General: GeneralConfig{
			ResidentialProxy: "", // Your resedental proxy
			WorkerPool:       5,
		},
		ClaimLombard: ClaimLombardConfig{
			ChainName: "Ethereum",
		},
		EvmClient: EvmClientConfig{
			GasMultiplier: 115,
			SendTxTimeout: 30,
		},
	}
}
