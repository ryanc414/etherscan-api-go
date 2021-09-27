package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const accountModule = "account"

type AccountsClient struct {
	api *apiClient
}

type ETHBalanceRequest struct {
	Address common.Address
	Tag     BlockParameter
}

type BlockParameter int32

const (
	BlockParameterLatest = iota
	BlockParameterEarliest
	BlockParameterPending
)

func (b BlockParameter) String() string {
	switch b {
	case BlockParameterEarliest:
		return "earliest"

	case BlockParameterPending:
		return "pending"

	case BlockParameterLatest:
		return "latest"

	default:
		panic(fmt.Sprintf("unknown block parameter %d", int32(b)))
	}
}

func (c *AccountsClient) GetETHBalance(
	ctx context.Context, req *ETHBalanceRequest,
) (*big.Int, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "balance",
		other:  marshalRequest(req),
	})
	if err != nil {
		return nil, err
	}

	var result bigInt
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result.unwrap(), nil
}

type MultiETHBalancesRequest struct {
	Addresses []common.Address `etherscan:"address"`
	Tag       BlockParameter
}

type MultiBalanceResponse struct {
	Account common.Address
	Balance *big.Int
}

func (c *AccountsClient) GetMultiETHBalances(
	ctx context.Context, req *MultiETHBalancesRequest,
) ([]MultiBalanceResponse, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "balancemulti",
		other:  marshalRequest(req),
	})
	if err != nil {
		return nil, err
	}

	var response []MultiBalanceResponse
	if err := unmarshalResponse(rspData, &response); err != nil {
		return nil, err
	}

	return response, nil
}

type ListTxRequest struct {
	Address    common.Address
	StartBlock uint64
	EndBlock   uint64
	Sort       SortingPreference
}

type SortingPreference int32

const (
	SortingPreferenceAscending = iota
	SortingPreferenceDescending
)

func (s SortingPreference) String() string {
	switch s {
	case SortingPreferenceAscending:
		return "asc"

	case SortingPreferenceDescending:
		return "desc"

	default:
		panic(fmt.Sprintf("unknown sorting preference %d", int32(s)))
	}
}

type TransactionInfo struct {
	BlockNumber     uint64    `etherscan:"blockNumber"`
	Timestamp       time.Time `etherscan:"timeStamp"`
	Hash            common.Hash
	From            common.Address
	To              common.Address
	Value           *big.Int
	ContractAddress *common.Address `etherscan:"contractAddress"`
	Input           []byte
	Gas             uint64
	GasUsed         uint64 `etherscan:"gasUsed"`
	IsError         bool   `etherscan:"isError"`
}

type NormalTxInfo struct {
	TransactionInfo
	Nonce             uint64
	BlockHash         common.Hash `etherscan:"blockHash"`
	TransactionIndex  uint64      `etherscan:"transactionIndex"`
	GasPrice          *big.Int    `etherscan:"gasPrice"`
	TxReceiptStatus   string      `etherscan:"txreceipt_status"`
	CumulativeGasUsed uint64      `etherscan:"cumulativeGasUsed"`
	Confirmations     uint64
}

type InternalTxInfo struct {
	TransactionInfo
	TraceID string `etherscan:"traceId"`
	Type    string
}

