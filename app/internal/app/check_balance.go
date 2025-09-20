package app

import (
	"fmt"
	"log"
	evmClient "lombardClaimer/internal/clients/evm_client"
	"lombardClaimer/internal/config"
	contractabis "lombardClaimer/internal/config/contract_abis"
	fr "lombardClaimer/pkg/fileReader"
	"math/big"
	"sync"
)

func CheckBard() {

	conf := config.NewConfig()
	reader := fr.NewFileReader()

	limiter := make(chan struct{}, conf.General.WorkerPool)
	wg := sync.WaitGroup{}

	wallets, err := reader.ScanFile("data/wallets.txt")
	if err != nil {
		log.Fatal(err.Error())
	}

	provider, err := config.GetProvider("Ethereum")
	if err != nil {
		log.Fatal(err.Error())
	}

	tokens, err := config.GetTokens("Ethereum")
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, key := range wallets {

		wg.Add(1)
		limiter <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-limiter }()

			evmClient, err := evmClient.NewEVMClient(
				"Ethereum",
				provider,
				key,
				conf.EvmClient,
			)

			if err != nil {
				fmt.Printf("error init evm client from pk\n")
				return
			}

			balance, err := evmClient.CallContract(
				tokens["BARD"],
				contractabis.ERC20_ABI,
				"balanceOf",
				evmClient.Wallet.Address,
			)

			bigBalance, ok := balance.(*big.Int)
			if !ok {
				fmt.Printf("cant parse balance to *Big.Int\n")
				return
			}

			if bigBalance.Cmp(big.NewInt(0)) == 1 {
				fmt.Printf("Wallet: %s | Have BARD balance %v\n", evmClient.Wallet.Address.Hex(), balance)
			}

		}()

	}

	wg.Wait()
	fmt.Printf("Soft Ends Work\n")

}
