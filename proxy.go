package etherscan

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

type ProxyClient struct {
	api *apiClient
}

const proxyModule = "proxy"

var errNotImplemented = errors.New("not implemented")

func (c *ProxyClient) BlockNumber(ctx context.Context) (uint64, error) {
	var result hexutil.Uint64
	err := c.api.call(ctx, &callParams{
		module: proxyModule,
		action: "eth_blockNumber",
		result: &result,
	})

	return uint64(result), err
}

type ProxyBaseBlockInfo struct {
	Difficulty       *big.Int `etherscan:"difficulty,hex"`
	ExtraData        []byte   `etherscan:"extraData,hex"`
	GasLimit         *big.Int `etherscan:"gasLimit,hex"`
	GasUsed          *big.Int `etherscan:"gasUsed,hex"`
	Hash             common.Hash
	LogsBloom        []byte `etherscan:"logsBloom,hex"`
	Miner            common.Address
	MixHash          common.Hash `etherscan:"mixHash"`
	Nonce            *big.Int    `etherscan:"nonce,hex"`
	Number           uint64      `etherscan:"number,hex"`
	ParentHash       common.Hash `etherscan:"parentHash"`
	ReceiptsRoot     common.Hash `etherscan:"receiptsRoot"`
	SHA3Uncles       common.Hash `etherscan:"sha3Uncles"`
	Size             uint64      `etherscan:"size,hex"`
	StateRoot        common.Hash `etherscan:"stateRoot"`
	Timestamp        time.Time   `etherscan:"timestamp,hex"`
	TransactionsRoot common.Hash `etherscan:"transactionsRoot"`
	Uncles           []common.Hash
}

type ProxyFullBlockInfo struct {
	ProxyBaseBlockInfo
	TotalDifficulty *big.Int `etherscan:"totalDifficulty,hex"`
	Transactions    []ProxyTransactionInfo
}

type getBlockByNumRequest struct {
	Number  uint64 `etherscan:"tag,hex"`
	Boolean bool
}

func (c *ProxyClient) GetBlockByNumberFull(
	ctx context.Context, number uint64,
) (*ProxyFullBlockInfo, error) {
	req := getBlockByNumRequest{Number: number, Boolean: true}
	result := new(ProxyFullBlockInfo)

	err := c.api.call(ctx, &callParams{
		module:  proxyModule,
		action:  "eth_getBlockByNumber",
		request: req,
		result:  result,
	})

	return result, err
}

type ProxySummaryBlockInfo struct {
	ProxyBaseBlockInfo
	TotalDifficulty *big.Int `etherscan:"totalDifficulty,hex"`
	Transactions    []common.Hash
}

func (c *ProxyClient) GetBlockByNumberSummary(
	ctx context.Context, number uint64,
) (*ProxySummaryBlockInfo, error) {
	req := getBlockByNumRequest{Number: number, Boolean: false}
	result := new(ProxySummaryBlockInfo)

	err := c.api.call(ctx, &callParams{
		module:  proxyModule,
		action:  "eth_getBlockByNumber",
		request: req,
		result:  result,
	})

	return result, err
}

type BlockNumberAndIndex struct {
	Number uint64 `etherscan:"tag,hex"`
	Index  uint32 `etherscan:"index,hex"`
}

type ProxyUncleBlockInfo struct {
	ProxyBaseBlockInfo
	BaseFeePerGas *big.Int `etherscan:"baseFeePerGas,hex"`
}

func (c *ProxyClient) GetUncleByBlockNumberAndIndex(
	ctx context.Context, req *BlockNumberAndIndex,
) (*ProxyUncleBlockInfo, error) {
	result := new(ProxyUncleBlockInfo)

	err := c.api.call(ctx, &callParams{
		module:  proxyModule,
		action:  "eth_getUncleByBlockNumberAndIndex",
		request: req,
		result:  result,
	})

	return result, err
}

type blockTxCountRequest struct {
	Tag uint64 `etherscan:"tag,hex"`
}