func (c *AccountsClient) ListNormalTransactions(
	ctx context.Context, req *ListTxRequest,
) ([]NormalTxInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "txlist",
		other:  marshalRequest(req),
	})
	if err != nil {
		return nil, err
	}

	var result []NormalTxInfo
	if err := unmarshalResponse(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *AccountsClient) ListInternalTransactions(
	ctx context.Context, req *ListTxRequest,
) ([]InternalTxInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "txlistinternal",
		other:  marshalRequest(req),
	})
	if err != nil {
		return nil, err
	}

	var result []InternalTxInfo
	if err := unmarshalResponse(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *AccountsClient) GetInternalTxsByHash(
	ctx context.Context, hash common.Hash,
) ([]InternalTxInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "txlistinternal",
		other:  map[string]string{"txhash": hash.String()},
	})
	if err != nil {
		return nil, err
	}

	var result []InternalTxInfo
	if err := unmarshalResponse(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

type BlockRangeRequest struct {
	StartBlock uint64
	EndBlock   uint64
	Sort       SortingPreference
}

func (c *AccountsClient) GetInternalTxsByBlockRange(
	ctx context.Context,
	req *BlockRangeRequest,
) ([]InternalTxInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "txlistinternal",
		other:  marshalRequest(req),
	})
	if err != nil {
		return nil, err
	}

	var result []InternalTxInfo
	if err := unmarshalResponse(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

type TokenTransfersRequest struct {
	Address         common.Address
	ContractAddress common.Address
	Sort            SortingPreference
}

type TokenTransferInfo struct {
	NormalTxInfo
	TokenName    string `json:"tokenName"`
	TokenSymbol  string `json:"tokenSymbol"`
	TokenDecimal uint32 `json:"tokenDecimal"`
}

func (c *AccountsClient) ListTokenTransfers(
	ctx context.Context, req *TokenTransfersRequest,
) ([]TokenTransferInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: "account",
		action: "tokentx",
		other:  marshalRequest(req),
	})
	if err != nil {
		return nil, err
	}

	var result []TokenTransferInfo
	if err := unmarshalResponse(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

type ListNFTTransferRequest struct {
	Address         *common.Address
	ContractAddress *common.Address
	Sort            SortingPreference
}

type NFTTransferInfo struct {
	TokenTransferInfo
	TokenID string `etherscan:"tokenID"`
}

func (c *AccountsClient) ListNFTTransfers(
	ctx context.Context, req *ListNFTTransferRequest,
) ([]NFTTransferInfo, error) {
	if req.Address == nil && req.ContractAddress == nil {
		return nil, errors.New("at least one of Address or ContractAddress must be specifide")
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "tokennfttx",
		other:  marshalRequest(req),
	})
	if err != nil {
		return nil, err
	}

	var result []NFTTransferInfo
	if err := unmarshalResponse(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

type ListBlocksRequest struct {
	Address common.Address
	Type    BlockType `etherscan:"blocktype"`
}

type BlockType int32

const (
	BlockTypeBlocks = iota
	BlockTypeUncles
)

func (b BlockType) String() string {
	switch b {
	case BlockTypeBlocks:
		return "blocks"

	case BlockTypeUncles:
		return "uncles"

	default:
		panic(fmt.Sprintf("unknown block type %d", int32(b)))
	}
}

type BlockInfo struct {
	BlockNumber uint64
	Timestamp   time.Time
	BlockReward *big.Int
}

type blockResult struct {
	BlockNumber uintStr
	Timestamp   unixTimestamp
	BlockReward *bigInt
}

func (res *blockResult) toInfo() *BlockInfo {
	return &BlockInfo{
		BlockNumber: res.BlockNumber.unwrap(),
		Timestamp:   res.Timestamp.unwrap(),
		BlockReward: res.BlockReward.unwrap(),
	}
}

func (c *AccountsClient) ListBlocksMined(
	ctx context.Context, req *ListBlocksRequest,
) ([]BlockInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "getminedblocks",
		other:  marshalRequest(req),
	})
	if err != nil {
		return nil, err
	}

	var result []blockResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	blocks := make([]BlockInfo, len(result))
	for i := range result {
		blocks[i] = *result[i].toInfo()
	}

	return blocks, nil
}

type HistoricalETHRequest struct {
	Address     common.Address
	BlockNumber uint64 `etherscan:"blockno"`
}

func (c *AccountsClient) GetHistoricalETHBalance(
	ctx context.Context, req *HistoricalETHRequest,
) (*big.Int, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "balancehistory",
		other:  marshalRequest(req),
	})
	if err != nil {
		return nil, err
	}

	var result bigInt
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result.unwrap(), nil
}
