package evmclient

import (
	"fmt"
	"lombardClaimer/internal/config"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EVMClient struct {
	Rpc       string
	ChainName string
	Client    *ethclient.Client
	Wallet    *Wallet
	Config    config.EvmClientConfig
}

func NewEVMClient(chainName, rpc, pk string, config config.EvmClientConfig) (*EVMClient, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, fmt.Errorf("error connect rpc: %s", err.Error())
	}

	wallet, err := NewWallet(pk)
	if err != nil {
		return nil, fmt.Errorf("error init wallet from private key: %s", err.Error())
	}

	return &EVMClient{
		Rpc:       rpc,
		ChainName: chainName,
		Client:    client,
		Wallet:    wallet,
		Config:    config,
	}, nil
}

func (ec *EVMClient) GetAbi(abiStr string) (abi.ABI, error) {
	contractAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("error read abi: %s", err.Error())
	}
	return contractAbi, nil
}
