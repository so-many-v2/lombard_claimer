package modules

import (
	evmclient "lombardClaimer/internal/clients/evm_client"
	httpclient "lombardClaimer/internal/clients/http_client"

	"github.com/ethereum/go-ethereum/common"
)

type Lombard struct {
	evmClient  *evmclient.EVMClient
	httpClient *httpclient.HttpClient
}

func NewLombard(evmClient *evmclient.EVMClient, httpClient *httpclient.HttpClient) *Lombard {
	return &Lombard{
		evmClient:  evmClient,
		httpClient: httpClient,
	}
}

func (l *Lombard) ClaimAllocation(contractAddress, contractAbi string) error {

	caHex := common.HexToAddress(contractAddress)

}
