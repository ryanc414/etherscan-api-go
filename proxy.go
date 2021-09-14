package etherscan

import (
	"context"
	"encoding/json"
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
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_blockNumber",
	})
	if err != nil {
		return 0, errors.Wrap(err, "while getting block number")
	}

	var result hexutil.Uint64
	if err := json.Unmarshal(rspData, &result); err != nil {
		return 0, err
	}

	return uint64(result), nil
}

type ProxyBaseBlockInfo struct {
	BaseFeePerGas    *big.Int
	Difficulty       *big.Int
	ExtraData        []byte
	GasLimit         *big.Int
	GasUsed          *big.Int
	Hash             common.Hash
	LogsBloom        []byte
	Miner            common.Address
	MixHash          common.Hash
	Nonce            *big.Int
	Number           uint64
	ParentHash       common.Hash
	ReceiptsRoot     common.Hash
	SHA3Uncles       common.Hash
	Size             uint64
	StateRoot        common.Hash
	Timestamp        time.Time
	TotalDifficulty  *big.Int
	TransactionsRoot common.Hash
	Uncles           []common.Hash
}

type proxyBaseBlockResult struct {
	BaseFeePerGas    *hexutil.Big   `json:"baseFeePerGas"`
	Difficulty       *hexutil.Big   `json:"difficulty"`
	ExtraData        hexutil.Bytes  `json:"extraData"`
	GasLimit         *hexutil.Big   `json:"gasLimit"`
	GasUsed          *hexutil.Big   `json:"gasUsed"`
	Hash             common.Hash    `json:"hash"`
	LogsBloom        hexutil.Bytes  `json:"logsBloom"`
	Miner            common.Address `json:"miner"`
	MixHash          common.Hash    `json:"mixHash"`
	Nonce            *hexutil.Big   `json:"nonce"`
	Number           hexutil.Uint64 `json:"number"`
	ParentHash       common.Hash    `json:"parentHash"`
	ReceiptsRoot     common.Hash    `json:"receiptsRoot"`
	SHA3Uncles       common.Hash    `json:"sha3Uncles"`
	Size             hexutil.Uint64 `json:"size"`
	StateRoot        common.Hash    `json:"stateRoot"`
	Timestamp        hexTimestamp   `json:"timestamp"`
	TotalDifficulty  *hexutil.Big   `json:"totalDifficulty"`
	TransactionsRoot common.Hash    `json:"transactionsRoot"`
	Uncles           []common.Hash  `json:"uncles"`
}

func (res *proxyBaseBlockResult) toInfo() *ProxyBaseBlockInfo {
	return &ProxyBaseBlockInfo{
		BaseFeePerGas:    res.BaseFeePerGas.ToInt(),
		Difficulty:       res.Difficulty.ToInt(),
		ExtraData:        res.ExtraData,
		GasLimit:         res.GasLimit.ToInt(),
		GasUsed:          res.GasUsed.ToInt(),
		Hash:             res.Hash,
		LogsBloom:        res.LogsBloom,
		Miner:            res.Miner,
		MixHash:          res.MixHash,
		Nonce:            res.Nonce.ToInt(),
		Number:           uint64(res.Number),
		ParentHash:       res.ParentHash,
		ReceiptsRoot:     res.ReceiptsRoot,
		SHA3Uncles:       res.SHA3Uncles,
		Size:             uint64(res.Size),
		StateRoot:        res.StateRoot,
		Timestamp:        res.Timestamp.unwrap(),
		TotalDifficulty:  res.TotalDifficulty.ToInt(),
		TransactionsRoot: res.TransactionsRoot,
		Uncles:           res.Uncles,
	}
}

type ProxyFullBlockInfo struct {
	ProxyBaseBlockInfo
	Transactions []ProxyTransactionInfo
}

type proxyFullBlockResult struct {
	proxyBaseBlockResult
	Transactions []proxyTransactionResult `json:"transactions"`
}

func (res *proxyFullBlockResult) toInfo() *ProxyFullBlockInfo {
	transactions := make([]ProxyTransactionInfo, len(res.Transactions))
	for i := range res.Transactions {
		transactions[i] = *res.Transactions[i].toInfo()
	}

	return &ProxyFullBlockInfo{
		ProxyBaseBlockInfo: *res.proxyBaseBlockResult.toInfo(),
		Transactions:       transactions,
	}
}

