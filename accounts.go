//go:generate go-enum -f=$GOFILE
package etherscan

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const accountModule = "account"

// AccountsClient is the client for accounts actions.
type AccountsClient struct {
	api *apiClient
}

// ETHBalanceRequest contains the request parameters for GetETHBalance.
type ETHBalanceRequest struct {
	Address common.Address
	Tag     BlockParameter
}

// BlockParameter is an enumeration of allowed block parameters.
// ENUM(latest,earliest,pending)
type BlockParameter int32

// GetETHBalance returns the Ether balance for a single address.
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

// MultiETHBalancesRequest contains the request parameters for GetMultiETHBalances.
type MultiETHBalancesRequest struct {
	Addresses []common.Address `etherscan:"address"`
	Tag       BlockParameter
}

// MultiBalanceResponse contains the Ether balance for a specific address.
type MultiBalanceResponse struct {
	Account common.Address
	Balance *big.Int
}

// GetMultiETHBalances returns the balance of accounts from a list of addresses.
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

// ListTxRequest contains the request parameters for ListNormalTransactions.
type ListTxRequest struct {
	Address    common.Address
	StartBlock uint64
	EndBlock   uint64
	Sort       SortingPreference
}

// SortingPreference is an enumeration of sorting preferences.
// ENUM(asc,desc)
type SortingPreference int32

// TransactionInfo conatins the base transaction info included in multiple
// response types.
type TransactionInfo struct {
	BlockNumber     uint64    `etherscan:"blockNumber"`
	Timestamp       time.Time `etherscan:"timeStamp"`
	From            common.Address
	To              common.Address
	Value           *big.Int
	ContractAddress *common.Address `etherscan:"contractAddress"`
	Input           []byte
	Gas             uint64
	GasUsed         uint64 `etherscan:"gasUsed"`
	IsError         bool   `etherscan:"isError,num"`
}

// NormalTxInfo contains information on normal transactions returned by ListNormalTransactions.
type NormalTxInfo struct {
	TransactionInfo
	Hash              common.Hash
	Nonce             uint64
	BlockHash         common.Hash `etherscan:"blockHash"`
	TransactionIndex  uint64      `etherscan:"transactionIndex"`
	GasPrice          *big.Int    `etherscan:"gasPrice"`
	TxReceiptStatus   string      `etherscan:"txreceipt_status"`
	CumulativeGasUsed uint64      `etherscan:"cumulativeGasUsed"`
	Confirmations     uint64
}

// InternalTxInfo contains information on internal transactions.
type InternalTxInfo struct {
	TransactionInfo
	Hash    common.Hash
	TraceID string `etherscan:"traceId"`
	Type    string
}

// InternalTxInfoByHash contains information on internal transactions returned
// by GetInternalTxsByHash
type InternalTxInfoByHash struct {
	TransactionInfo
	Type string
}

// ListNormalTransactions returns the list of transactions performed by an address.
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

// ListInternalTransactions returns the list of internal transactions performed by an address.
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

// GetInternalTxsByHash returns the list of internal transactions performed within a transaction.
func (c *AccountsClient) GetInternalTxsByHash(
	ctx context.Context, hash common.Hash,
) (result []InternalTxInfoByHash, err error) {
	req := struct{ TxHash common.Hash }{hash}
	err = c.api.call(ctx, &callParams{
		module:  accountModule,
		action:  "txlistinternal",
		request: req,
		result:  &result,
	})

	return result, err
}

// BlockRangeRequest contains the request parameters for GetInternalTxsByBlockRange.
type BlockRangeRequest struct {
	StartBlock uint64
	EndBlock   uint64
	Sort       SortingPreference
}

// GetInternalTxsByBlockRange returns the list of internal transactions performed within a block range.
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

// TokenTransferRequest
type TokenTransfersRequest struct {
	Address         common.Address
	ContractAddress common.Address
	Sort            SortingPreference
}

// BaseTokenTransferInfo contains common token transfer information.
type BaseTokenTransferInfo struct {
	BlockNumber       uint64    `etherscan:"blockNumber"`
	Timestamp         time.Time `etherscan:"timeStamp"`
	Hash              common.Hash
	Nonce             uint64
	BlockHash         common.Hash `etherscan:"blockHash"`
	From              common.Address
	ContractAddress   common.Address `etherscan:"contractAddress"`
	To                common.Address
	TokenName         string `etherscan:"tokenName"`
	TokenSymbol       string `etherscan:"tokenSymbol"`
	TokenDecimal      uint32 `etherscan:"tokenDecimal"`
	TransactionIndex  uint32 `etherscan:"transactionIndex"`
	Gas               uint64
	GasPrice          *big.Int `etherscan:"gasPrice"`
	GasUsed           uint64   `etherscan:"gasUsed"`
	CumulativeGasUsed uint64   `etherscan:"cumulativeGasUsed"`
	Confirmations     uint64
}

// TokenTransferInfo contains information on an ERC20 token transfer.
type TokenTransferInfo struct {
	BaseTokenTransferInfo
	Value *big.Int
}

// ListTokenTransfers lists the ERC20 token transfers for an address.
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

// ListNFTTransferRequest contains the request parameters for ListNFTTransfers.
type ListNFTTransferRequest struct {
	Address         *common.Address
	ContractAddress *common.Address
	Sort            SortingPreference
}

// NFTTransferInfo contains the information on an NFT token transfer.
type NFTTransferInfo struct {
	BaseTokenTransferInfo
	TokenID string `etherscan:"tokenID"`
}

// ListNFTTransfers lists the NFT token transfers for an address.
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

// ListBlocksRequest contains the request parameters for ListBlocksMined.
type ListBlocksRequest struct {
	Address common.Address
	Type    BlockType `etherscan:"blocktype"`
}

// BlockType is an enumeration of block types.
// ENUM(blocks,uncles)
type BlockType int32

// BlockInfo contains information on a specific ethereum block.
type BlockInfo struct {
	BlockNumber uint64    `etherscan:"blockNumber"`
	Timestamp   time.Time `etherscan:"timeStamp"`
	BlockReward *big.Int  `etherscan:"blockReward"`
}

// ListBlocksMined lists blocks that were mined by a specific address.
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

// HistoricalETHRequest contains the request parameters for GetHistoricalETHBalance.
type HistoricalETHRequest struct {
	Address     common.Address
	BlockNumber uint64 `etherscan:"blockno"`
}

// GetHistoricalETHBalance retrieves the ETH balance for an address at a
// particular block number.
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
