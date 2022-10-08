package biz

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"
	"github.com/leondevpt/wallet/trxservice/global"
	"unicode/utf8"

	"github.com/fbsobreira/gotron-sdk/pkg/abi"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	trxapi "github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

const (
	trc20TransferMethodSignature = "0xa9059cbb"
	trc20TransferEventSignature  = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	trc20NameSignature           = "0x06fdde03"
	trc20SymbolSignature         = "0x95d89b41"
	trc20DecimalsSignature       = "0x313ce567"
	trc20BalanceOf               = "0x70a08231"
)

type TronCli struct {
	Conn          *grpc.ClientConn
	TronWalletCli trxapi.WalletClient
	ApiKey        string
	GrpcTimeout   time.Duration
}

func NewTronCli() *TronCli {
	// 建立对应 scheme 的连接, 并且配置负载均衡
	conn, err := grpc.Dial(fmt.Sprintf("%s:///%s", global.StaticSchema, global.TronNode), // "static:///tron_node"
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		panic(err)
	}
	cli := trxapi.NewWalletClient(conn)

	defaultTimeout := 30 * time.Second
	return &TronCli{TronWalletCli: cli, Conn: conn, GrpcTimeout: defaultTimeout}
}

func (t *TronCli) Stop() {
	if t.Conn != nil {
		t.Conn.Close()
	}
}

// SetTimeout for Client connections
func (c *TronCli) SetTimeout(timeout time.Duration) {
	c.GrpcTimeout = timeout
}

// SetAPIKey enable API on connection
func (c *TronCli) SetAPIKey(apiKey string) error {
	c.ApiKey = apiKey
	return nil
}

func (c TronCli) GetBalance(ctx context.Context, addr string) (*big.Int, error) {
	var (
		err     error
		account = new(core.Account)
	)
	account.Address, err = common.DecodeCheck(addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	acc, err := c.TronWalletCli.GetAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(acc.Address, account.Address) {
		return nil, fmt.Errorf("account not found")
	}
	return big.NewInt(acc.Balance), nil
}

func (c TronCli) GetTRC20TokenBalance(ctx context.Context, from, contractAddr string) (*big.Int, error) {
	return c.TRC20ContractBalance(from, contractAddr)
}

// TRC20Call make cosntant calll
func (c *TronCli) TRC20Call(from, contractAddress, data string, constant bool, feeLimit int64) (*api.TransactionExtention, error) {
	var err error
	fromDesc := address.HexToAddress("410000000000000000000000000000000000000000")
	if len(from) > 0 {
		fromDesc, err = address.Base58ToAddress(from)
		if err != nil {
			return nil, err
		}
	}
	contractDesc, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}
	dataBytes, err := common.FromHex(data)
	if err != nil {
		return nil, err
	}
	ct := &core.TriggerSmartContract{
		OwnerAddress:    fromDesc.Bytes(),
		ContractAddress: contractDesc.Bytes(),
		Data:            dataBytes,
	}
	result := &api.TransactionExtention{}
	if constant {
		result, err = c.triggerConstantContract(ct)

	} else {
		result, err = c.triggerContract(ct, feeLimit)
	}
	if err != nil {
		return nil, err
	}
	if result.Result.Code > 0 {
		return result, fmt.Errorf(string(result.Result.Message))
	}
	return result, nil

}

// TRC20ContractBalance get Address balance
func (c *TronCli) TRC20ContractBalance(addr, contractAddress string) (*big.Int, error) {
	addrB, err := address.Base58ToAddress(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address %s: %v", addr, addr)
	}
	req := trc20BalanceOf + "0000000000000000000000000000000000000000000000000000000000000000"[len(addrB.Hex())-2:] + addrB.Hex()[2:]
	result, err := c.TRC20Call("", contractAddress, req, true, 0)
	if err != nil {
		return nil, err
	}
	data := common.ToHex(result.GetConstantResult()[0])
	r, err := c.ParseTRC20NumericProperty(data)
	if err != nil {
		return nil, fmt.Errorf("contract address %s: %v", contractAddress, err)
	}
	if r == nil {
		return nil, fmt.Errorf("contract address %s: invalid balance of %s", contractAddress, addr)
	}
	return r, nil
}

// TRC20Send send toke to address
func (c *TronCli) TRC20Send(from, to, contract string, amount *big.Int, feeLimit int64) (*api.TransactionExtention, error) {
	addrB, err := address.Base58ToAddress(to)
	if err != nil {
		return nil, err
	}
	ab := common.LeftPadBytes(amount.Bytes(), 32)
	req := trc20TransferMethodSignature + "0000000000000000000000000000000000000000000000000000000000000000"[len(addrB.Hex())-4:] + addrB.Hex()[4:]
	req += common.Bytes2Hex(ab)
	return c.TRC20Call(from, contract, req, false, feeLimit)
}

// ParseTRC20NumericProperty get number from data
func (c *TronCli) ParseTRC20NumericProperty(data string) (*big.Int, error) {
	if common.Has0xPrefix(data) {
		data = data[2:]
	}
	if len(data) == 64 {
		var n big.Int
		_, ok := n.SetString(data, 16)
		if ok {
			return &n, nil
		}
	}
	return nil, fmt.Errorf("Cannot parse %s", data)
}

