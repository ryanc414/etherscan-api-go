package etherscan

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

type TransactionsClient struct {
	api *apiClient
}

type ExecutionStatus struct {
	IsError        bool
	ErrDescription string
}

type executionStatusResult struct {
	IsError        string `json:"isError"`
	ErrDescription string `json:"errDescription"`
}

func (res *executionStatusResult) toStatus() *ExecutionStatus {
	return &ExecutionStatus{
		IsError:        res.IsError != "0",
		ErrDescription: res.ErrDescription,
	}
}

const transactionModule = "transaction"

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

	var result executionStatusResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result.toStatus(), nil
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
