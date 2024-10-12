package Client

import (
	"encoding/json"
	"fmt"
	"math/big"
	"runtime"
	"strconv"
	"strings"

	"github.com/ethclient/common"
	"github.com/ethclient/core/types"
	"github.com/ethclient/crypto"
	"github.com/ethclient/models"
	"github.com/shopspring/decimal"
)

// to 16进制串
func ToHexString(q string) (*string, error) {
	v, success := big.NewInt(0).SetString(q, 0)
	if !success {
		return nil, fmt.Errorf("invalid input")
	}
	if v.Cmp(big.NewInt(0)) < 0 {
		return nil, fmt.Errorf("quantity must be larger than -1")
	}
	hexString := "0x" + v.Text(16)
	return &hexString, nil
}

// Block to MixedBlock
func BlockToMixedBlock(c *SipcClient, block *models.Block) (*models.MixedBlock, error) {
	mixedTransactions := make([]models.MixTransaction, 0)
	for _, tx := range block.Transactions {
		receipt, err := c.GetTransactionDetail(tx.Hash)
		if err != nil {

			return nil, err
		}
		mixedTx, err := TransactionToMixedTransaction(&tx, receipt)
		if err != nil {

			return nil, err
		}
		mixedTransactions = append(mixedTransactions, *mixedTx)
	}
	var mixedBlock models.MixedBlock
	mixedBlock.Bloom = block.Bloom
	mixedBlock.Coinbase = block.Coinbase
	mixedBlock.Difficulty = block.Difficulty
	mixedBlock.Extra = block.Extra
	mixedBlock.GasLimit = block.GasLimit
	mixedBlock.Hash = block.Hash
	mixedBlock.MixDigest = block.MixDigest
	mixedBlock.Nonce = block.Nonce
	mixedBlock.Number = block.Number
	mixedBlock.ReceiptHash = block.ReceiptHash
	mixedBlock.ParentHash = block.ParentHash
	mixedBlock.Root = block.Root
	mixedBlock.Size = block.Size
	mixedBlock.Time = block.Time
	mixedBlock.TotalDifficulty = block.TotalDifficulty
	mixedBlock.Transactions = mixedTransactions
	mixedBlock.TxHash = block.TxHash
	mixedBlock.UncleHash = block.UncleHash
	mixedBlock.Uncles = block.Uncles
	mixedBlock.GasUsed = block.GasUsed
	return &mixedBlock, nil
}

// Transaction to MixTransaction
func TransactionToMixedTransaction(transaction *models.Transaction, receipt *models.Receipt) (*models.MixTransaction, error) {
	var mixedTransaction models.MixTransaction
	mixedTransaction.Tx = *transaction
	mixedTransaction.Status = receipt.Status
	mixedTransaction.ContractAddress = receipt.ContractAddress
	price, err := strconv.ParseInt(transaction.Price[2:], 16, 64)
	if err != nil {

		return nil, err
	}

	gasUsed, err := strconv.ParseInt(receipt.GasUsed, 10, 64)
	if err != nil {

		return nil, err
	}
	sipc := decimal.NewFromFloat((float64)(price)).Mul(decimal.NewFromFloat((float64)(gasUsed))).Div(decimal.NewFromFloat(10).Pow(decimal.NewFromFloat(18)))
	mixedTransaction.Sipc = sipc.String()
	return &mixedTransaction, nil
}