// ParseTRC20StringProperty get string from data
func (g *TronCli) ParseTRC20StringProperty(data string) (string, error) {
	if common.Has0xPrefix(data) {
		data = data[2:]
	}
	if len(data) > 128 {
		n, _ := g.ParseTRC20NumericProperty(data[64:128])
		if n != nil {
			l := n.Uint64()
			if 2*int(l) <= len(data)-128 {
				b, err := hex.DecodeString(data[128 : 128+2*l])
				if err == nil {
					return string(b), nil
				}
			}
		}
	} else if len(data) == 64 {
		// allow string properties as 32 bytes of UTF-8 data
		b, err := hex.DecodeString(data)
		if err == nil {
			i := bytes.Index(b, []byte{0})
			if i > 0 {
				b = b[:i]
			}
			if utf8.Valid(b) {
				return string(b), nil
			}
		}
	}
	return "", fmt.Errorf("Cannot parse %s,", data)
}

// TRC20GetName get token name
func (c *TronCli) TRC20GetName(contractAddress string) (string, error) {
	result, err := c.TRC20Call("", contractAddress, trc20NameSignature, true, 0)
	if err != nil {
		return "", err
	}
	data := common.ToHex(result.GetConstantResult()[0])
	return c.ParseTRC20StringProperty(data)
}

// TRC20GetSymbol get contract symbol
func (c *TronCli) TRC20GetSymbol(contractAddress string) (string, error) {
	result, err := c.TRC20Call("", contractAddress, trc20SymbolSignature, true, 0)
	if err != nil {
		return "", err
	}
	data := common.ToHex(result.GetConstantResult()[0])
	return c.ParseTRC20StringProperty(data)
}

// TRC20GetDecimals get contract decimals
func (c *TronCli) TRC20GetDecimals(contractAddress string) (*big.Int, error) {
	result, err := c.TRC20Call("", contractAddress, trc20DecimalsSignature, true, 0)
	if err != nil {
		return nil, err
	}
	data := common.ToHex(result.GetConstantResult()[0])
	return c.ParseTRC20NumericProperty(data)
}

// TriggerConstantContract and return tx result
func (c *TronCli) TriggerConstantContract(from, contractAddress, method, jsonString string) (*api.TransactionExtention, error) {
	var err error
	fromDesc := address.HexToAddress("410000000000000000000000000000000000000000")
	if len(from) > 0 {
		fromDesc, err = address.Base58ToAddress(from)
		if err != nil {
			return nil, err
		}
	}
	contractDesc, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}

	param, err := abi.LoadFromJSON(jsonString)
	if err != nil {
		return nil, err
	}

	dataBytes, err := abi.Pack(method, param)
	if err != nil {
		return nil, err
	}

	ct := &core.TriggerSmartContract{
		OwnerAddress:    fromDesc.Bytes(),
		ContractAddress: contractDesc.Bytes(),
		Data:            dataBytes,
	}

	return c.triggerConstantContract(ct)
}

// triggerConstantContract and return tx result
func (c *TronCli) triggerConstantContract(ct *core.TriggerSmartContract) (*api.TransactionExtention, error) {
	ctx, cancel := c.getContext()
	defer cancel()

	return c.TronWalletCli.TriggerConstantContract(ctx, ct)
}

// TriggerContract and return tx result
func (c *TronCli) TriggerContract(from, contractAddress, method, jsonString string,
	feeLimit, tAmount int64, tTokenID string, tTokenAmount int64) (*api.TransactionExtention, error) {
	fromDesc, err := address.Base58ToAddress(from)
	if err != nil {
		return nil, err
	}

	contractDesc, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}

	param, err := abi.LoadFromJSON(jsonString)
	if err != nil {
		return nil, err
	}

	dataBytes, err := abi.Pack(method, param)
	if err != nil {
		return nil, err
	}

	ct := &core.TriggerSmartContract{
		OwnerAddress:    fromDesc.Bytes(),
		ContractAddress: contractDesc.Bytes(),
		Data:            dataBytes,
	}
	if tAmount > 0 {
		ct.CallValue = tAmount
	}
	if len(tTokenID) > 0 && tTokenAmount > 0 {
		ct.CallTokenValue = tTokenAmount
		ct.TokenId, err = strconv.ParseInt(tTokenID, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return c.triggerContract(ct, feeLimit)
}

// triggerContract and return tx result
func (c *TronCli) triggerContract(ct *core.TriggerSmartContract, feeLimit int64) (*api.TransactionExtention, error) {
	ctx, cancel := c.getContext()
	defer cancel()

	tx, err := c.TronWalletCli.TriggerConstantContract(ctx, ct)
	if err != nil {
		return nil, err
	}

	if tx.Result.Code > 0 {
		return nil, fmt.Errorf("%s", string(tx.Result.Message))
	}
	if feeLimit > 0 {
		tx.Transaction.RawData.FeeLimit = feeLimit
		// update hash
		c.UpdateHash(tx)
	}
	return tx, err
}

// UpdateHash after local changes
func (c *TronCli) UpdateHash(tx *api.TransactionExtention) error {
	rawData, err := proto.Marshal(tx.Transaction.GetRawData())
	if err != nil {
		return err
	}

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	tx.Txid = hash
	return nil
}

func (c *TronCli) getContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), c.GrpcTimeout)
	if len(c.ApiKey) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, "TRON-PRO-API-KEY", c.ApiKey)
	}
	return ctx, cancel
}
