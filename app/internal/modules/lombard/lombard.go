package lombard

import (
	evmclient "lombardClaimer/internal/clients/evm_client"
	httpclient "lombardClaimer/internal/clients/http_client"
)

type Lombard struct {
	evmClient    *evmclient.EVMClient
	httpClient   *httpclient.HttpClient
	VaultAddress string
	LBTCvAddress string
}

func NewLombard(evmClient *evmclient.EVMClient, httpClient *httpclient.HttpClient, va, la string) *Lombard {
	return &Lombard{
		evmClient:    evmClient,
		httpClient:   httpClient,
		VaultAddress: va,
		LBTCvAddress: la,
	}
}
