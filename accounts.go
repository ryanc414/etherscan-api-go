package etherscan

import (
	"context"
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
	result := new(bigInt)
	err := c.api.call(ctx, &callParams{
		module:  accountModule,
		action:  "balance",
		request: req,
		result:  result,
	})
	if err != nil {
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
) (result []MultiBalanceResponse, err error) {
	err = c.api.call(ctx, &callParams{
		module:  accountModule,
		action:  "balancemulti",
		request: req,
		result:  &result,
	})
	return result, err
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
	IsError         bool   `etherscan:"isError,num"`
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
) (result []NormalTxInfo, err error) {
	err = c.api.call(ctx, &callParams{
		module:  accountModule,
		action:  "txlist",
		request: req,
		result:  &result,
	})

	return result, err
}

func (c *AccountsClient) ListInternalTransactions(
	ctx context.Context, req *ListTxRequest,
) (result []InternalTxInfo, err error) {
	err = c.api.call(ctx, &callParams{
		module:  accountModule,
		action:  "txlistinternal",
		request: req,
		result:  &result,
	})

	return result, err
}

func (c *AccountsClient) GetInternalTxsByHash(
	ctx context.Context, hash common.Hash,
) (result []InternalTxInfo, err error) {
	req := struct{ TxHash common.Hash }{hash}
	err = c.api.call(ctx, &callParams{
		module:  accountModule,
		action:  "txlistinternal",
		request: req,
		result:  &result,
	})

	return result, err
}

type BlockRangeRequest struct {
	StartBlock uint64
	EndBlock   uint64
	Sort       SortingPreference
}

func (c *AccountsClient) GetInternalTxsByBlockRange(
	ctx context.Context,
	req *BlockRangeRequest,
) (result []InternalTxInfo, err error) {
	err = c.api.call(ctx, &callParams{
		module:  accountModule,
		action:  "txlistinternal",
		request: req,
		result:  &result,
	})

	return result, err
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
) (result []TokenTransferInfo, err error) {
	err = c.api.call(ctx, &callParams{
		module:  "account",
		action:  "tokentx",
		request: req,
		result:  &result,
	})

	return result, err
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
) (result []NFTTransferInfo, err error) {
	if req.Address == nil && req.ContractAddress == nil {
		return nil, errors.New("at least one of Address or ContractAddress must be specifide")
	}

	err = c.api.call(ctx, &callParams{
		module:  accountModule,
		action:  "tokennfttx",
		request: req,
		result:  &result,
	})

	return result, err
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
	BlockNumber uint64 `etherscan:"blockNumber"`
	Timestamp   time.Time
	BlockReward *big.Int `etherscan:"blockReward"`
}

func (c *AccountsClient) ListBlocksMined(
	ctx context.Context, req *ListBlocksRequest,
) (result []BlockInfo, err error) {
	err = c.api.call(ctx, &callParams{
		module:  accountModule,
		action:  "getminedblocks",
		request: req,
		result:  &result,
	})

	return result, err
}

type HistoricalETHRequest struct {
	Address     common.Address
	BlockNumber uint64 `etherscan:"blockno"`
}

func (c *AccountsClient) GetHistoricalETHBalance(
	ctx context.Context, req *HistoricalETHRequest,
) (result *big.Int, err error) {
	err = c.api.call(ctx, &callParams{
		module:  accountModule,
		action:  "balancehistory",
		request: req,
		result:  &result,
	})

	return result, err
}
