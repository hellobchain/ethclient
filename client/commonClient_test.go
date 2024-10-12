package Client

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethclient/models"
)

func TestSignal(t *testing.T) {
	signal := make(chan bool, 100)
	signal <- true
	si := <-signal
	if si {
		fmt.Println("wsw")
	}
	signal <- true
	signal <- true
	for si := range signal {
		fmt.Println("wsw", si)
	}
}

var Client *EthClient

func GetClient(tls bool) *EthClient {
	address := "192.168.0.196:8548"

	signTxPara := &models.SignTxPara{
		SignPrikeyFile: "keystore/pri_key",
		PasswdFile:     "keystore/passwd.txt",
	}
	if Client != nil {
		return Client
	}
	sipcClient, err := NewClient(address, signTxPara)
	if err != nil {
		return nil
	}
	Client = sipcClient
	return Client
}

// 获取链上矿工账号
func TestGetSigners(t *testing.T) {
	signers, err := GetClient(true).GetSigners()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("%v\n", signers)
}

// 获取节点账户地址
func TestGetAccounts(t *testing.T) {
	accounts, err := GetClient(true).GetAccounts()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("%v\n", accounts)
}

// 获取节点信息
func TestGetNodeInfo(t *testing.T) {
	nodeInfos, err := GetClient(true).GetNodeInfo()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	res, _ := json.MarshalIndent(nodeInfos, "", "	")
	fmt.Printf("%v\n", string(res))
	fmt.Printf("%v\n", (nodeInfos.Protocols["eth"]).(map[string]interface{})["config"].(map[string]interface{})["clique"])
}

// 连接peer
func TestAddPeer(t *testing.T) {

}

// 最新块高
func TestBlockNumber(t *testing.T) {
	blockNumber, err := GetClient(true).BlockNumber()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("%v\n", blockNumber)
}

// 解锁账号
func TestUnlockAccount(t *testing.T) {
	accounts, err := GetClient(true).GetAccounts()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("%v\n", accounts)
	res, err := GetClient(true).UnlockAccount(accounts[0], "", 10000)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("%v\n", res)
}

// 连接的节点数
func TestPeerCount(t *testing.T) {
	peerCount, err := GetClient(true).PeerCount()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("%v\n", peerCount)
}

// 连接的节点信息
func TestPeers(t *testing.T) {
	peers, err := GetClient(true).Peers()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("%v\n", peers)
}

// 矿工投票 auth为false 删除矿工 为true 添加矿工 address为矿工账号
func TestPropose(t *testing.T) {
	accounts, err := GetClient(true).GetAccounts()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("%v\n", accounts)
	res, err := GetClient(true).Propose(accounts[0], true)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("%v\n", res)
}

// 获取余额
func TestGetBalance(t *testing.T) {
	accounts, err := GetClient(true).GetAccounts()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("%v\n", accounts)
	res, err := GetClient(true).GetBalance(accounts[0], "latest")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("%v\n", res)
}

// 通过交易ID获取交易信息
func TestGetTransactionByHash(t *testing.T) {
	res, err := GetClient(true).GetTransactionByHash("")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("%v\n", res)
}

// 通过块hash或者块高获取块
func TestGetBlockByBlockNumOrHash(t *testing.T) {
	res, err := GetClient(true).GetBlockByBlockNumOrHash("0")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	resMashal, _ := json.MarshalIndent(res, "", "	")
	fmt.Printf("%v\n", string(resMashal))
}

// 通过块hash或者块高获取块（返回块信息结构体信息不一样）
func TestGetMixedBlockByBlockNumOrHash(t *testing.T) {
	res, err := GetClient(true).GetMixedBlockByBlockNumOrHash("0")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	resMashal, _ := json.MarshalIndent(res, "", "	")
	fmt.Printf("%v\n", string(resMashal))
}

// 设置矿工账号
func TestSetEtherbase(t *testing.T) {
}

// 开启挖矿
func TestStartMiner(t *testing.T) {
}

// 开启挖矿
func TestMinerStart(t *testing.T) {
}

// 停止挖矿
func TestMinerStop(t *testing.T) {
}

// 挖矿状态
func TestMining(t *testing.T) {
}

// 获取交易receipt
func TestGetTransactionReceipt(t *testing.T) {
}

// 获取交易receipt
func TestGetTransactionDetail(t *testing.T) {
}
