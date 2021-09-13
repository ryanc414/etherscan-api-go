package etherscan

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type ProxyClient struct {
	api *apiClient
}

const proxyModule = "proxy"

var errNotImplemented = errors.New("not implemented")

func (c *ProxyClient) BlockNumber(ctx context.Context) (uint64, error) {
	return 0, errNotImplemented
}

type ProxyBaseBlockInfo struct {
	BaseFeePerGas     *big.Int
	Difficulty        *big.Int
	ExtraData         []byte
	GasLimit          *big.Int
	GasUsed           *big.Int
	Hash              common.Hash
	LogsBloom         []byte
	Miner             common.Address
	MixHash           common.Hash
	Nonce             *big.Int
	Number            uint64
	ParentHash        common.Hash
	ReceiptsRoot      common.Hash
	SHA3Uncles        common.Hash
	Size              uint64
	StateRoot         common.Hash
	Timestamp         time.Time
	TotalDifficulty   *big.Int
	TransactsionsRoot common.Hash
	Uncles            []common.Hash
}

type ProxyFullBlockInfo struct {
	ProxyBaseBlockInfo
	Transactions []ProxyTransactionInfo
}

func (c *ProxyClient) GetBlockByNumberFull(
	ctx context.Context, number uint64,
) (*ProxyFullBlockInfo, error) {
	return nil, errNotImplemented
}

type ProxySummaryBlockInfo struct {
	ProxyBaseBlockInfo
	Transactions []common.Hash
}

func (c *ProxyClient) GetBlockByNumberSummary(
	ctx context.Context, number uint64,
) (*ProxySummaryBlockInfo, error) {
	return nil, errNotImplemented
}

type BlockNumberAndIndex struct {
	Number uint64
	Index  uint32
}

func (c *ProxyClient) GetUncleByBlockNumberAndIndex(
	ctx context.Context, req *BlockNumberAndIndex,
) (*ProxyBaseBlockInfo, error) {
	return nil, errNotImplemented
}

func (c *ProxyClient) GetBlockTransactionCountByNumber(
	ctx context.Context, number uint64,
) (uint32, error) {
	return 0, errNotImplemented
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
	Type             []byte
	V                uint32
	R                *big.Int
	S                *big.Int
}

func (c *ProxyClient) GetTransactionByHash(
	ctx context.Context, txhash common.Hash,
) (*ProxyTransactionInfo, error) {
	return nil, errNotImplemented
}

func (c *ProxyClient) GetTransactionByBlockNumberAndIndex(
	ctx context.Context, req *BlockNumberAndIndex,
) (*ProxyTransactionInfo, error) {
	return nil, errNotImplemented
}

type TxCountRequest struct {
	Address common.Address
	Tag     BlockParameter
}

func (c *ProxyClient) GetTransactionCount(
	ctx context.Context, req *TxCountRequest,
) (uint64, error) {
	return 0, errNotImplemented
}

func (c *ProxyClient) SendRawTransaction(
	ctx context.Context, data []byte,
) (common.Hash, error) {
	return common.Hash{}, errNotImplemented
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

func (c *ProxyClient) GetTransactionReceipt(
	ctx context.Context, txhash common.Hash,
) (*ProxyTransactionReceipt, error) {
	return nil, errNotImplemented
}

type CallRequest struct {
	To   common.Address
	Data common.Hash
	Tag  BlockParameter
}

func (c *ProxyClient) Call(
	ctx context.Context, req *CallRequest,
) ([]byte, error) {
	return nil, errNotImplemented
}

type GetCodeRequest struct {
	Address common.Address
	Tag     BlockParameter
}

func (c *ProxyClient) GetCode(
	ctx context.Context, req *GetCodeRequest,
) ([]byte, error) {
	return nil, errNotImplemented
}

type GetStorageRequest struct {
	Address  common.Address
	Position uint32
	Tag      BlockParameter
}

func (c *ProxyClient) GetStorageAt(
	ctx context.Context, req *GetStorageRequest,
) ([]byte, error) {
	return nil, errNotImplemented
}

func (c *ProxyClient) GasPrice(ctx context.Context) (*big.Int, error) {
	return nil, errNotImplemented
}

type EstimateGasRequest struct {
	Data     []byte
	To       common.Address
	Value    *big.Int
	Gas      *big.Int
	GasPrice *big.Int
}

func (c *ProxyClient) EstimateGas(
	ctx context.Context, req *EstimateGasRequest,
) (*big.Int, error) {
	return nil, errNotImplemented
}
