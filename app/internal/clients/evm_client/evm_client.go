package evmclient

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EVMClient struct {
	Rpc       string
	ChainName string
	Client    *ethclient.Client
	Wallet    *Wallet
}

func NewEVMClient(chainName, rpc, pk string) (*EVMClient, error) {
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
	}, nil
}

func (ec *EVMClient) GetAbi(abiStr string) (abi.ABI, error) {
	contractAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("error read abi: %s", err.Error())
	}
	return contractAbi, nil
}

func (ec *EVMClient) SendTransaction(toAddress common.Address, txData []byte) (common.Hash, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	nonce, err := ec.Client.NonceAt(ctx, ec.Wallet.Address, nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error getting nonce: %s", err.Error())
	}

	gasPrice, err := ec.Client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error getting SuggestGasPrice: %s", err.Error())
	}

	msg := ethereum.CallMsg{
		From: ec.Wallet.Address,
		To:   &toAddress,
		Data: txData,
	}

	gasLimit, err := ec.Client.EstimateGas(ctx, msg)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error getting gasLimit: %s", err.Error())
	}

	chainId, err := ec.Client.ChainID(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error getting chainID: %s", err.Error())
	}

	tx := types.NewTransaction(nonce, toAddress, big.NewInt(0), gasLimit, gasPrice, txData)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), ec.Wallet.PrivateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error getting signedTx: %s", err.Error())
	}

	if err := ec.Client.SendTransaction(ctx, signedTx); err != nil {
		return common.Hash{}, fmt.Errorf("error sending transaction: %s", err.Error())
	}

	return signedTx.Hash(), nil

}

func (ec *EVMClient) GetBalance() (*big.Int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	balance, err := ec.Client.BalanceAt(ctx, ec.Wallet.Address, nil)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (ec *EVMClient) CallContract(ctx context.Context, address, abiStr string, methodName string, args ...interface{}) (any, error) {

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

func (ec *EVMClient) Approve(tokenAddress, spenderAddress, abiStr string, amount int64) error {

	tokenAddressHex := common.HexToAddress(tokenAddress)
	spenderAddressHex := common.HexToAddress(spenderAddress)

	decimalsResp, err := ec.CallContract(
		context.Background(),
		tokenAddress,
		abiStr,
		"decimals",
	)

	if err != nil {
		return fmt.Errorf("error getting decimals: %s", err.Error())
	}

	intDecimals, ok := decimalsResp.([]interface{})[0].(uint8)
	if !ok {
		return fmt.Errorf("error get int64 from decimal interface: %s", fmt.Errorf("error type decimal switch"))
	}

	regularAmount := big.NewInt(amount)
	weiMul := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(intDecimals)), nil)
	weiAmount := new(big.Int).Mul(regularAmount, weiMul)

	contractAbi, err := ec.GetAbi(abiStr)
	if err != nil {
		return err
	}

	txData, err := contractAbi.Pack("approve", spenderAddressHex, weiAmount)
	if err != nil {
		return fmt.Errorf("error getting txData: %s", err.Error())
	}

	txHash, err := ec.SendTransaction(tokenAddressHex, txData)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	fmt.Printf("Make approve | Tx Hash: %sending", txHash.String())
	return nil
}
