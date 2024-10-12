package Client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethclient/abi"
	"github.com/ethclient/common"
	"github.com/ethclient/common/compiler"
	"github.com/ethclient/common/hexutil"
	"github.com/ethclient/core/types"
	"github.com/ethclient/crypto"
	"github.com/ethclient/ethclient"
	"github.com/ethclient/models"
	"github.com/ethclient/rlp"
)

var versionRegexp = regexp.MustCompile(`([0-9]+)\.([0-9]+)\.([0-9]+)`)

// 调整solc编译参数
func makeArgs(s *compiler.Solidity) []string {
	p := []string{
		"--combined-json", "bin,bin-runtime,srcmap,srcmap-runtime,abi,userdoc,devdoc",
		"--optimize",                  // code optimizer switched on
		"--allow-paths", "., ./, ../", // default to support relative paths
	}
	if s.Major > 0 || s.Minor > 4 || s.Patch > 6 {
		p[1] += ",metadata,hashes"
	}
	return p
}

// 以docker形式获取solc版本
func solidityVersionForDocker(solcVersion string) (*compiler.Solidity, error) {
	var out bytes.Buffer
	tmp := fmt.Sprintf("docker run --rm -i --privileged=true --net=host --name solc ethereum/solc:%s --version", solcVersion)
	cmd := exec.Command("bash", "-c", tmp)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	matches := versionRegexp.FindStringSubmatch(out.String())
	if len(matches) != 4 {
		return nil, fmt.Errorf("can't parse solc version %q", out.String())
	}
	s := &compiler.Solidity{Path: cmd.Path, FullVersion: out.String(), Version: matches[0]}
	if s.Major, err = strconv.Atoi(matches[1]); err != nil {
		return nil, err
	}
	if s.Minor, err = strconv.Atoi(matches[2]); err != nil {
		return nil, err
	}
	if s.Patch, err = strconv.Atoi(matches[3]); err != nil {
		return nil, err
	}
	return s, nil
}

// 编译合约 基于docker形式 linux
func compilerContractForDocker(contractPath string, contractFileName string, solcVersion string) (map[string]models.ContractConfig, error) {
	if contractPath == "" || contractFileName == "" || solcVersion == "" {
		return nil, fmt.Errorf("para is empty")
	}
	_, err := os.Stat(filepath.Join(contractPath, contractFileName))
	if err != nil {
		return nil, err
	}
	s, err := solidityVersionForDocker(solcVersion)
	if err != nil {
		return nil, err
	}
	var compilerOptions string
	for _, arg := range makeArgs(s) {
		compilerOptions = fmt.Sprintf("%s %s", compilerOptions, arg)
	}
	var stderr, stdout bytes.Buffer
	arg := fmt.Sprintf("docker run --rm -i --privileged=true --workdir=/contract --net=host -v %s:/contract --name solc ethereum/solc:%s %s %s", contractPath, solcVersion, compilerOptions, contractFileName)
	cmd := exec.Command("bash", "-c", arg)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("solc: %v\n%s", err, stderr.Bytes())
	}
	contracts, err := compiler.ParseCombinedJSON(stdout.Bytes(), "", "", "", "")
	if err != nil {
		return nil, fmt.Errorf("Failed to read contract information from json output: %v\n", err)
	}
	if len(contracts) != 1 {
		return nil, fmt.Errorf("one contract expected, got %d", len(contracts))
	}
	contractMap := make(map[string]models.ContractConfig)
	for name, contract := range contracts {
		abiValue, err := json.MarshalIndent(contract.Info.AbiDefinition, "", "	") // Flatten the compiler parse
		if err != nil {
			return nil, fmt.Errorf("Failed to parse ABIs from compiler output: %v", err)
		}
		fmt.Printf("abiValue:%v\n", string(abiValue))
		contractMap[name] = models.ContractConfig{AbiData: abiValue, ContractCode: contract.Code}
	}
	return contractMap, nil
}

// 以命令行方式获取solc版本
func solidityVersion(solc string) (*compiler.Solidity, error) {
	if solc == "" {
		solc = "solc"
	}
	var out bytes.Buffer
	cmd := exec.Command(solc, "--version")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	matches := versionRegexp.FindStringSubmatch(out.String())
	if len(matches) != 4 {
		return nil, fmt.Errorf("can't parse solc version %q", out.String())
	}
	s := &compiler.Solidity{Path: cmd.Path, FullVersion: out.String(), Version: matches[0]}
	if s.Major, err = strconv.Atoi(matches[1]); err != nil {
		return nil, err
	}
	if s.Minor, err = strconv.Atoi(matches[2]); err != nil {
		return nil, err
	}
	if s.Patch, err = strconv.Atoi(matches[3]); err != nil {
		return nil, err
	}
	return s, nil
}