func (c *ProxyClient) GetBlockTransactionCountByNumber(
	ctx context.Context, number uint64,
) (uint32, error) {
	req := blockTxCountRequest{number}
	var result hexutil.Uint

	err := c.api.call(ctx, &callParams{
		module:  proxyModule,
		action:  "eth_getBlockTransactionCountByNumber",
		request: req,
		result:  &result,
	})

	return uint32(result), err
}

type ProxyTransactionInfo struct {
	BlockHash        common.Hash `etherscan:"blockHash"`
	BlockNumber      uint64      `etherscan:"blockNumber,hex"`
	From             common.Address
	Gas              *big.Int `etherscan:"gas,hex"`
	GasPrice         *big.Int `etherscan:"gasPrice,hex"`
	Hash             common.Hash
	Input            []byte
	Nonce            uint64 `etherscan:"nonce,hex"`
	To               common.Address
	TransactionIndex uint64   `etherscan:"transactionIndex,hex"`
	Value            *big.Int `etherscan:"value,hex"`
	Type             uint32   `etherscan:"type,hex"`
	V                uint32   `etherscan:"v,hex"`
	R                *big.Int `etherscan:"r,hex"`
	S                *big.Int `etherscan:"s,hex"`
}

func (c *ProxyClient) GetTransactionByHash(
	ctx context.Context, txHash common.Hash,
) (*ProxyTransactionInfo, error) {
	req := struct{ TxHash common.Hash }{txHash}
	result := new(ProxyTransactionInfo)

	err := c.api.call(ctx, &callParams{
		module:  proxyModule,
		action:  "eth_getTransactionByHash",
		request: req,
		result:  result,
	})

	return result, err
}

func (c *ProxyClient) GetTransactionByBlockNumberAndIndex(
	ctx context.Context, req *BlockNumberAndIndex,
) (*ProxyTransactionInfo, error) {
	result := new(ProxyTransactionInfo)
	err := c.api.call(ctx, &callParams{
		module:  proxyModule,
		action:  "eth_getTransactionByBlockNumberAndIndex",
		request: req,
		result:  result,
	})

	return result, err
}

type TxCountRequest struct {
	Address common.Address
	Tag     BlockParameter
}

func (c *ProxyClient) GetTransactionCount(
	ctx context.Context, req *TxCountRequest,
) (uint64, error) {
	var result hexutil.Uint64
	err := c.api.call(ctx, &callParams{
		module:  proxyModule,
		action:  "eth_getTransactionCount",
		request: req,
		result:  &result,
	})

	return uint64(result), err
}

func (c *ProxyClient) SendRawTransaction(
	ctx context.Context, data []byte,
) (result common.Hash, err error) {
	req := struct{ Hex []byte }{data}
	err = c.api.call(ctx, &callParams{
		module:  proxyModule,
		action:  "eth_sendRawTransaction",
		request: req,
		result:  &result,
	})

	return result, err
}

type ProxyTransactionReceipt struct {
	BlockHash         common.Hash     `etherscan:"blockHash"`
	BlockNumber       uint64          `etherscan:"blockNumber,hex"`
	ContractAddress   *common.Address `etherscan:"contractAddress"`
	CumulativeGasUsed *big.Int        `etherscan:"cumulativeGasUsed,hex"`
	EffectiveGasPrice *big.Int        `etherscan:"effectiveGasPrice,hex"`
	From              common.Address
	GasUsed           *big.Int `etherscan:"gasUsed,hex"`
	Logs              []ProxyTxLog
	LogsBloom         []byte `etherscan:"logsBloom"`
	Status            bool   `etherscan:"status,hex"`
	To                common.Address
	TransactionHash   common.Hash `etherscan:"transactionHash"`
	TransactionIndex  uint32      `etherscan:"transactionIndex,hex"`
	Type              uint32      `etherscan:"type,hex"`
}

