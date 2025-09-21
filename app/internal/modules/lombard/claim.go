package lombard

import (
	"encoding/json"
	"fmt"
	"io"
	lombardabi "lombardClaimer/internal/config/contract_abis/lombardAbi"
	"lombardClaimer/internal/modules/lombard/entities"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var LOMBARD_HEADERS = map[string]string{
	"accept":             "*/*",
	"accept-language":    "ru,en-US;q=0.9,en;q=0.8,zh-TW;q=0.7,zh;q=0.6",
	"content-type":       "application/json",
	"origin":             "https://claim.lombard.finance",
	"priority":           "u=1, i",
	"referer":            "https://claim.lombard.finance/",
	"sec-ch-ua":          "\"Chromium\";v=\"140\", \"Not=A?Brand\";v=\"24\", \"Google Chrome\";v=\"140\"",
	"sec-ch-ua-mobile":   "?0",
	"sec-ch-ua-platform": "\"Windows\"",
	"sec-fetch-dest":     "empty",
	"sec-fetch-mode":     "cors",
	"sec-fetch-site":     "same-site",
	"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36",
}

func (l *Lombard) GetSignMessage() (*entities.ResponseDistributionData, error) {

	url := fmt.Sprintf("https://mainnet.prod.lombard.finance/api/v1/bard/distributor/WAVE_1/%s/recipient", l.evmClient.Wallet.Address.Hex())

	resp, err := l.httpClient.Get(url, LOMBARD_HEADERS)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var MsgData entities.ResponseDistributionData

	if err := json.Unmarshal(body, &MsgData); err != nil {
		return nil, nil
	}

	return &MsgData, nil
}

func (l *Lombard) ConfirmSign(payload *entities.SignConfirmationData) error {

	url := "https://mainnet.prod.lombard.finance/api/v1/bard/distributor/sign-terms"

	resp, err := l.httpClient.Post(url, LOMBARD_HEADERS, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var SignResponse entities.SignConfirmationResponse

	if err := json.Unmarshal(body, &SignResponse); err != nil {
		return err
	}

	if !SignResponse.Status {
		return fmt.Errorf("error sign lombard distribution message | Status: %v", SignResponse.Status)
	}

	return nil
}

func (l *Lombard) GetHardClaims() (*entities.HardClaimResponse, error) {

	url := fmt.Sprintf("https://mainnet.prod.lombard.finance/api/v1/bard/distributor/%s/hard-claims", l.evmClient.Wallet.Address.Hex())

	resp, err := l.httpClient.Get(url, LOMBARD_HEADERS)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ClaimsData entities.HardClaimResponse

	if err := json.Unmarshal(body, &ClaimsData); err != nil {
		return nil, err
	}

	return &ClaimsData, nil
}

func (l *Lombard) MakeClaim(claimAddress string) error {

	distributionData, err := l.GetSignMessage()
	if err != nil {
		return err
	}

	if !distributionData.Status {
		return fmt.Errorf("wallet: %s | cat get lombard distribution data", l.evmClient.Wallet.Address.Hex())
	}

	if distributionData.Recipient.AcceptedAt == "0001-01-01T00:00:00Z" {
		signedMsg, err := l.evmClient.SignMessage(distributionData.Recipient.Msg)
		if err != nil {
			return err
		}

		signConfirmationPayload := entities.SignConfirmationData{
			Address: distributionData.Recipient.Address,
			Slug:    distributionData.Recipient.OrgSlug,
			Sign:    signedMsg,
		}

		if err := l.ConfirmSign(&signConfirmationPayload); err != nil {
			return err
		}

		fmt.Printf("Wallet: %s | Sign message for BARD claim\n\n", l.evmClient.Wallet.Address.Hex())
	}

	claimsData, err := l.GetHardClaims()
	if err != nil {
		return err
	}

	claimAbi, err := l.evmClient.GetAbi(lombardabi.LOMBARD_CLAIM_ABI)
	if err != nil {
		return nil
	}

	amount := new(big.Int)
	amount, ok := amount.SetString(claimsData.Claims[0].Amount, 10)
	if !ok {
		return fmt.Errorf("cant parse claim amount to *Big.Int")
	}

	hashes := []common.Hash{}

	for _, s := range claimsData.Claims[0].Proofs.Hashes {
		hashes = append(hashes, common.HexToHash(s))
	}

	txData, err := claimAbi.Pack(
		"claim",
		l.evmClient.Wallet.Address,
		amount,
		hashes,
	)

	txHash, err := l.evmClient.SendTransaction(
		common.HexToAddress(claimAddress),
		txData,
	)

	if err != nil {
		return err
	}

	fmt.Printf(
		"Wallet: %s | Claim Bard Allocation | Hash: %s\n\n",
		l.evmClient.Wallet.Address.Hex(),
		txHash.Hex(),
	)

	return nil
}
