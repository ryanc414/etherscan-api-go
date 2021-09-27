package etherscan

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

const transactionModule = "transaction"

type TransactionsClient struct {
	api *apiClient
}

type ExecutionStatus struct {
	IsError        bool   `etherscan:"isError"`
	ErrDescription string `etherscan:"errDescription"`
}

func (c *TransactionsClient) GetExecutionStatus(
	ctx context.Context, txHash common.Hash,
) (*ExecutionStatus, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: transactionModule,
		action: "getstatus",
		other:  map[string]string{"txhash": txHash.String()},
	})
	if err != nil {
		return nil, err
	}

	result := new(ExecutionStatus)
	if err := unmarshalResponse(rspData, result); err != nil {
		return nil, err
	}

	return result, nil
}

type txReceiptStatusResult struct {
	Status string `json:"status"`
}

func (c *TransactionsClient) GetTxReceiptStatus(
	ctx context.Context, txHash common.Hash,
) (bool, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: transactionModule,
		action: "gettxreceiptstatus",
		other:  map[string]string{"txhash": txHash.String()},
	})
	if err != nil {
		return false, err
	}

	var result txReceiptStatusResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return false, err
	}

	return result.Status != "0", nil
}
