package app

import (
	"context"
	"fmt"
	"log"
	evmclient "lombardClaimer/internal/clients/evm_client"
	httpclient "lombardClaimer/internal/clients/http_client"
	"lombardClaimer/internal/config"
	abis "lombardClaimer/internal/config/contract_abis"
	"lombardClaimer/internal/modules"
	fr "lombardClaimer/pkg/fileReader"

	"github.com/ethereum/go-ethereum/common"
)

var LOMBARD_CLAIM_ADDRESS string = "0x6fF742845D45d29cb38fa075EFc889247A52Eb02"

func ClaimLombard() {

	conf := config.NewConfig()
	reader := fr.NewFileReader()

	wallets := reader.ScanFile("data/wallets.txt")

	provider, err := config.GetProvider(conf.ClaimLombard.ChainName)
	if err != nil {
		log.Fatal(err.Error())
	}

	// tokens, err := config.GetTokens(conf.ClaimLombard.ChainName)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	httpClient, err := httpclient.NewFetcher(conf.Common.ResidentalProxy)

	for wallets.Scan() {

		evmClient, err := evmclient.NewEVMClient(
			conf.ClaimLombard.ChainName,
			provider,
			wallets.Text(),
		)

		if err != nil {
			fmt.Println("error init evm client from pk")
			continue
		}

		lombard := modules.NewLombard(evmClient, httpClient)

		err := lombard.ClaimAllocation()

	}

}