// 编译合约 以solc命令行形式
func compilerContract(solc, source string) (map[string]models.ContractConfig, error) {
	s, err := solidityVersion(solc)
	if err != nil {
		return nil, err
	}
	args := append(makeArgs(s), "--")
	cmd := exec.Command(s.Path, append(args, "-")...)
	cmd.Stdin = strings.NewReader(source)
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("solc: %v\n%s", err, stderr.Bytes())
	}
	contracts, err := compiler.ParseCombinedJSON(stdout.Bytes(), source, s.Version, s.Version, strings.Join(makeArgs(s), " "))
	if err != nil {
		return nil, fmt.Errorf("Failed to read contract information from json output: %v\n", err)
	}
	if len(contracts) != 1 {
		return nil, fmt.Errorf("one contract expected, got %d", len(contracts))
	}
	contractMap := make(map[string]models.ContractConfig)
	for name, contract := range contracts {
		abiValue, err := json.MarshalIndent(contract.Info.AbiDefinition, "", "	") // Flatten the compiler parse
		if err != nil {
			return nil, fmt.Errorf("Failed to parse ABIs from compiler output: %v", err)
		}
		fmt.Printf("abiValue:%v\n", string(abiValue))
		contractMap[name] = models.ContractConfig{AbiData: abiValue, ContractCode: contract.Code}
	}
	return contractMap, nil
}

// 调用合约
func (c *SipcClient) InvokeContract(contractAddressString string, abiData string, nonce uint64, method string, args ...interface{}) (string, error) {
	abiValue, err := abi.JSON(bytes.NewReader([]byte(abiData)))
	if err != nil {
		return "", err
	}
	out, err := abiValue.Pack(method, args...)
	if err != nil {
		return "", err
	}
	contractAddress := common.HexToAddress(contractAddressString)
	address := crypto.PubkeyToAddress(c.SignPrikey.PublicKey)
	gasPrice, err := c.ClientPara.Client.SuggestGasPrice(*c.Ctx)
	if err != nil {
		return "", err
	}
	msg := ethclient.CallMsg{
		From:     address,
		To:       &contractAddress,
		Data:     out,
		GasPrice: gasPrice,
	}
	gasLimit, err := c.ClientPara.Client.EstimateGas(*c.Ctx, msg)
	if err != nil {
		return "", err
	}

	transaction := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, out)
	transaction, err = types.SignTx(transaction, types.HomesteadSigner{}, c.SignPrikey)
	if err != nil {
		return "", err
	}
	content, err := rlp.EncodeToBytes(transaction)

	if err != nil {
		return "", err
	}
	var result common.Hash

	err = c.ClientPara.RpcClient.CallContext(*c.Ctx, &result, "eth_sendRawTransaction", hexutil.Bytes(content))

	if err != nil {
		return "", err
	}
	return result.String(), nil
}

// 查询合约
func (c *SipcClient) QueryContract(contractAddressString string, abiData string, result interface{}, method string, args ...interface{}) error {
	abiValue, err := abi.JSON(bytes.NewReader([]byte(abiData)))
	if err != nil {
		return err
	}
	out, err := abiValue.Pack(method, args...)
	if err != nil {
		return err
	}
	contractAddress := common.HexToAddress(contractAddressString)
	msg := ethclient.CallMsg{
		To:   &contractAddress,
		Data: out,
	}
	res, err := c.ClientPara.Client.CallContract(*c.Ctx, msg, nil)
	if err != nil {
		return err
	}
	return abiValue.Unpack(result, method, res)
}

// 部署合约
func (c *SipcClient) DeployContract(contractData string, contractName string) error {
	contractMap, err := compilerContract("", contractData)
	if err != nil {
		return err
	}
	i := 0
	nonce, err := c.GetNonce()
	if err != nil {
		return err
	}
	contractAddress := ""
	for _, contractData := range contractMap {
		txid, err := c.SendTransaction(models.CREATE_CONTRACT, nonce+uint64(i), "", "", common.FromHex(contractData.ContractCode))
		if err != nil {
			return err
		}
		opType := "DeployContract"
		var wg sync.WaitGroup
		txResultStatus := make(chan models.TxResultStatus, 1)
		wg.Add(1)
		go c.JudgeUpChainStatus(*txid, opType, txResultStatus, &wg)
		wg.Wait()
		select {
		case receipt := <-txResultStatus:
			if receipt.Err != nil {
				return err
			}
			contractAddress = receipt.Receipt.ContractAddress.Hex()
		default:
			return fmt.Errorf(opType)
		}
		i++
	}
	log.Infof("contractAddress:%s,contractName:%s", contractAddress, contractName)
	return nil
}

// 判断上链状态
func (c *SipcClient) JudgeUpChainStatus(txId string, opType string, txResultStatusChan chan models.TxResultStatus, wg *sync.WaitGroup) {
	defer wg.Done()
	txResultStatus := models.TxResultStatus{}
	receipt, err := c.internalJudgeUpChainStatus(txId, opType)
	if err == nil {
		txResultStatus.Receipt = *receipt
	}
	txResultStatus.Err = err
	txResultStatusChan <- txResultStatus
}

// 判断上链状态
func (c *SipcClient) internalJudgeUpChainStatus(txId string, opType string) (*types.Receipt, error) {
	ticker := time.NewTicker(30 * time.Second)
	meterCount := 0
	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			return nil, fmt.Errorf("get txId:%v opType(%v) timeout(30s)", txId, opType)
		default:
			receipt, err := c.GetTransactionReceipt(txId)
			if err != nil {
				return nil, err
			}
			if receipt == nil {
				if meterCount >= 30 {
					return nil, fmt.Errorf("get txId:%v opType(%v) times(30)", txId, opType)
				}
				time.Sleep(1 * time.Second)
				meterCount++
				continue
			}
			if receipt.Status == 0 {
				return nil, fmt.Errorf("%v up chain success but transaction is invalid,err:%v", opType, receipt.Bloom)
			}
			return receipt, nil
		}
	}
}
