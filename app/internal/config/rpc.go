package config

import "fmt"

func GetProvider(chanName string) (string, error) {

	providers := map[string]string{
		"Ethereum": "https://eth.llamarpc.com",
		"Base":     "https://base-rpc.publicnode.com",
		"Bsc":      "https://base.llamarpc.com",
	}

	provider, ok := providers[chanName]
	if !ok {
		return "", fmt.Errorf("provider not found: %s", chanName)
	}

	return provider, nil
}
