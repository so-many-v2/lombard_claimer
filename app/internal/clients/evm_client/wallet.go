package evmclient

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
}

func NewWallet(privateKey string) (*Wallet, error) {

	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, fmt.Errorf("error decode privat key: %s", err.Error())
	}

	pub := pk.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*pub)

	return &Wallet{
		PrivateKey: pk,
		Address:    address,
	}, nil
}
