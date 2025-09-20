package config

import "fmt"

func GetProvider(chanName string) (string, error) {

	providers := map[string]string{
		"Ethereum": "https://eth.llamarpc.com",
		"Base":     "https://rpc.ankr.com/base/a08b52a16806752b7cce3d78aec157bd12c176474abea1e32d2683b74ed57638",
		"Bsc":      "https://base.llamarpc.com",
	}

	provider, ok := providers[chanName]
	if !ok {
		return "", fmt.Errorf("provider not found: %s", chanName)
	}

	return provider, nil
}
