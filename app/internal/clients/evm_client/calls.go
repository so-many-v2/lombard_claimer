package evmclient

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func (ec *EVMClient) CallContract(address, abiStr string, methodName string, args ...interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(ec.Config.SendTxTimeout))
	defer cancel()

	addr := common.HexToAddress(address)

	contractAbi, err := ec.GetAbi(abiStr)
	if err != nil {
		return nil, err
	}

	data, err := contractAbi.Pack(methodName, args...)
	if err != nil {
		return nil, fmt.Errorf("error pack abi method %s: %s", methodName, err.Error())
	}

	msg := ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}

	resp, err := ec.Client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting response from rpc: %s", err.Error())
	}

	var out interface{}

	if err := contractAbi.UnpackIntoInterface(&out, methodName, resp); err != nil {
		return nil, fmt.Errorf("error uppack response data: %s", err.Error())
	}

	return out, nil
}

func (ec *EVMClient) GetDecimals(tokenAddress, tokenAbi string) (int64, error) {

	response, err := ec.CallContract(
		tokenAddress,
		tokenAbi,
		"decimals",
	)

	if err != nil {
		return 0, err
	}

	decimals, ok := response.(int64)
	if !ok {
		return 0, fmt.Errorf("can't parse decimals to int64 | response %v", response)
	}

	return decimals, nil
}

func (ec *EVMClient) GetBalanceOf(tokenAddress, tokenAbi string) (*big.Int, error) {

	response, err := ec.CallContract(
		tokenAddress,
		tokenAbi,
		"balanceOf",
		ec.Wallet.Address,
	)

	if err != nil {
		return nil, err
	}

	decimals, ok := response.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("can't parse decimals to *Big.Int | response %v", response)
	}

	return decimals, nil
}

func (ec *EVMClient) GetNativeBalance() (*big.Int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(ec.Config.SendTxTimeout))
	defer cancel()

	balance, err := ec.Client.BalanceAt(ctx, ec.Wallet.Address, nil)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
