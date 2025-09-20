package app

import (
	"fmt"
	"log"
	evmClient "lombardClaimer/internal/clients/evm_client"
	httpClient "lombardClaimer/internal/clients/http_client"
	"lombardClaimer/internal/config"
	"lombardClaimer/internal/modules/lombard"
	fr "lombardClaimer/pkg/fileReader"
	"sync"
)

var LOMBARD_CLAIM_ADDRESS string = "0x6fF742845D45d29cb38fa075EFc889247A52Eb02"

var LOMBARD_VAULT_ADDRESSES = map[string]string{
	"Base":     "0x3b4acd8879fb60586ccd74bc2f831a4c5e7dbbf8",
	"Ethereum": "",
	"Bcs":      "",
}

func WithdrawLombard() {

	conf := config.NewConfig()
	reader := fr.NewFileReader()

	limiter := make(chan struct{}, conf.General.WorkerPool)
	wg := sync.WaitGroup{}

	wallets, err := reader.ScanFile("data/wallets.txt")
	if err != nil {
		log.Fatal(err.Error())
	}

	provider, err := config.GetProvider(conf.ClaimLombard.ChainName)
	if err != nil {
		log.Fatal(err.Error())
	}

	tokens, err := config.GetTokens("Base")
	if err != nil {
		log.Fatal(err.Error())
	}

	httpClient, err := httpClient.NewFetcher(conf.General.ResidentialProxy)

	for _, key := range wallets {

		wg.Add(1)
		limiter <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-limiter }()

			evmClient, err := evmClient.NewEVMClient(
				conf.ClaimLombard.ChainName,
				provider,
				key,
				conf.EvmClient,
			)

			if err != nil {
				fmt.Println("error init evm client from pk")
				return
			}

			lombard := lombard.NewLombard(
				evmClient,
				httpClient,
				LOMBARD_VAULT_ADDRESSES[evmClient.ChainName],
				tokens["LBTCv"],
			)

			err = lombard.Withdraw()
			if err != nil {
				fmt.Println(err.Error())
			}

		}()

	}

	wg.Wait()
	fmt.Printf("Soft Ends Work\n")

}