type ProxyTxLog struct {
	Address             common.Address
	BlockHash           common.Hash `etherscan:"blockHash"`
	BlockNumber         uint64      `etherscan:"blockNumber,hex"`
	Data                []byte
	LogIndex            uint32 `etherscan:"logIndex,hex"`
	Removed             bool
	Topics              []common.Hash
	TransactionHash     common.Hash `etherscan:"transactionHash"`
	TransactionIndex    uint32      `etherscan:"transactionIndex,hex"`
	TransactionLogIndex uint32      `etherscan:"transactionLogIndex,hex"`
	Type                string
}

type proxyTxLogResult struct {
	Address             common.Address `json:"address"`
	BlockHash           common.Hash    `json:"blockHash"`
	BlockNumber         hexutil.Uint64 `json:"blockNumber"`
	Data                hexutil.Bytes  `json:"data"`
	LogIndex            hexutil.Uint   `json:"logIndex"`
	Removed             bool           `json:"removed"`
	Topics              []common.Hash  `json:"topics"`
	TransactionHash     common.Hash    `json:"transactionHash"`
	TransactionIndex    hexutil.Uint   `json:"transactionIndex"`
	TransactionLogIndex hexutil.Uint   `json:"transactionLogIndex"`
	Type                string         `json:"type"`
}

func (c *ProxyClient) GetTransactionReceipt(
	ctx context.Context, txHash common.Hash,
) (*ProxyTransactionReceipt, error) {
	req := struct{ TxHash common.Hash }{txHash}
	result := new(ProxyTransactionReceipt)

	err := c.api.call(ctx, &callParams{
		module:  proxyModule,
		action:  "eth_getTransactionReceipt",
		request: req,
		result:  result,
	})

	return result, err
}

type CallRequest struct {
	To   common.Address
	Data []byte
	Tag  BlockParameter
}

func (c *ProxyClient) Call(
	ctx context.Context, req *CallRequest,
) ([]byte, error) {
	var result hexutil.Bytes

	err := c.api.call(ctx, &callParams{
		module:  proxyModule,
		action:  "eth_call",
		request: req,
		result:  &result,
	})

	return result, err
}

type GetCodeRequest struct {
	Address common.Address
	Tag     BlockParameter
}

func (c *ProxyClient) GetCode(
	ctx context.Context, req *GetCodeRequest,
) ([]byte, error) {
	var result hexutil.Bytes
	err := c.api.call(ctx, &callParams{
		module:  proxyModule,
		action:  "eth_getCode",
		request: req,
		result:  &result,
	})

	return result, err
}

type GetStorageRequest struct {
	Address  common.Address
	Position uint32 `etherscan:"position,hex"`
	Tag      BlockParameter
}

func (c *ProxyClient) GetStorageAt(
	ctx context.Context, req *GetStorageRequest,
) ([]byte, error) {
	var result hexutil.Bytes
	err := c.api.call(ctx, &callParams{
		module:  proxyModule,
		action:  "eth_getStorageAt",
		request: req,
		result:  &result,
	})

	return result, err
}

func (c *ProxyClient) GasPrice(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	err := c.api.call(ctx, &callParams{
		module: proxyModule,
		action: "eth_gasPrice",
		result: &result,
	})
	if err != nil {
		return nil, err
	}

	return result.ToInt(), nil
}

type EstimateGasRequest struct {
	Data     []byte
	To       common.Address
	Value    *big.Int `etherscan:"value,hex"`
	Gas      *big.Int `etherscan:"gas,hex"`
	GasPrice *big.Int `etherscan:"gasPrice,hex"`
}

func (c *ProxyClient) EstimateGas(
	ctx context.Context, req *EstimateGasRequest,
) (*big.Int, error) {
	var result hexutil.Big

	err := c.api.call(ctx, &callParams{
		module:  proxyModule,
		action:  "eth_estimateGas",
		request: req,
		result:  &result,
	})
	if err != nil {
		return nil, err
	}

	return result.ToInt(), nil
}
