package genesis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/ethclient/common"
	"github.com/ethclient/common/flogging"
	"github.com/ethclient/core"
	"github.com/ethclient/core/types"
	"github.com/ethclient/models"
	"github.com/ethclient/params"
	"github.com/ethclient/rlp"
)

var log = flogging.MustGetLogger("sipcclient.keystore.genesis")
var GenesisMap = map[string]int{
	"poa":  models.POA,
	"raft": models.RAFT,
	"pow":  models.POW,
}

// 生成创世区块
// para
// consensus int  #共识类型 三种 RAFT POW POA
// networkID int  #网络id
// packageBlocksTime int # 打块时间 秒为单位
// sealAddresses []string # 挖矿账号
// isAdd256Balances bool # 是否分配256 余额
// genesisFilePath string # 创世区块的文件路径
// gasPoolAccount string # 统一分配给其他管理员gas的gas池
// gasLimit uint64  # 每个块的最大gaslimit

func MakeGenesis(consensus int, networkID int, packageBlocksTime int, sealAddresses []string, isAdd256Balances bool, genesisFilePath string, gasLimit uint64) error {
	// Construct a default genesis block
	genesis := &core.Genesis{
		Timestamp:  uint64(time.Now().Unix()),
		GasLimit:   gasLimit,
		Difficulty: big.NewInt(0),
		Alloc:      make(core.GenesisAlloc),
		Config: &params.ChainConfig{
			SingularityBlock: big.NewInt(0),
		},
	}
	switch consensus {
	case models.RAFT:
		// In case of ethash, we're pretty much done
		genesis.Config.Raft = true
		genesis.ExtraData = make([]byte, 32)
	case models.POW:
		// In case of ethash, we're pretty much done
		genesis.Config.Ethash = new(params.EthashConfig)
		genesis.ExtraData = make([]byte, 32)
	case models.PBFT:
		// In the case of clique, configure the consensus parameters
		genesis.Difficulty = big.NewInt(1)
		genesis.Config.Pbft = &params.PbftConfig{
			ProposerPolicy: 0,
			Epoch:          30000,
		}
		var buf bytes.Buffer
		var extra []byte
		var addresses []common.Address
		for _, sealAddress := range sealAddresses {
			addresses = append(addresses, common.HexToAddress(sealAddress))
		}

		// compensate the lack bytes if header.Extra is not enough ByzantineExtraVanity bytes.
		if len(extra) < types.ByzantineExtraVanity {
			extra = append(extra, bytes.Repeat([]byte{0x00}, types.ByzantineExtraVanity-len(extra))...)
		}
		buf.Write(extra[:types.ByzantineExtraVanity])

		ist := &types.ByzantineExtra{
			Validators:    addresses,
			Seal:          []byte{},
			CommittedSeal: [][]byte{},
		}

		payload, err := rlp.EncodeToBytes(&ist)
		if err != nil {
			return err
		}

		extra = append(buf.Bytes(), payload...)
		genesis.ExtraData = extra
	case models.POA:
		// In the case of clique, configure the consensus parameters
		genesis.Difficulty = big.NewInt(1)
		genesis.Config.Clique = &params.CliqueConfig{
			Period: 15,
			Epoch:  30000,
		}
		genesis.Config.Clique.Period = uint64(packageBlocksTime)

		var signers []common.Address
		for _, sealAddress := range sealAddresses {
			if !common.IsHexAddress(sealAddress) {
				return fmt.Errorf("addr is not HexAddress")
			}
			if address := common.FromHex(sealAddress); address != nil {
				signers = append(signers, common.HexToAddress(sealAddress))
				continue
			}
			if len(signers) > 0 {
				break
			}
		}
		// Sort the signers and embed into the extra-data section
		for i := 0; i < len(signers); i++ {
			for j := i + 1; j < len(signers); j++ {
				if bytes.Compare(signers[i][:], signers[j][:]) > 0 {
					signers[i], signers[j] = signers[j], signers[i]
				}
			}
		}
		genesis.ExtraData = make([]byte, 32+len(signers)*common.AddressLength+65)
		for i, signer := range signers {
			copy(genesis.ExtraData[32+i*common.AddressLength:], signer[:])
		}

	default:
		log.Fatal("Invalid consensus engine choice", "choice", consensus)
	}
	for _, sealAddress := range sealAddresses {
		if !common.IsHexAddress(sealAddress) {
			return fmt.Errorf("addr is not HexAddress")
		}
		if address := common.FromHex(sealAddress); address != nil {
			genesis.Alloc[common.HexToAddress(sealAddress)] = core.GenesisAccount{
				Balance: new(big.Int).Lsh(big.NewInt(1), 256-7), // 2^256 / 128 (allow many pre-funds without balance overflows)
			}
			continue
		}
		break
	}

	if isAdd256Balances {
		// Add a batch of precompile balances to avoid them getting deleted
		for i := int64(0); i < 256; i++ {
			genesis.Alloc[common.BigToAddress(big.NewInt(i))] = core.GenesisAccount{Balance: big.NewInt(1)}
		}
	}
	// Query the user for some custom extras
	genesis.Config.ChainID = new(big.Int).SetUint64(uint64(networkID))
	out, _ := json.MarshalIndent(genesis, "", "  ")
	if err := ioutil.WriteFile(genesisFilePath, out, 0644); err != nil {
		return err
	}
	return nil
}
