package evmclient

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func (ec *EVMClient) Approve(tokenAddress, spenderAddress, abiStr string, amount int64) error {

	tokenAddressHex := common.HexToAddress(tokenAddress)
	spenderAddressHex := common.HexToAddress(spenderAddress)

	decimalsResp, err := ec.CallContract(
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
