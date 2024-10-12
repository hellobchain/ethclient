package Client

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethclient/common/flogging"
	"github.com/ethclient/ethclient"
	key2 "github.com/ethclient/keystore/key"
	"github.com/ethclient/models"
	"github.com/ethclient/rpc"
)

var log = flogging.MustGetLogger("sipcclient.Client")

type SipcClient struct {
	// input
	Address string `json:"address"` // 节点的地址 IP+rpc port
	// output
	ClientPara *models.ClientPara  `json:"clientPara"` // client 参数
	SignPrikey *ecdsa.PrivateKey   `json:"signPrikey"` // 交易签名参数
	Ctx        *context.Context    `json:"ctx"`
	Cancel     *context.CancelFunc `json:"cancel"`
}

// new 一个client
func NewClient(address string, signTxPara *models.SignTxPara) (*SipcClient, error) {
	key, err := key2.GetKey(signTxPara.SignPrikeyFile, signTxPara.PasswdFile)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	client := &SipcClient{
		Address:    address,
		SignPrikey: key.PrivateKey,
		Ctx:        &ctx,
		Cancel:     &cancel,
	}
	err = client.clientInit()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// client初始化
func (c *SipcClient) clientInit() error {
	cli, err := rpc.Dial("http://"+c.Address, "", "", nil)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	c.ClientPara = &models.ClientPara{
		RpcClient: cli,
		Client:    ethclient.NewClient(cli),
	}
	return nil
}

// 关闭client
func (c *SipcClient) Close() {
	if c.ClientPara != nil {
		c.ClientPara.RpcClient.Close()
	}
	if c.ClientPara != nil {
		c.ClientPara.Client.Close()
	}
	if c.Cancel != nil {
		(*c.Cancel)()
	}
}
