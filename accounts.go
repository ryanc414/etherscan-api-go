package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
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

func (req *ETHBalanceRequest) toParams() (map[string]string, error) {
	if req.Address == (common.Address{}) {
		return nil, errors.New("address is required")
	}

	tag := req.Tag
	if tag == blockParameterUnspecified {
		tag = BlockParameterLatest
	}

	return map[string]string{
		"address": req.Address.String(),
		"tag":     tag.String(),
	}, nil
}

type BlockParameter int32

const (
	blockParameterUnspecified = iota
	BlockParameterEarliest
	BlockParameterPending
	BlockParameterLatest
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
	params, err := req.toParams()
	if err != nil {
		return nil, err
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "balance",
		other:  params,
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
	Addresses []common.Address
	Tag       BlockParameter
}

func (req *MultiETHBalancesRequest) toParams() (map[string]string, error) {
	if len(req.Addresses) == 0 {
		return nil, errors.New("Addresses must be provided")
	}

	joinedAddrs := joinAddresses(req.Addresses)

	tag := req.Tag
	if tag == blockParameterUnspecified {
		tag = BlockParameterLatest
	}

	return map[string]string{
		"address": joinedAddrs,
		"tag":     tag.String(),
	}, nil
}

type MultiBalanceResponse struct {
	Account common.Address
	Balance *big.Int
}

type multiBalanceResult struct {
	Account string  `json:"account"`
	Balance *bigInt `json:"balance"`
}

func (r *multiBalanceResult) toResponse() (*MultiBalanceResponse, error) {
	return &MultiBalanceResponse{
		Account: common.HexToAddress(r.Account),
		Balance: r.Balance.unwrap(),
	}, nil
}

func (c *AccountsClient) GetMultiETHBalances(
	ctx context.Context, req *MultiETHBalancesRequest,
) ([]MultiBalanceResponse, error) {
	params, err := req.toParams()
	if err != nil {
		return nil, err
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "balancemulti",
		other:  params,
	})
	if err != nil {
		return nil, err
	}

	var result []multiBalanceResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	response := make([]MultiBalanceResponse, len(result))
	for i := range result {
		rsp, err := result[i].toResponse()
		if err != nil {
			return nil, err
		}

		response[i] = *rsp
	}

	return response, nil
}

func joinAddresses(addresses []common.Address) string {
	addrStrs := make([]string, len(addresses))
	for i := range addresses {
		addrStrs[i] = addresses[i].String()
	}

	return strings.Join(addrStrs, ",")
}

type ListTxRequest struct {
	Address    common.Address
	StartBlock uint64
	EndBlock   uint64
	Sort       SortingPreference
}

func (req *ListTxRequest) toParams() (map[string]string, error) {
	if req.Address == (common.Address{}) {
		return nil, errors.New("address must be specified")
	}

	return map[string]string{
		"address":    req.Address.String(),
		"startblock": strconv.FormatUint(req.StartBlock, 10),
		"endblock":   strconv.FormatUint(req.EndBlock, 10),
		"sort":       req.Sort.String(),
	}, nil
}

type SortingPreference int32

const (
	sortingPreferenceUnspecified = iota
	SortingPreferenceAscending
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
	BlockNumber     uint64
	Timestamp       time.Time
	Hash            common.Hash
	From            common.Address
	To              common.Address
	Value           *big.Int
	ContractAddress *common.Address
	Input           []byte
	Gas             uint64
	GasUsed         uint64
	IsError         bool
	ErrCode         string
}

type transactionResult struct {
	BlockNumber     uintStr `json:"blockNumber"`
	Timestamp       string  `json:"timeStamp"`
	Hash            string  `json:"hash"`
	From            string  `json:"from"`
	To              string  `json:"to"`
	Value           *bigInt `json:"value"`
	ContractAddress string  `json:"contractAddress"`
	Input           string  `json:"input"`
	Gas             uintStr `json:"gas"`
	GasUsed         uintStr `json:"gasUsed"`
	IsError         string  `json:"isError"`
}