func (c *ProxyClient) GetBlockByNumberFull(
	ctx context.Context, number uint64,
) (*ProxyFullBlockInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getBlockByNumber",
		other: map[string]string{
			"tag":     hexutil.EncodeUint64(number),
			"boolean": "true",
		},
	})
	if err != nil {
		return nil, err
	}

	var result proxyFullBlockResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result.toInfo(), nil
}

type ProxySummaryBlockInfo struct {
	ProxyBaseBlockInfo
	Transactions []common.Hash
}

type proxySummaryBlockResult struct {
	proxyBaseBlockResult
	Transactions []common.Hash `json:"transactions"`
}

func (res *proxySummaryBlockResult) toInfo() *ProxySummaryBlockInfo {
	return &ProxySummaryBlockInfo{
		ProxyBaseBlockInfo: *res.proxyBaseBlockResult.toInfo(),
		Transactions:       res.Transactions,
	}
}

func (c *ProxyClient) GetBlockByNumberSummary(
	ctx context.Context, number uint64,
) (*ProxySummaryBlockInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getBlockByNumber",
		other: map[string]string{
			"tag":     hexutil.EncodeUint64(number),
			"boolean": "false",
		},
	})
	if err != nil {
		return nil, err
	}

	var result proxySummaryBlockResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result.toInfo(), nil
}

type BlockNumberAndIndex struct {
	Number uint64
	Index  uint32
}

func (b BlockNumberAndIndex) toParams() map[string]string {
	return map[string]string{
		"tag":   hexutil.EncodeUint64(b.Number),
		"index": hexutil.EncodeUint64(uint64(b.Index)),
	}
}

func (c *ProxyClient) GetUncleByBlockNumberAndIndex(
	ctx context.Context, req *BlockNumberAndIndex,
) (*ProxyBaseBlockInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getUncleByBlockNumberAndIndex",
		other:  req.toParams(),
	})
	if err != nil {
		return nil, err
	}

	var result proxyBaseBlockResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result.toInfo(), nil
}

func (c *ProxyClient) GetBlockTransactionCountByNumber(
	ctx context.Context, number uint64,
) (uint32, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getBlockTransactionCountByNumber",
		other:  map[string]string{"tag": hexutil.EncodeUint64(number)},
	})
	if err != nil {
		return 0, err
	}

	var result hexutil.Uint
	if err := json.Unmarshal(rspData, &result); err != nil {
		return 0, err
	}

	return uint32(result), nil
}

type ProxyTransactionInfo struct {
	BlockHash        common.Hash
	BlockNumber      uint64
	From             common.Address
	Gas              *big.Int
	GasPrice         *big.Int
	Hash             common.Hash
	Input            []byte
	Nonce            uint64
	To               common.Address
	TransactionIndex uint64
	Value            *big.Int
	Type             uint32
	V                uint32
	R                *big.Int
	S                *big.Int
}

type proxyTransactionResult struct {
	BlockHash        common.Hash    `json:"blockHash"`
	BlockNumber      hexutil.Uint64 `json:"blockNumber"`
	From             common.Address `json:"from"`
	Gas              *hexutil.Big   `json:"gas"`
	GasPrice         *hexutil.Big   `json:"gasPrice"`
	Hash             common.Hash    `json:"hash"`
	Input            hexutil.Bytes  `json:"input"`
	Nonce            hexutil.Uint64 `json:"nonce"`
	To               common.Address `json:"to"`
	TransactionIndex hexutil.Uint64 `json:"transactionIndex"`
	Value            *hexutil.Big   `json:"value"`
	Type             hexutil.Uint   `json:"type"`
	V                hexutil.Uint   `json:"v"`
	R                *hexutil.Big   `json:"r"`
	S                *hexutil.Big   `json:"s"`
}

func (res *proxyTransactionResult) toInfo() *ProxyTransactionInfo {
	return &ProxyTransactionInfo{
		BlockHash:        res.BlockHash,
		BlockNumber:      uint64(res.BlockNumber),
		From:             res.From,
		Gas:              res.Gas.ToInt(),
		GasPrice:         res.GasPrice.ToInt(),
		Hash:             res.Hash,
		Input:            res.Input,
		Nonce:            uint64(res.Nonce),
		To:               res.To,
		TransactionIndex: uint64(res.TransactionIndex),
		Value:            res.Value.ToInt(),
		Type:             uint32(res.Type),
		V:                uint32(res.V),
		R:                res.R.ToInt(),
		S:                res.S.ToInt(),
	}
}

