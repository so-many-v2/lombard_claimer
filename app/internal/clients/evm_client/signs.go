package evmclient

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func (ec *EVMClient) SignMessage(msg string) (string, error) {

	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msg))
	prefixedMsg := []byte(prefix + msg)

	hash := crypto.Keccak256Hash(prefixedMsg)

	signature, err := crypto.Sign(hash.Bytes(), ec.Wallet.PrivateKey)
	if err != nil {
		return "", err
	}

	signature[64] += 27

	return hexutil.Encode(signature), nil
}

func (ec *EVMClient) SendTransaction(toAddress common.Address, txData []byte) (common.Hash, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(ec.Config.SendTxTimeout))
	defer cancel()

	nonce, err := ec.Client.NonceAt(ctx, ec.Wallet.Address, nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error getting nonce: %s", err.Error())
	}

	gasPrice, err := ec.Client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error getting SuggestGasPrice: %s", err.Error())
	}

	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(ec.Config.GasMultiplier))
	gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))

	msg := ethereum.CallMsg{
		From: ec.Wallet.Address,
		To:   &toAddress,
		Data: txData,
	}

	gasLimit, err := ec.Client.EstimateGas(ctx, msg)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error getting gasLimit: %s", err.Error())
	}
	gasLimit = gasLimit * uint64(ec.Config.GasMultiplier) / 100

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
