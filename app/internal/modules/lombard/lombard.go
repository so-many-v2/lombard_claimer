package lombard

import (
	evmclient "lombardClaimer/internal/clients/evm_client"
	httpclient "lombardClaimer/internal/clients/http_client"
)

type Lombard struct {
	evmClient        *evmclient.EVMClient
	httpClient       *httpclient.HttpClient
	BaseVaultAddress string
	BscVaultAddress  string
}

func NewLombard(evmClient *evmclient.EVMClient, httpClient *httpclient.HttpClient) *Lombard {
	return &Lombard{
		evmClient:  evmClient,
		httpClient: httpClient,
	}
}