func (c *ProxyClient) GetTransactionByHash(
	ctx context.Context, txHash common.Hash,
) (*ProxyTransactionInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getTransactionByHash",
		other:  map[string]string{"txhash": txHash.String()},
	})
	if err != nil {
		return nil, err
	}

	var result proxyTransactionResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result.toInfo(), nil
}

func (c *ProxyClient) GetTransactionByBlockNumberAndIndex(
	ctx context.Context, req *BlockNumberAndIndex,
) (*ProxyTransactionInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getTransactionByBlockNumberAndIndex",
		other:  req.toParams(),
	})
	if err != nil {
		return nil, err
	}

	var result proxyTransactionResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result.toInfo(), nil
}

type TxCountRequest struct {
	Address common.Address
	Tag     BlockParameter
}

func (req *TxCountRequest) toParams() map[string]string {
	return map[string]string{
		"address": req.Address.String(),
		"tag":     req.Tag.String(),
	}
}

func (c *ProxyClient) GetTransactionCount(
	ctx context.Context, req *TxCountRequest,
) (uint64, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getTransactionCount",
		other:  req.toParams(),
	})
	if err != nil {
		return 0, err
	}

	var result hexutil.Uint64
	if err := json.Unmarshal(rspData, &result); err != nil {
		return 0, err
	}

	return uint64(result), nil
}

func (c *ProxyClient) SendRawTransaction(
	ctx context.Context, data []byte,
) (common.Hash, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_sendRawTransaction",
		other:  map[string]string{"hex": hexutil.Encode(data)},
	})
	if err != nil {
		return common.Hash{}, err
	}

	var result common.Hash
	if err := json.Unmarshal(rspData, &result); err != nil {
		return common.Hash{}, err
	}

	return result, nil
}

type ProxyTransactionReceipt struct {
	BlockHash         common.Hash
	BlockNumber       uint64
	ContractAddress   *common.Address
	CumulativeGasUsed *big.Int
	EffectiveGasPrice *big.Int
	From              common.Address
	GasUsed           *big.Int
	Logs              []ProxyTxLog
	LogsBloom         []byte
	Status            bool
	To                common.Address
	TransactionHash   common.Hash
	TransactionIndex  uint32
	Type              uint32
}

type ProxyTxLog struct {
	Address             common.Address
	BlockHash           common.Hash
	BlockNumber         uint64
	Data                []byte
	LogIndex            uint32
	Removed             bool
	Topics              []common.Hash
	TransactionHash     common.Hash
	TransactionIndex    uint32
	TransactionLogIndex uint32
	Type                string
}

type proxyTxReceiptResult struct {
	BlockHash         common.Hash        `json:"blockHash"`
	BlockNumber       hexutil.Uint64     `json:"blockNumber"`
	ContractAddress   *common.Address    `json:"contractAddress"`
	CumulativeGasUsed *hexutil.Big       `json:"cumulativeGasUsed"`
	EffectiveGasPrice *hexutil.Big       `json:"effectiveGasPrice"`
	From              common.Address     `json:"from"`
	GasUsed           *hexutil.Big       `json:"gasUsed"`
	Logs              []proxyTxLogResult `json:"logs"`
	LogsBloom         hexutil.Bytes      `json:"logsBloom"`
	Status            hexutil.Uint       `json:"status"`
	To                common.Address     `json:"to"`
	TransactionHash   common.Hash        `json:"transactionHash"`
	TransactionIndex  hexutil.Uint       `json:"transactionIndex"`
	Type              hexutil.Uint       `json:"type"`
}

