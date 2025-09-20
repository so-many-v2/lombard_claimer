package config

type CommonConfig struct {
	ResidentalProxy string
}

type ClaimLombardConfig struct {
	ChainName string
}

type Config struct {
	Common       CommonConfig
	ClaimLombard ClaimLombardConfig
}

func NewConfig() *Config {
	return &Config{
		Common: CommonConfig{
			ResidentalProxy: "", // Your resedental proxy
		},
		ClaimLombard: ClaimLombardConfig{
			ChainName: "Ethereum",
		},
	}
}
