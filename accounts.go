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
	if req.Address == (common.Address{}) {
		return nil, errors.New("address is required")
	}

	tag := req.Tag
	if tag == blockParameterUnspecified {
		tag = BlockParameterLatest
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "balance",
		other: map[string]string{
			"address": req.Address.String(),
			"tag":     tag.String(),
		},
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
	if len(req.Addresses) == 0 {
		return nil, errors.New("Addresses must be provided")
	}

	joinedAddrs := joinAddresses(req.Addresses)

	tag := req.Tag
	if tag == blockParameterUnspecified {
		tag = BlockParameterLatest
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "balancemulti",
		other: map[string]string{
			"address": joinedAddrs,
			"tag":     tag.String(),
		},
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

type ListNormalTxRequest struct {
	Address    common.Address
	StartBlock uint64
	EndBlock   uint64
	Page       uint32
	Offset     uint32
	Sort       SortingPreference
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
	BlockNumber       uint64
	Timestamp         time.Time
	Hash              common.Hash
	Nonce             uint64
	BlockHash         common.Hash
	TransactionIndex  uint64
	From              common.Address
	To                common.Address
	Value             *big.Int
	Gas               uint64
	GasPrice          *big.Int
	IsError           bool
	TxReceiptStatus   string
	Input             []byte
	ContractAddress   *common.Address
	CumulativeGasUsed uint64
	GasUsed           uint64
	Confirmations     uint64
}

type transactionResult struct {
	BlockNumber       uintStr `json:"blockNumber"`
	Timestamp         string  `json:"timeStamp"`
	Hash              string  `json:"hash"`
	Nonce             uintStr `json:"nonce"`
	BlockHash         string  `json:"blockHash"`
	TransactionIndex  uintStr `json:"transactionIndex"`
	From              string  `json:"from"`
	To                string  `json:"to"`
	Value             *bigInt `json:"value"`
	Gas               uintStr `json:"gas"`
	GasPrice          *bigInt `json:"gasPrice"`
	IsError           string  `json:"isError"`
	TxReceiptStatus   string  `json:"txreceipt_status"`
	Input             string  `json:"input"`
	ContractAddress   string  `json:"contractAddress"`
	CumulativeGasUsed uintStr `json:"cumulativeGasUsed"`
	GasUsed           uintStr `json:"gasUsed"`
	Confirmations     uintStr `json:"confirmations"`
}

func (r *transactionResult) toInfo() (*TransactionInfo, error) {
	timestampUnix, err := strconv.ParseInt(r.Timestamp, 10, 64)
	if err != nil {
		return nil, err
	}

	var contractAddress *common.Address
	if r.ContractAddress != "" {
		addr := common.HexToAddress(r.ContractAddress)
		contractAddress = &addr
	}

	return &TransactionInfo{
		BlockNumber:       r.BlockNumber.unwrap(),
		Timestamp:         time.Unix(timestampUnix, 0),
		Hash:              common.HexToHash(r.Hash),
		Nonce:             r.Nonce.unwrap(),
		BlockHash:         common.HexToHash(r.BlockHash),
		TransactionIndex:  r.TransactionIndex.unwrap(),
		From:              common.HexToAddress(r.From),
		To:                common.HexToAddress(r.To),
		Value:             r.Value.unwrap(),
		Gas:               r.Gas.unwrap(),
		GasPrice:          r.GasPrice.unwrap(),
		IsError:           r.IsError != "0",
		TxReceiptStatus:   r.TxReceiptStatus,
		Input:             common.Hex2Bytes(r.Input),
		ContractAddress:   contractAddress,
		CumulativeGasUsed: r.CumulativeGasUsed.unwrap(),
		GasUsed:           r.GasUsed.unwrap(),
		Confirmations:     r.Confirmations.unwrap(),
	}, nil
}

func (c *AccountsClient) ListNormalTransactions(
	ctx context.Context, req *ListNormalTxRequest,
) ([]TransactionInfo, error) {
	if req.Address == (common.Address{}) {
		return nil, errors.New("address must be specified")
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: accountModule,
		action: "txlist",
		other: map[string]string{
			"address":    req.Address.String(),
			"startblock": strconv.FormatUint(req.StartBlock, 10),
			"endblock":   strconv.FormatUint(req.EndBlock, 10),
			"page":       strconv.FormatUint(uint64(req.Page), 10),
			"offset":     strconv.FormatUint(uint64(req.Offset), 10),
			"sort":       req.Sort.String(),
		},
	})
	if err != nil {
		return nil, err
	}

	var result []transactionResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	txInfos := make([]TransactionInfo, len(result))
	for i := range result {
		info, err := result[i].toInfo()
		if err != nil {
			return nil, err
		}

		txInfos[i] = *info
	}

	return txInfos, nil
}
