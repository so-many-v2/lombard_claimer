package evmclient

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func (ec *EVMClient) Approve(tokenAddress, spenderAddress, abiStr string, amount *big.Int) error {

	taAddress := common.HexToAddress(tokenAddress)
	sa := common.HexToAddress(spenderAddress)

	contractAbi, err := ec.GetAbi(abiStr)
	if err != nil {
		return err
	}

	txData, err := contractAbi.Pack("approve", sa, amount)
	if err != nil {
		return fmt.Errorf("error getting txData: %s", err.Error())
	}

	txHash, err := ec.SendTransaction(taAddress, txData)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	fmt.Printf("Made approve | Tx Hash: %s\n", txHash.String())
	return nil
}

func (ec *EVMClient) CheckAllowance(tokenAddress, spender, tokenAbi string, amount *big.Int) error {

	allowance, err := ec.GetAllowance(tokenAddress, spender, tokenAbi)
	if err != nil {
		return err
	}

	if allowance.Cmp(amount) == -1 {
		err = ec.Approve(tokenAddress, spender, tokenAbi, new(big.Int).Mul(amount, big.NewInt(2)))
		if err != nil {
			return err
		}
		return nil
	}

	fmt.Println("Enough allowance for transaction")

	return nil
}
