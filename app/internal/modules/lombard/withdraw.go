package lombard

import (
	"fmt"
	contractAbis "lombardClaimer/internal/config/contract_abis"
	lombardAbi "lombardClaimer/internal/config/contract_abis/lombardAbi"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func (l *Lombard) Withdraw() error {

	balanceOfLbtc, err := l.evmClient.GetBalanceOf(l.LBTCvAddress, contractAbis.ERC20_ABI)
	if err != nil {
		return err
	}

	err = l.evmClient.CheckAllowance(
		l.LBTCvAddress,
		l.VaultAddress,
		contractAbis.ERC20_ABI,
		balanceOfLbtc,
	)

	if err != nil {
		return err
	}

	if l.evmClient.ChainName == "Base" {
		err = l.WithdrawFromBase(balanceOfLbtc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Lombard) WithdrawFromBase(amount *big.Int) error {

	abi, err := l.evmClient.GetAbi(lombardAbi.BASE_VAULT_ABI)
	if err != nil {
		return err
	}

	userReq := struct {
		Deadline    uint64
		AtomicPrice *big.Int
		OfferAmount *big.Int
		InSolve     bool
	}{
		Deadline:    uint64(time.Now().Add(72 * time.Hour).Unix()),
		AtomicPrice: big.NewInt(0),
		OfferAmount: amount,
		InSolve:     false,
	}

	txData, err := abi.Pack(
		"safeUpdateAtomicRequest",
		common.HexToAddress("0x5401b8620E5FB570064CA9114fd1e135fd77D57c"),
		common.HexToAddress("0xecAc9C5F704e954931349Da37F60E39f515c11c1"),
		userReq,
		common.HexToAddress("0x28634D0c5edC67CF2450E74deA49B90a4FF93dCE"),
		big.NewInt(100),
	)

	vaultAddressHash := common.HexToAddress(l.VaultAddress)

	hash, err := l.evmClient.SendTransaction(vaultAddressHash, txData)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Wallet: %s | Start Withdraw from lombard | Network: %s | Hash: %s\n\n",
		l.evmClient.Wallet.Address.Hex(),
		l.evmClient.ChainName,
		hash.Hex(),
	)

	return nil
}
