package proxy

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ecommon "github.com/ryanc414/etherscan-api-go/common"
	"github.com/ryanc414/etherscan-api-go/httpapi"
)

// ProxyClient is the client for ethereum proxy actions.
type ProxyClient struct {
	API *httpapi.APIClient
}

// BlockNumber returns the current block number.
func (c *ProxyClient) BlockNumber(ctx context.Context) (uint64, error) {
	var result hexutil.Uint64
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module: ecommon.ProxyModule,
		Action: "eth_blockNumber",
		Result: &result,
	})

	return uint64(result), err
}

// ProxyBaseBlockInfo contains common information on a block.
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

// ProxyFullBlockInfo contains the full information on a block and its
// transactions.
type ProxyFullBlockInfo struct {
	ProxyBaseBlockInfo
	TotalDifficulty *big.Int `etherscan:"totalDifficulty,hex"`
	Transactions    []ProxyTransactionInfo
}

type getBlockByNumRequest struct {
	Number  uint64 `etherscan:"tag,hex"`
	Boolean bool
}

// GetBlockByNumberFull returns full information about a block by block number.
func (c *ProxyClient) GetBlockByNumberFull(
	ctx context.Context, number uint64,
) (*ProxyFullBlockInfo, error) {
	req := getBlockByNumRequest{Number: number, Boolean: true}
	result := new(ProxyFullBlockInfo)

	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ProxyModule,
		Action:  "eth_getBlockByNumber",
		Request: req,
		Result:  result,
	})

	return result, err
}

// ProxySummaryBlockInfo contains summary information on a block, including
// a slice of transaction hashes.
type ProxySummaryBlockInfo struct {
	ProxyBaseBlockInfo
	TotalDifficulty *big.Int `etherscan:"totalDifficulty,hex"`
	Transactions    []common.Hash
}

// GetBlockByNumberSummary returns summary information about a block by block number.
func (c *ProxyClient) GetBlockByNumberSummary(
	ctx context.Context, number uint64,
) (*ProxySummaryBlockInfo, error) {
	req := getBlockByNumRequest{Number: number, Boolean: false}
	result := new(ProxySummaryBlockInfo)

	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ProxyModule,
		Action:  "eth_getBlockByNumber",
		Request: req,
		Result:  result,
	})

	return result, err
}

// BlockNumberAndIndex uniquely identifies a transaction's position in a block.
type BlockNumberAndIndex struct {
	Number uint64 `etherscan:"tag,hex"`
	Index  uint32 `etherscan:"index,hex"`
}

// ProxyUncleBlockInfo contains information about an uncle block.
type ProxyUncleBlockInfo struct {
	ProxyBaseBlockInfo
	BaseFeePerGas *big.Int `etherscan:"baseFeePerGas,hex"`
}

// GetUncleByBlockNumberAndIndex returns information about a uncle by block number.
func (c *ProxyClient) GetUncleByBlockNumberAndIndex(
	ctx context.Context, req *BlockNumberAndIndex,
) (*ProxyUncleBlockInfo, error) {
	result := new(ProxyUncleBlockInfo)

	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ProxyModule,
		Action:  "eth_getUncleByBlockNumberAndIndex",
		Request: req,
		Result:  result,
	})

	return result, err
}

type blockTxCountRequest struct {
	Tag uint64 `etherscan:"tag,hex"`
}

// GetBlockTransactionCountByNumber returns the number of transactions in a block.
func (c *ProxyClient) GetBlockTransactionCountByNumber(
	ctx context.Context, number uint64,
) (uint32, error) {
	req := blockTxCountRequest{number}
	var result hexutil.Uint

	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ProxyModule,
		Action:  "eth_getBlockTransactionCountByNumber",
		Request: req,
		Result:  &result,
	})

	return uint32(result), err
}

// ProxyTransactionInfo contains information about a transaction.
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

// GetTransactionsByHash returns the information about a transaction requested by transaction hash.
func (c *ProxyClient) GetTransactionByHash(
	ctx context.Context, txHash common.Hash,
) (*ProxyTransactionInfo, error) {
	req := struct{ TxHash common.Hash }{txHash}
	result := new(ProxyTransactionInfo)

	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ProxyModule,
		Action:  "eth_getTransactionByHash",
		Request: req,
		Result:  result,
	})

	return result, err
}

