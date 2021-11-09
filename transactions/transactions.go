package transactions

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ryanc414/etherscan-api-go/httpapi"
)

const transactionModule = "transaction"

// TransactionsClient is the client for transaction actions.
type TransactionsClient struct {
	API *httpapi.APIClient
}

// ExecutionStatus describes the current state of a transaction.
type ExecutionStatus struct {
	IsError        bool   `etherscan:"isError,num"`
	ErrDescription string `etherscan:"errDescription"`
}

// GetExecutionStatus returns he status code of a contract execution.
func (c *TransactionsClient) GetExecutionStatus(
	ctx context.Context, txHash common.Hash,
) (*ExecutionStatus, error) {
	req := struct{ TxHash common.Hash }{txHash}
	result := new(ExecutionStatus)

	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  transactionModule,
		Action:  "getstatus",
		Request: req,
		Result:  result,
	})

	return result, err
}

type txReceiptStatusResult struct {
	Status bool `etherscan:"status,num"`
}

// GetTxReceiptStatus returns the status code of a transaction execution.
func (c *TransactionsClient) GetTxReceiptStatus(
	ctx context.Context, txHash common.Hash,
) (bool, error) {
	req := struct{ TxHash common.Hash }{txHash}
	result := new(txReceiptStatusResult)
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  transactionModule,
		Action:  "gettxreceiptstatus",
		Request: req,
		Result:  result,
	})
	if err != nil {
		return false, err
	}

	return result.Status, nil
}
