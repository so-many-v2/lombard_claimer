package app

import (
	"fmt"
	"log"
	evmclient "lombardClaimer/internal/clients/evm_client"
	"lombardClaimer/internal/config"
	abis "lombardClaimer/internal/config/contract_abis"
	fr "lombardClaimer/pkg/fileReader"

	"github.com/ethereum/go-ethereum/common"
)

func TestScript() {

	conf := config.NewConfig()
	reader := fr.NewFileReader()

	wallets, err := reader.ScanFile("data/wallets.txt")
	if err != nil {
		log.Fatal(err.Error())
	}

	provider, err := config.GetProvider(conf.ClaimLombard.ChainName)
	if err != nil {
		log.Fatal(err.Error())
	}

	tokens, err := config.GetTokens(conf.ClaimLombard.ChainName)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, key := range wallets {

		evmClient, err := evmclient.NewEVMClient(
			conf.ClaimLombard.ChainName,
			provider,
			key,
			config.NewConfig().EvmClient,
		)

		if err != nil {
			fmt.Println("error init evm client from pk")
			continue
		}

		response, err := evmClient.CallContract(
			tokens["USDC"],
			abis.ERC20_ABI,
			"balanceOf",
			common.HexToAddress("0x11a0f7449c4E1cB3200bC63Ce4E75f4859205230"),
		)

		if err != nil {
			fmt.Printf("error gating balance USDC: %v | %s\n", evmClient.Wallet.Address, err.Error())
			continue
		}

		fmt.Printf("USDC Balance: %v\n", response)
	}

}
