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
	BaseFeePerGas    *big.Int `etherscan:"baseFeePerGas,hex"`
	Difficulty       *big.Int `etherscan:"difficulty,hex"`
	ExtraData        []byte   `etherscan:"extraData,hex"`
	GasLimit         *big.Int `etherscan:"gasLimit,hex"`
	GasUsed          *big.Int `etherscan:"gasUsed,hex"`
	Hash             common.Hash
	LogsBloom        []byte `etherscan:"logsBloom,hex"`
	Miner            common.Address
	MixHash          common.Hash
	Nonce            *big.Int    `etherscan:"nonce,hex"`
	Number           uint64      `etherscan:"number,hex"`
	ParentHash       common.Hash `etherscan:"parentHash"`
	ReceiptsRoot     common.Hash `etherscan:"receiptsRoot"`
	SHA3Uncles       common.Hash `etherscan:"sha3Uncles"`
	Size             uint64      `etherscan:"size,hex"`
	StateRoot        common.Hash `etherscan:"stateRoot"`
	Timestamp        time.Time   `etherscan:"timestamp,hex"`
	TotalDifficulty  *big.Int    `etherscan:"totalDifficulty,hex"`
	TransactionsRoot common.Hash `etherscan:"transactionsRoot"`
	Uncles           []common.Hash
}

type ProxyFullBlockInfo struct {
	ProxyBaseBlockInfo
	Transactions []ProxyTransactionInfo
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

	result := new(ProxyFullBlockInfo)
	if err := unmarshalResponse(rspData, result); err != nil {
		return nil, err
	}

	return result, nil
}

type ProxySummaryBlockInfo struct {
	ProxyBaseBlockInfo
	Transactions []common.Hash
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

	result := new(ProxySummaryBlockInfo)
	if err := unmarshalResponse(rspData, result); err != nil {
		return nil, err
	}

	return result, nil
}

type BlockNumberAndIndex struct {
	Number uint64 `etherscan:"tag,hex"`
	Index  uint32 `etherscan:"index,hex"`
}

func (c *ProxyClient) GetUncleByBlockNumberAndIndex(
	ctx context.Context, req *BlockNumberAndIndex,
) (*ProxyBaseBlockInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getUncleByBlockNumberAndIndex",
		other:  marshalRequest(req),
	})
	if err != nil {
		return nil, err
	}

	result := new(ProxyBaseBlockInfo)
	if err := unmarshalResponse(rspData, result); err != nil {
		return nil, err
	}

	return result, nil
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
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getTransactionByHash",
		other:  map[string]string{"txhash": txHash.String()},
	})
	if err != nil {
		return nil, err
	}

	result := new(ProxyTransactionInfo)
	if err := unmarshalResponse(rspData, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *ProxyClient) GetTransactionByBlockNumberAndIndex(
	ctx context.Context, req *BlockNumberAndIndex,
) (*ProxyTransactionInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getTransactionByBlockNumberAndIndex",
		other:  marshalRequest(req),
	})
	if err != nil {
		return nil, err
	}

	result := new(ProxyTransactionInfo)
	if err := unmarshalResponse(rspData, result); err != nil {
		return nil, err
	}

	return result, nil
}

type TxCountRequest struct {
	Address common.Address
	Tag     BlockParameter
}

func (c *ProxyClient) GetTransactionCount(
	ctx context.Context, req *TxCountRequest,
) (uint64, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getTransactionCount",
		other:  marshalRequest(req),
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
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getTransactionReceipt",
		other:  map[string]string{"txhash": txHash.String()},
	})
	if err != nil {
		return nil, err
	}

	result := new(ProxyTransactionReceipt)
	if err := unmarshalResponse(rspData, result); err != nil {
		return nil, err
	}

	return result, nil
}

type CallRequest struct {
	To   common.Address
	Data []byte
	Tag  BlockParameter
}

func (c *ProxyClient) Call(
	ctx context.Context, req *CallRequest,
) ([]byte, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_call",
		other:  marshalRequest(req),
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

func (c *ProxyClient) GetCode(
	ctx context.Context, req *GetCodeRequest,
) ([]byte, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getCode",
		other:  marshalRequest(req),
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
	Position uint32 `etherscan:"position,hex"`
	Tag      BlockParameter
}

func (c *ProxyClient) GetStorageAt(
	ctx context.Context, req *GetStorageRequest,
) ([]byte, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_getStorageAt",
		other:  marshalRequest(req),
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
	Value    *big.Int `etherscan:"value,hex"`
	Gas      *big.Int `etherscan:"gas,hex"`
	GasPrice *big.Int `etherscan:"gasPrice,hex"`
}

func (c *ProxyClient) EstimateGas(
	ctx context.Context, req *EstimateGasRequest,
) (*big.Int, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: proxyModule,
		action: "eth_estimateGas",
		other:  marshalRequest(req),
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