func (tx *transactionResult) toInfo() (*TransactionInfo, error) {
	timestampUnix, err := strconv.ParseInt(tx.Timestamp, 10, 64)
	if err != nil {
		return nil, err
	}

	var contractAddress *common.Address
	if tx.ContractAddress != "" {
		addr := common.HexToAddress(tx.ContractAddress)
		contractAddress = &addr
	}

	return &TransactionInfo{
		BlockNumber:     tx.BlockNumber.unwrap(),
		Timestamp:       time.Unix(timestampUnix, 0),
		Hash:            common.HexToHash(tx.Hash),
		From:            common.HexToAddress(tx.From),
		To:              common.HexToAddress(tx.To),
		Value:           tx.Value.unwrap(),
		Gas:             tx.Gas.unwrap(),
		IsError:         tx.IsError != "0",
		Input:           common.Hex2Bytes(tx.Input),
		ContractAddress: contractAddress,
		GasUsed:         tx.GasUsed.unwrap(),
	}, nil
}

type NormalTxInfo struct {
	TransactionInfo
	Nonce             uint64
	BlockHash         common.Hash
	TransactionIndex  uint64
	GasPrice          *big.Int
	TxReceiptStatus   string
	CumulativeGasUsed uint64
	Confirmations     uint64
}

type normalTxResult struct {
	transactionResult
	Nonce             uintStr `json:"nonce"`
	BlockHash         string  `json:"blockHash"`
	TransactionIndex  uintStr `json:"transactionIndex"`
	GasPrice          *bigInt `json:"gasPrice"`
	TxReceiptStatus   string  `json:"txreceipt_status"`
	CumulativeGasUsed uintStr `json:"cumulativeGasUsed"`
	Confirmations     uintStr `json:"confirmations"`
}

func (tx *normalTxResult) toInfo() (*NormalTxInfo, error) {
	baseTx, err := tx.transactionResult.toInfo()
	if err != nil {
		return nil, err
	}

	return &NormalTxInfo{
		TransactionInfo:   *baseTx,
		Nonce:             tx.Nonce.unwrap(),
		BlockHash:         common.HexToHash(tx.BlockHash),
		TransactionIndex:  tx.TransactionIndex.unwrap(),
		GasPrice:          tx.GasPrice.unwrap(),
		TxReceiptStatus:   tx.TxReceiptStatus,
		CumulativeGasUsed: tx.CumulativeGasUsed.unwrap(),
		Confirmations:     tx.Confirmations.unwrap(),
	}, nil
}

type InternalTxInfo struct {
	TransactionInfo
	TraceID string
	Type    string
}

type internalTxResult struct {
	transactionResult
	TraceID string `json:"traceId"`
	Type    string `json:"type"`
}

func (tx *internalTxResult) toInfo() (*InternalTxInfo, error) {
	baseTx, err := tx.transactionResult.toInfo()
	if err != nil {
		return nil, err
	}

	return &InternalTxInfo{
		TransactionInfo: *baseTx,
		TraceID:         tx.TraceID,
		Type:            tx.Type,
	}, nil
}

func (c *AccountsClient) ListNormalTransactions(
	ctx context.Context, req *ListTxRequest,
) ([]NormalTxInfo, error) {
	params, err := req.toParams()
	if err != nil {
		return nil, err
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "txlist",
		other:  params,
	})
	if err != nil {
		return nil, err
	}

	var result []normalTxResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	txInfos := make([]NormalTxInfo, len(result))
	for i := range result {
		info, err := result[i].toInfo()
		if err != nil {
			return nil, err
		}

		txInfos[i] = *info
	}

	return txInfos, nil
}

func (c *AccountsClient) ListInternalTransactions(
	ctx context.Context, req *ListTxRequest,
) ([]InternalTxInfo, error) {
	params, err := req.toParams()
	if err != nil {
		return nil, err
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "txlistinternal",
		other:  params,
	})
	if err != nil {
		return nil, err
	}

	return unmarshalInternalTxs(rspData)
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

	return unmarshalInternalTxs(rspData)
}

type BlockRangeRequest struct {
	StartBlock uint64
	EndBlock   uint64
	Sort       SortingPreference
}

func (req *BlockRangeRequest) toParams() map[string]string {
	return map[string]string{
		"startblock": strconv.FormatUint(req.StartBlock, 10),
		"endblock":   strconv.FormatUint(req.EndBlock, 10),
		"sort":       req.Sort.String(),
	}
}

