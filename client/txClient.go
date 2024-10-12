package Client

import (
	"fmt"
	"math/big"

	"github.com/ethclient/common"
	"github.com/ethclient/core/types"
	"github.com/ethclient/crypto"
	"github.com/ethclient/ethclient"
	"github.com/ethclient/models"
)

// 发送交易
func (c *EthClient) SendTransaction(opType int, nonce uint64, to string, amount string, data []byte) (*string, error) {
	var rawTx *types.Transaction
	amountBigInt, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		amountBigInt = new(big.Int)
	}
	from := crypto.PubkeyToAddress(c.SignPrikey.PublicKey)
	gasPrice, err := c.ClientPara.Client.SuggestGasPrice(*c.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %v", err)
	}
	var gasLimit uint64
	switch opType {
	case models.CREATE_CONTRACT:
		msg := ethclient.CallMsg{From: from, To: nil, GasPrice: gasPrice, Value: amountBigInt, Data: data}
		gasLimit, err = c.ClientPara.Client.EstimateGas(*c.Ctx, msg)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas needed: %v", err)
		}
		rawTx = types.NewContractCreation(nonce, amountBigInt, gasLimit, gasPrice, data)
	case models.NORMAL_TRANSACTION:
		if !common.IsHexAddress(to) {
			return nil, fmt.Errorf("%s is not HexAddress", to)
		}
		to := common.HexToAddress(to)
		msg := ethclient.CallMsg{From: from, To: &to, GasPrice: gasPrice, Value: amountBigInt, Data: data}
		gasLimit, err = c.ClientPara.Client.EstimateGas(*c.Ctx, msg)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas needed: %v", err)
		}
		rawTx = types.NewTransaction(nonce, to, amountBigInt, gasLimit, gasPrice, data)
	}
	signedTx, err := types.SignTx(rawTx, types.HomesteadSigner{}, c.SignPrikey)
	if err != nil {
		return nil, err
	}
	if err := c.ClientPara.Client.SendTransaction(*c.Ctx, signedTx); err != nil {
		return nil, err
	}
	txHash := signedTx.Hash().Hex()
	return &txHash, nil
}
