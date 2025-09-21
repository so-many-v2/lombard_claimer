package app

import (
	"fmt"
	"log"
	evmClient "lombardClaimer/internal/clients/evm_client"
	httpClient "lombardClaimer/internal/clients/http_client"
	"lombardClaimer/internal/config"
	contractabis "lombardClaimer/internal/config/contract_abis"
	"lombardClaimer/internal/modules/lombard"
	fr "lombardClaimer/pkg/fileReader"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

var LOMBARD_CLAIM_ADDRESS string = "0x6fF742845D45d29cb38fa075EFc889247A52Eb02"

func ClaimLombard() {

	conf := config.NewConfig()
	reader := fr.NewFileReader()

	limiter := make(chan struct{}, conf.General.WorkerPool)
	wg := sync.WaitGroup{}

	wallets, err := reader.ScanFile("data/wallets.txt")
	if err != nil {
		log.Fatal(err.Error())
	}

	depositAddresses := []string{}

	if conf.ClaimLombard.SendTokensAfterClaim {

		depositAddresses, err = reader.ScanFile("data/deposit_addresses.txt")
		if err != nil {
			log.Fatal(err.Error())
		}

		if len(wallets) != len(depositAddresses) {
			log.Fatalf("amount wallets and deposit address are different\n")
		}

	}

	provider, err := config.GetProvider("Ethereum")
	if err != nil {
		log.Fatal(err.Error())
	}

	tokens, err := config.GetTokens("Ethereum")
	if err != nil {
		log.Fatal(err.Error())
	}

	httpClient, err := httpClient.NewFetcher(conf.General.ResidentialProxy)

	for i, key := range wallets {

		wg.Add(1)
		limiter <- struct{}{}

		go func(idx int, key string) {
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

			err = lombard.MakeClaim(LOMBARD_CLAIM_ADDRESS)
			if err != nil {
				fmt.Println(err.Error())
			}

			if conf.ClaimLombard.SendTokensAfterClaim {

				receiverAddress := depositAddresses[idx]
				bardTokenAddress := tokens["BARD"]

				bardTokenAmount, err := evmClient.CallContract(
					bardTokenAddress,
					contractabis.ERC20_ABI,
					"balanceOf",
					evmClient.Wallet.Address,
				)

				bardTokenAmountInt, ok := bardTokenAmount.(*big.Int)
				if !ok {
					fmt.Printf("unexpected balanceOf result type: %T\n", bardTokenAmount)
					return
				}

				if err != nil {
					fmt.Println(err.Error())
					return
				}

				tokenAbi, err := evmClient.GetAbi(contractabis.ERC20_ABI)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				txData, err := tokenAbi.Pack(
					"transfer",
					common.HexToAddress(receiverAddress),
					bardTokenAmountInt,
				)

				txHash, err := evmClient.SendTransaction(
					common.HexToAddress(bardTokenAddress),
					txData,
				)

				fmt.Printf(
					"Wallet: %s Send Bard tokens to %s | Hash: %s\n\n",
					evmClient.Wallet.Address.Hex(),
					receiverAddress,
					txHash.Hex(),
				)

			}

		}(i, key)

	}

	wg.Wait()
	fmt.Printf("Soft Ends Work\n")

}