func (c *AccountsClient) GetInternalTxsByBlockRange(
	ctx context.Context,
	req *BlockRangeRequest,
) ([]InternalTxInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "txlistinternal",
		other:  req.toParams(),
	})
	if err != nil {
		return nil, err
	}

	return unmarshalInternalTxs(rspData)
}

func unmarshalInternalTxs(rspData []byte) ([]InternalTxInfo, error) {
	var result []internalTxResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	txInfos := make([]InternalTxInfo, len(result))
	for i := range result {
		info, err := result[i].toInfo()
		if err != nil {
			return nil, err
		}

		txInfos[i] = *info
	}

	return txInfos, nil
}

type TokenTransfersRequest struct {
	Address         common.Address
	ContractAddress common.Address
	Sort            SortingPreference
}

func (req *TokenTransfersRequest) toParams() map[string]string {
	return map[string]string{
		"address":         req.Address.String(),
		"contractaddress": req.ContractAddress.String(),
		"sort":            req.Sort.String(),
	}
}

type TokenTransferInfo struct {
	NormalTxInfo
	TokenName    string
	TokenSymbol  string
	TokenDecimal uint32
}

type tokenTransferResult struct {
	normalTxResult
	TokenName    string  `json:"tokenName"`
	TokenSymbol  string  `json:"tokenSymbol"`
	TokenDecimal uintStr `json:"tokenDecimal"`
}

func (res *tokenTransferResult) toInfo() (*TokenTransferInfo, error) {
	baseTx, err := res.normalTxResult.toInfo()
	if err != nil {
		return nil, err
	}

	return &TokenTransferInfo{
		NormalTxInfo: *baseTx,
		TokenName:    res.TokenName,
		TokenSymbol:  res.TokenSymbol,
		TokenDecimal: uint32(res.TokenDecimal),
	}, nil
}

func (c *AccountsClient) ListTokenTransfers(
	ctx context.Context, req *TokenTransfersRequest,
) ([]TokenTransferInfo, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: "account",
		action: "tokentx",
		other:  req.toParams(),
	})
	if err != nil {
		return nil, err
	}

	var result []tokenTransferResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	tokenInfos := make([]TokenTransferInfo, len(result))
	for i := range result {
		info, err := result[i].toInfo()
		if err != nil {
			return nil, err
		}

		tokenInfos[i] = *info
	}

	return tokenInfos, nil
}

type ListNFTTransferRequest struct {
	Address         *common.Address
	ContractAddress *common.Address
	Sort            SortingPreference
}

func (req *ListNFTTransferRequest) toParams() (map[string]string, error) {
	if req.Address == nil && req.ContractAddress == nil {
		return nil, errors.New("at least one of Address or ContractAddress must be specifide")
	}

	params := map[string]string{"sort": req.Sort.String()}
	if req.Address != nil {
		params["address"] = req.Address.String()
	}

	if req.ContractAddress != nil {
		params["contractaddress"] = req.ContractAddress.String()
	}

	return params, nil
}

type NFTTransferInfo struct {
	TokenTransferInfo
	TokenID string
}

type nftTransferResult struct {
	tokenTransferResult
	TokenID string `json:"tokenID"`
}

func (res *nftTransferResult) toInfo() (*NFTTransferInfo, error) {
	baseTx, err := res.tokenTransferResult.toInfo()
	if err != nil {
		return nil, err
	}

	return &NFTTransferInfo{
		TokenTransferInfo: *baseTx,
		TokenID:           res.TokenID,
	}, nil
}

func (c *AccountsClient) ListNFTTransfers(
	ctx context.Context, req *ListNFTTransferRequest,
) ([]NFTTransferInfo, error) {
	params, err := req.toParams()
	if err != nil {
		return nil, err
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "tokennfttx",
		other:  params,
	})
	if err != nil {
		return nil, err
	}

	var result []nftTransferResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	nftInfos := make([]NFTTransferInfo, len(result))
	for i := range result {
		info, err := result[i].toInfo()
		if err != nil {
			return nil, err
		}

		nftInfos[i] = *info
	}

	return nftInfos, nil
}