func (res *proxyTxReceiptResult) toReceipt() *ProxyTransactionReceipt {
	logs := make([]ProxyTxLog, len(res.Logs))
	for i := range res.Logs {
		logs[i] = *res.Logs[i].toLog()
	}

	return &ProxyTransactionReceipt{
		BlockHash:         res.BlockHash,
		BlockNumber:       uint64(res.BlockNumber),
		ContractAddress:   res.ContractAddress,
		CumulativeGasUsed: res.CumulativeGasUsed.ToInt(),
		EffectiveGasPrice: res.EffectiveGasPrice.ToInt(),
		From:              res.From,
		GasUsed:           res.GasUsed.ToInt(),
		Logs:              logs,
		LogsBloom:         res.LogsBloom,
		Status:            res.Status != 0,
		To:                res.To,
		TransactionHash:   res.TransactionHash,
		TransactionIndex:  uint32(res.TransactionIndex),
		Type:              uint32(res.Type),
	}
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

func (log *proxyTxLogResult) toLog() *ProxyTxLog {
	return &ProxyTxLog{
		Address:             log.Address,
		BlockHash:           log.BlockHash,
		BlockNumber:         uint64(log.BlockNumber),
		Data:                log.Data,
		LogIndex:            uint32(log.LogIndex),
		Removed:             log.Removed,
		Topics:              log.Topics,
		TransactionHash:     log.TransactionHash,
		TransactionIndex:    uint32(log.TransactionIndex),
		TransactionLogIndex: uint32(log.TransactionLogIndex),
		Type:                log.Type,
	}
}

func (c *ProxyClient) GetTransactionReceipt(
	ctx context.Context, txHash common.Hash,
) (*ProxyTransactionReceipt, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getTransactionReceipt",
		other:  map[string]string{"txhash": txHash.String()},
	})
	if err != nil {
		return nil, err
	}

	var result proxyTxReceiptResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result.toReceipt(), nil
}

type CallRequest struct {
	To   common.Address
	Data []byte
	Tag  BlockParameter
}

func (req *CallRequest) toParams() map[string]string {
	return map[string]string{
		"to":   req.To.String(),
		"data": hexutil.Encode(req.Data),
		"tag":  req.Tag.String(),
	}
}

func (c *ProxyClient) Call(
	ctx context.Context, req *CallRequest,
) ([]byte, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_call",
		other:  req.toParams(),
	})
	if err != nil {
		return nil, err
	}

	var result hexutil.Bytes
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

type GetCodeRequest struct {
	Address common.Address
	Tag     BlockParameter
}

func (req *GetCodeRequest) toParams() map[string]string {
	return map[string]string{
		"address": req.Address.String(),
		"tag":     req.Tag.String(),
	}
}

func (c *ProxyClient) GetCode(
	ctx context.Context, req *GetCodeRequest,
) ([]byte, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getCode",
		other:  req.toParams(),
	})
	if err != nil {
		return nil, err
	}

	var result hexutil.Bytes
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

type GetStorageRequest struct {
	Address  common.Address
	Position uint32
	Tag      BlockParameter
}

func (req *GetStorageRequest) toParams() map[string]string {
	return map[string]string{
		"address":  req.Address.String(),
		"position": hexutil.EncodeUint64(uint64(req.Position)),
		"tag":      req.Tag.String(),
	}
}

func (c *ProxyClient) GetStorageAt(
	ctx context.Context, req *GetStorageRequest,
) ([]byte, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getStorageAt",
		other:  req.toParams(),
	})
	if err != nil {
		return nil, err
	}

	var result hexutil.Bytes
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *ProxyClient) GasPrice(ctx context.Context) (*big.Int, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_gasPrice",
	})
	if err != nil {
		return nil, err
	}

	var result hexutil.Big
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result.ToInt(), nil
}

type EstimateGasRequest struct {
	Data     []byte
	To       common.Address
	Value    *big.Int
	Gas      *big.Int
	GasPrice *big.Int
}

func (req *EstimateGasRequest) toParams() map[string]string {
	return map[string]string{
		"data":     hexutil.Encode(req.Data),
		"to":       req.To.String(),
		"value":    hexutil.EncodeBig(req.Value),
		"gas":      hexutil.EncodeBig(req.Gas),
		"gasPrice": hexutil.EncodeBig(req.GasPrice),
	}
}

func (c *ProxyClient) EstimateGas(
	ctx context.Context, req *EstimateGasRequest,
) (*big.Int, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_estimateGas",
		other:  req.toParams(),
	})
	if err != nil {
		return nil, err
	}

	var result hexutil.Big
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result.ToInt(), nil
}
