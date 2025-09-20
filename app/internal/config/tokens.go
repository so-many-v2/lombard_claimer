package config

import "fmt"

func GetTokens(chainName string) (map[string]string, error) {
	data := map[string]map[string]string{
		"Base": map[string]string{
			"USDC":  "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
			"BARD":  "",
			"LBTCv": "0x5401b8620E5FB570064CA9114fd1e135fd77D57c",
		},
		"Ethereum": map[string]string{
			"USDC":  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
			"USDT":  "0xdac17f958d2ee523a2206206994597c13d831ec7",
			"BARD":  "0xf0DB65D17e30a966C2ae6A21f6BBA71cea6e9754",
			"LBTCv": "0x5401b8620E5FB570064CA9114fd1e135fd77D57c",
			"LBTC":  "0xecAc9C5F704e954931349Da37F60E39f515c11c1",
		},
		"BSC": map[string]string{
			"USDC": "0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d",
			"USDT": "0x55d398326f99059ff775485246999027b3197955",
			"BARD": "",
		},
		"Arbitrum": map[string]string{
			"USDC": "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
			"USDT": "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
		},
		"Optimism": map[string]string{
			"USDC": "0x0b2C639c533813f4Aa9D7837CAf62653d097Ff85",
			"USDT": "0x94b008aa00579c1307b0ef2c499ad98a8ce58e58",
		},
	}

	tokens, ok := data[chainName]
	if !ok {
		return nil, fmt.Errorf("incorrect chain name: %s", chainName)
	}

	return tokens, nil
}