// GetTransactionByBlockNumberAndIndex returns information about a transaction
// by block number and transaction index position.
func (c *ProxyClient) GetTransactionByBlockNumberAndIndex(
	ctx context.Context, req *BlockNumberAndIndex,
) (*ProxyTransactionInfo, error) {
	result := new(ProxyTransactionInfo)
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ProxyModule,
		Action:  "eth_getTransactionByBlockNumberAndIndex",
		Request: req,
		Result:  result,
	})

	return result, err
}

// TxCountRequest contains request parameters for GetTransactionCount.
type TxCountRequest struct {
	Address common.Address
	Tag     ecommon.BlockParameter
}

// GetTransactionCount returns the number of transactions performed by an address.
func (c *ProxyClient) GetTransactionCount(
	ctx context.Context, req *TxCountRequest,
) (uint64, error) {
	var result hexutil.Uint64
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ProxyModule,
		Action:  "eth_getTransactionCount",
		Request: req,
		Result:  &result,
	})

	return uint64(result), err
}

// SendRawTransaction submits a pre-signed transaction for broadcast to the Ethereum network.
func (c *ProxyClient) SendRawTransaction(
	ctx context.Context, data []byte,
) (result common.Hash, err error) {
	req := struct{ Hex []byte }{data}
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ProxyModule,
		Action:  "eth_sendRawTransaction",
		Request: req,
		Result:  &result,
	})

	return result, err
}

// ProxyTransactionReceipt describes a transaction receipt.
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

// ProxyTxLog describes a transaction log.
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

// GetTransactionReceipt returns the receipt of a transaction by transaction hash.
func (c *ProxyClient) GetTransactionReceipt(
	ctx context.Context, txHash common.Hash,
) (*ProxyTransactionReceipt, error) {
	req := struct{ TxHash common.Hash }{txHash}
	result := new(ProxyTransactionReceipt)

	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ProxyModule,
		Action:  "eth_getTransactionReceipt",
		Request: req,
		Result:  result,
	})

	return result, err
}

// CallRequest contains the request parameters for Call.
type CallRequest struct {
	To   common.Address
	Data []byte
	Tag  ecommon.BlockParameter
}

// Call executes a new message call immediately without creating a transaction on the block chain.
func (c *ProxyClient) Call(
	ctx context.Context, req *CallRequest,
) ([]byte, error) {
	var result hexutil.Bytes

	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ProxyModule,
		Action:  "eth_call",
		Request: req,
		Result:  &result,
	})

	return result, err
}

// GetCodeRequest contains the request parameters for GetCode.
type GetCodeRequest struct {
	Address common.Address
	Tag     ecommon.BlockParameter
}

// GetCode returns code at a given address.
func (c *ProxyClient) GetCode(
	ctx context.Context, req *GetCodeRequest,
) ([]byte, error) {
	var result hexutil.Bytes
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ProxyModule,
		Action:  "eth_getCode",
		Request: req,
		Result:  &result,
	})

	return result, err
}

// GetStorageRequest contains the request parameters for GetStorageAt.
type GetStorageRequest struct {
	Address  common.Address
	Position uint32 `etherscan:"position,hex"`
	Tag      ecommon.BlockParameter
}

// GetStorageAt returns the value from a storage position at a given address.
func (c *ProxyClient) GetStorageAt(
	ctx context.Context, req *GetStorageRequest,
) ([]byte, error) {
	var result hexutil.Bytes
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ProxyModule,
		Action:  "eth_getStorageAt",
		Request: req,
		Result:  &result,
	})

	return result, err
}

// GasPrice returns the current price per gas in wei.
func (c *ProxyClient) GasPrice(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module: ecommon.ProxyModule,
		Action: "eth_gasPrice",
		Result: &result,
	})
	if err != nil {
		return nil, err
	}

	return result.ToInt(), nil
}

// EstimateGasRequest contains the request parameters for EstimateGas.
type EstimateGasRequest struct {
	Data     []byte
	To       common.Address
	Value    *big.Int `etherscan:"value,hex"`
	Gas      *big.Int `etherscan:"gas,hex"`
	GasPrice *big.Int `etherscan:"gasPrice,hex"`
}

// EstimateGas makes a call or transaction, which won't be added to the blockchain and returns the used gas.
func (c *ProxyClient) EstimateGas(
	ctx context.Context, req *EstimateGasRequest,
) (*big.Int, error) {
	var result hexutil.Big

	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.ProxyModule,
		Action:  "eth_estimateGas",
		Request: req,
		Result:  &result,
	})
	if err != nil {
		return nil, err
	}

	return result.ToInt(), nil
}