// 块转换成前端要的结构
func ExchangeBlock(inputBlock interface{}) error {
	switch inputBlock.(type) {
	case *models.Block:
		block := inputBlock.(*models.Block)
		err := ExchangeBlockHeader(block)
		if err != nil {
			return err
		}
		for i := range block.Uncles {
			err = ExchangeUncleHeader(&block.Uncles[i])
			if err != nil {
				return err
			}
		}
		for i := range block.Transactions {
			err = ExchangeBlockTransaction(&block.Transactions[i])
			if err != nil {
				return err
			}
		}
	case *models.MixedBlock:
		block := inputBlock.(*models.MixedBlock)
		err := ExchangeBlockHeader(block)
		if err != nil {
			return err
		}
		for i := range block.Uncles {
			err = ExchangeUncleHeader(&block.Uncles[i])
			if err != nil {
				return err
			}
		}
		for i := range block.Transactions {
			err = ExchangeBlockTransaction(&block.Transactions[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 块头转换成前端要的结构
func ExchangeBlockHeader(inputBlock interface{}) error {
	switch inputBlock.(type) {
	case *models.Block:
		block := inputBlock.(*models.Block)
		tmp, err := HexStringToDecString(block.GasLimit)
		if err != nil {
			return err
		}
		block.GasLimit = *tmp
		tmp, err = HexStringToDecString(block.Difficulty)
		if err != nil {
			return err
		}
		block.Difficulty = *tmp

		tmp, err = HexStringToDecString(block.GasUsed)
		if err != nil {
			return err
		}
		block.GasUsed = *tmp

		tmp, err = HexStringToDecString(block.Nonce)
		if err != nil {
			return err
		}
		block.Nonce = *tmp

		tmp, err = HexStringToDecString(block.Number)
		if err != nil {
			return err
		}
		block.Number = *tmp

		tmp, err = HexStringToDecString(block.Size)
		if err != nil {
			return err
		}
		block.Size = *tmp

		tmp, err = HexStringToDecString(block.Time)
		if err != nil {
			return err
		}
		block.Time = *tmp

		tmp, err = HexStringToDecString(block.TotalDifficulty)
		if err != nil {
			return err
		}
		block.TotalDifficulty = *tmp
	case *models.MixedBlock:
		block := inputBlock.(*models.MixedBlock)
		tmp, err := HexStringToDecString(block.GasLimit)
		if err != nil {
			return err
		}
		block.GasLimit = *tmp
		tmp, err = HexStringToDecString(block.Difficulty)
		if err != nil {
			return err
		}
		block.Difficulty = *tmp

		tmp, err = HexStringToDecString(block.GasUsed)
		if err != nil {
			return err
		}
		block.GasUsed = *tmp

		tmp, err = HexStringToDecString(block.Nonce)
		if err != nil {
			return err
		}
		block.Nonce = *tmp

		tmp, err = HexStringToDecString(block.Number)
		if err != nil {
			return err
		}
		block.Number = *tmp

		tmp, err = HexStringToDecString(block.Size)
		if err != nil {
			return err
		}
		block.Size = *tmp

		tmp, err = HexStringToDecString(block.Time)
		if err != nil {
			return err
		}
		block.Time = *tmp

		tmp, err = HexStringToDecString(block.TotalDifficulty)
		if err != nil {
			return err
		}
		block.TotalDifficulty = *tmp
	}

	return nil
}

// 块uncle头转换成前端要的结构
func ExchangeUncleHeader(header *models.Header) error {
	tmp, err := HexStringToDecString(header.GasLimit)
	if err != nil {
		return err
	}
	header.GasLimit = *tmp
	tmp, err = HexStringToDecString(header.Difficulty)
	if err != nil {
		return err
	}
	header.Difficulty = *tmp

	tmp, err = HexStringToDecString(header.GasUsed)
	if err != nil {
		return err
	}
	header.GasUsed = *tmp

	tmp, err = HexStringToDecString(header.Nonce)
	if err != nil {
		return err
	}
	header.Nonce = *tmp

	tmp, err = HexStringToDecString(header.Number)
	if err != nil {
		return err
	}
	header.Number = *tmp

	tmp, err = HexStringToDecString(header.Size)
	if err != nil {
		return err
	}
	header.Size = *tmp

	tmp, err = HexStringToDecString(header.Time)
	if err != nil {
		return err
	}
	header.Time = *tmp

	tmp, err = HexStringToDecString(header.TotalDifficulty)
	if err != nil {
		return err
	}
	header.TotalDifficulty = *tmp
	return nil
}

// 块里交易转换成前端要的结构
func ExchangeBlockTransaction(inputTransaction interface{}) error {
	switch inputTransaction.(type) {
	case *models.Transaction:
		transaction := inputTransaction.(*models.Transaction)
		tmp, err := HexStringToDecString(transaction.BlockNumber)
		if err != nil {
			return err
		}
		transaction.BlockNumber = *tmp
		tmp, err = HexStringToDecString(transaction.AccountNonce)
		if err != nil {
			return err
		}
		transaction.AccountNonce = *tmp

		tmp, err = HexStringToDecString(transaction.Amount)
		if err != nil {
			return err
		}
		transaction.Amount = *tmp

		tmp, err = HexStringToDecString(transaction.GasLimit)
		if err != nil {
			return err
		}
		transaction.GasLimit = *tmp

		tmp, err = HexStringToDecString(transaction.Price)
		if err != nil {
			return err
		}
		transaction.Price = *tmp

		tmp, err = HexStringToDecString(transaction.TransactionIndex)
		if err != nil {
			return err
		}
		transaction.TransactionIndex = *tmp
	case *models.MixTransaction:
		transaction := inputTransaction.(*models.MixTransaction)
		tmp, err := HexStringToDecString(transaction.Tx.BlockNumber)
		if err != nil {
			return err
		}
		transaction.Tx.BlockNumber = *tmp
		tmp, err = HexStringToDecString(transaction.Tx.AccountNonce)
		if err != nil {
			return err
		}
		transaction.Tx.AccountNonce = *tmp

		tmp, err = HexStringToDecString(transaction.Tx.Amount)
		if err != nil {
			return err
		}
		transaction.Tx.Amount = *tmp

		tmp, err = HexStringToDecString(transaction.Tx.GasLimit)
		if err != nil {
			return err
		}
		transaction.Tx.GasLimit = *tmp

		tmp, err = HexStringToDecString(transaction.Tx.Price)
		if err != nil {
			return err
		}
		transaction.Tx.Price = *tmp

		tmp, err = HexStringToDecString(transaction.Tx.TransactionIndex)
		if err != nil {
			return err
		}
		transaction.Tx.TransactionIndex = *tmp
	}
	return nil
}

// 块里receipt转换成前端要的结构
func ExchangeBlockReceipt(receipt *models.Receipt) error {
	tmp, err := HexStringToDecString(receipt.BlockNumber)
	if err != nil {
		return err
	}
	receipt.BlockNumber = *tmp
	tmp, err = HexStringToDecString(receipt.CumulativeGasUsed)
	if err != nil {
		return err
	}
	receipt.CumulativeGasUsed = *tmp

	tmp, err = HexStringToDecString(receipt.GasUsed)
	if err != nil {
		return err
	}
	receipt.GasUsed = *tmp

	tmp, err = HexStringToDecString(receipt.TransactionIndex)
	if err != nil {
		return err
	}
	receipt.TransactionIndex = *tmp
	tmp, err = HexStringToDecString(receipt.Status)
	if err != nil {
		return err
	}
	receipt.Status = *tmp
	return nil
}

// 16进制转换
func HexStringToDecString(hexString string) (*string, error) {
	v, success := big.NewInt(0).SetString(hexString, 0)
	if !success {
		return nil, fmt.Errorf("invalid input")
	}
	if v.Cmp(big.NewInt(0)) < 0 {
		return nil, fmt.Errorf("hex must be larger than 0")
	}
	decString := v.Text(10)
	return &decString, nil
}

// 获取链上矿工账号
func (c *SipcClient) GetSigners() ([]string, error) {
	var signers []string
	err := c.ClientPara.RpcClient.Call(&signers, "clique_getSigners")
	if err != nil {

		return nil, err
	}
	return signers, nil
}

// 获取节点账户地址
func (c *SipcClient) GetAccounts() ([]string, error) {
	var account []string
	err := c.ClientPara.RpcClient.Call(&account, "eth_accounts")
	if err != nil {

		return nil, err
	}
	return account, nil
}

// 获取节点信息
func (c *SipcClient) GetNodeInfo() (*models.NodeInfo, error) {
	node := models.NodeInfo{}
	err := c.ClientPara.RpcClient.Call(&node, "admin_nodeInfo")
	if err != nil {

		return nil, err
	}
	return &node, nil
}

// 连接peer
func (c *SipcClient) AddPeer(enode string, from string) (string, error) {
	var ok interface{}
	if enode == "" {

		return "", fmt.Errorf("enode is empty")
	}
	if !common.IsHexAddress(from) {

		return "", fmt.Errorf("addr is not HexAddress")
	}
	err := c.ClientPara.RpcClient.Call(&ok, "permission_addPeer", enode, common.HexToAddress(from))
	if err != nil {

		return "", err
	}
	if ok.(string) == "false" {
		return "", fmt.Errorf("addPeer enode %v failed", enode)
	}
	return ok.(string), nil
}

// 最新块高
func (c *SipcClient) BlockNumber() (uint64, error) {
	var n string
	err := c.ClientPara.RpcClient.Call(&n, "eth_blockNumber")
	if err != nil {

		return 0, err
	}
	afterN := strings.Split(n, "0x")
	blockNumber, err := strconv.ParseInt(afterN[1], 16, 64)
	if err != nil {

		return 0, err
	}
	return uint64(blockNumber), nil
}

// 解锁账号
func (c *SipcClient) UnlockAccount(account string, passwd string, time uint64) (bool, error) {
	var res interface{}
	err := c.ClientPara.RpcClient.Call(&res, "personal_unlockAccount", account, passwd, time)
	if err != nil {

		return res.(bool), err
	}
	return res.(bool), nil
}

// 连接的节点数
func (c *SipcClient) PeerCount() (int64, error) {
	var res interface{}
	err := c.ClientPara.RpcClient.Call(&res, "net_peerCount")
	if err != nil {

		return 0, err
	}
	afterN := strings.Split(res.(string), "0x")
	peerCount, err := strconv.ParseInt(afterN[1], 16, 64)
	return peerCount, nil
}

// 连接的节点信息
func (c *SipcClient) Peers() (interface{}, error) {
	var res interface{}
	err := c.ClientPara.RpcClient.Call(&res, "admin_peers")
	if err != nil {

		return 0, err
	}
	return res, nil
}

// 矿工投票 auth为false 删除矿工 为true 添加矿工 address为矿工账号
func (c *SipcClient) Propose(address string, auth bool) (interface{}, error) {
	var res interface{}
	err := c.ClientPara.RpcClient.Call(&res, "clique_propose", common.HexToAddress(address), auth)
	if err != nil {

		return 0, err
	}
	return res, nil
}

// 获取余额
func (c *SipcClient) GetBalance(addr string, status string) (string, error) {
	var balance string
	err := c.ClientPara.RpcClient.Call(&balance, "eth_getBalance", addr, status)
	if err != nil {

		return balance, err
	}
	return balance, nil
}

// 通过交易ID获取交易信息
func (c *SipcClient) GetTransactionByHash(txId string) (*models.Transaction, error) {
	var res interface{}
	err := c.ClientPara.RpcClient.Call(&res, "eth_getTransactionByHash", txId)
	if err != nil {

		return nil, err
	}
	transactionInfo, err := json.MarshalIndent(res, "", "	")
	if err != nil {

		return nil, err
	}
	var transaction models.Transaction
	err = json.Unmarshal(transactionInfo, &transaction)
	if err != nil {

		return nil, err
	}
	err = ExchangeBlockTransaction(&transaction)
	if err != nil {

		return nil, err
	}
	return &transaction, nil
}

// 通过块hash或者块高获取块
func (c *SipcClient) GetBlockByBlockNumOrHash(input string) (*models.Block, error) {
	var arg string
	var method string
	tmp := strings.Split(input, "0x")
	switch len(tmp) {
	case 1: // 是块高
		method = "eth_getBlockByNumber"
		hexString, err := ToHexString(input)
		if err != nil {

			return nil, err
		}
		arg = *hexString
	case 2: // 是hash
		method = "eth_getBlockByHash"
		arg = input
	default:
		funcName, file, line, ok := runtime.Caller(0)
		if ok {
			log.Error("func name: " + runtime.FuncForPC(funcName).Name())
			log.Error("file:", file, "line:", line)
		}
		log.Error("arg is invalid")
		return nil, fmt.Errorf("arg is invalid")
	}
	var res interface{}
	err := c.ClientPara.RpcClient.Call(&res, method, arg, true)
	if err != nil {

		return nil, err
	}
	if res == nil {
		funcName, file, line, ok := runtime.Caller(0)
		if ok {
			log.Error("func name: " + runtime.FuncForPC(funcName).Name())
			log.Error("file:", file, "line:", line)
		}
		return nil, fmt.Errorf("当前的提供的参数%v链上不存在对应的块", arg)
	}
	blockInfo, err := json.MarshalIndent(res, "", "	")
	if err != nil {

		return nil, err
	}
	var block models.Block
	err = json.Unmarshal(blockInfo, &block)
	if err != nil {

		return nil, err
	}
	err = ExchangeBlock(&block)
	if err != nil {

		return nil, err
	}
	return &block, nil
}

// 通过块hash或者块高获取块（返回块信息结构体信息不一样）
func (c *SipcClient) GetMixedBlockByBlockNumOrHash(input string) (*models.MixedBlock, error) {
	var arg string
	var method string
	tmp := strings.Split(input, "0x")
	switch len(tmp) {
	case 1: // 是块高
		method = "eth_getBlockByNumber"
		hexString, err := ToHexString(input)
		if err != nil {

			return nil, err
		}
		arg = *hexString
	case 2: // 是hash
		method = "eth_getBlockByHash"
		arg = input
	default:
		funcName, file, line, ok := runtime.Caller(0)
		if ok {
			log.Error("func name: " + runtime.FuncForPC(funcName).Name())
			log.Error("file:", file, "line:", line)
		}
		log.Error("arg is invalid")
		return nil, fmt.Errorf("arg is invalid")
	}
	var res interface{}
	err := c.ClientPara.RpcClient.Call(&res, method, arg, true)
	if err != nil {

		return nil, err
	}
	if res == nil {

		return nil, fmt.Errorf("当前的提供的参数%v链上不存在对应的块", arg)
	}
	blockInfo, err := json.MarshalIndent(res, "", "	")
	if err != nil {

		return nil, err
	}
	var block models.Block
	err = json.Unmarshal(blockInfo, &block)
	if err != nil {

		return nil, err
	}
	mixedBlock, err := BlockToMixedBlock(c, &block)
	if err != nil {

		return nil, err
	}
	err = ExchangeBlock(mixedBlock)
	if err != nil {

		return nil, err
	}
	return mixedBlock, nil
}

// 设置矿工账号
func (c *SipcClient) SetEtherbase(address string) (bool, error) {
	ok := false
	err := c.ClientPara.RpcClient.Call(&ok, "miner_setEtherbase", address)
	if err != nil {

		return false, err
	}
	return ok, nil
}

// 开启挖矿
func (c *SipcClient) StartMiner() {
	a, err := c.GetAccounts()
	if err != nil {

		return
	}
	ret, err := c.SetEtherbase(a[0])
	if err != nil || !ret {
		funcName, file, line, ok := runtime.Caller(0)
		if ok {
			log.Error("func name: " + runtime.FuncForPC(funcName).Name())
			log.Error("file:", file, "line:", line)
		}
		if err != nil {
			log.Error(err.Error())
		}
		return
	}
	_, err = c.MinerStart()
	if err != nil {

		return
	}
}

// 开启挖矿
func (c *SipcClient) MinerStart() (interface{}, error) {
	var res interface{}
	err := c.ClientPara.RpcClient.Call(&res, "miner_start")
	if err != nil {

		return nil, err
	}
	return res, nil
}

// 停止挖矿
func (c *SipcClient) MinerStop() (interface{}, error) {
	var res interface{}
	err := c.ClientPara.RpcClient.Call(&res, "miner_stop")
	if err != nil {

		return nil, err
	}
	return res, nil
}

// 挖矿状态
func (c *SipcClient) Mining() (interface{}, error) {
	var res interface{}
	err := c.ClientPara.RpcClient.Call(&res, "eth_mining")
	if err != nil {

		return nil, err
	}
	return res, nil
}

// 获取交易receipt
func (c *SipcClient) GetTransactionReceipt(txId string) (*types.Receipt, error) {
	var r *types.Receipt
	err := c.ClientPara.RpcClient.Call(&r, "eth_getTransactionReceipt", common.HexToHash(txId))
	if err != nil {

		return nil, err
	}
	return r, nil
}

// 获取交易receipt
func (c *SipcClient) GetTransactionDetail(txId string) (*models.Receipt, error) {
	var res interface{}
	err := c.ClientPara.RpcClient.Call(&res, "eth_getTransactionReceipt", common.HexToHash(txId))
	if err != nil {
		return nil, err
	}
	receiptInfo, err := json.MarshalIndent(res, "", "	")
	if err != nil {

		return nil, err
	}
	var receipt models.Receipt
	err = json.Unmarshal(receiptInfo, &receipt)
	if err != nil {

		return nil, err
	}
	err = ExchangeBlockReceipt(&receipt)
	if err != nil {

		return nil, err
	}
	return &receipt, nil
}

// 获取共识类型 poa raft pow
func (c *SipcClient) GetConsensus() (string, error) {
	var consensus string
	nodeInfo, err := c.GetNodeInfo()
	if err != nil {
		return consensus, err
	}
	if nodeInfo.Protocols["eth"].(map[string]interface{})["config"].(map[string]interface{})["clique"] != nil {
		consensus = "poa"
	} else if nodeInfo.Protocols["eth"].(map[string]interface{})["config"].(map[string]interface{})["raft"] != nil {
		consensus = "raft"
	} else if nodeInfo.Protocols["eth"].(map[string]interface{})["config"].(map[string]interface{})["ethash"] != nil {
		consensus = "pow"
	} else if nodeInfo.Protocols["eth"].(map[string]interface{})["config"].(map[string]interface{})["scrypt"] != nil {
		consensus = "scrypt"
	} else {
		consensus = "unkown"
	}
	return consensus, nil
}

// 获取nonce
func (c *SipcClient) GetNonce() (uint64, error) {
	from := crypto.PubkeyToAddress(c.SignPrikey.PublicKey)
	return c.ClientPara.Client.PendingNonceAt(*c.Ctx, from)
}

// 调用RPC API
func (c *SipcClient) CallRpcApi(method string, para ...interface{}) (interface{}, error) {
	var res interface{}
	err := c.ClientPara.RpcClient.Call(&res, method, para)
	if err != nil {
		return nil, err
	}
	return res, nil
}
