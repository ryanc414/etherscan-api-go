package etherscan

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

const transactionModule = "transaction"

type TransactionsClient struct {
	api *apiClient
}

type ExecutionStatus struct {
	IsError        bool   `etherscan:"isError,num"`
	ErrDescription string `etherscan:"errDescription"`
}

func (c *TransactionsClient) GetExecutionStatus(
	ctx context.Context, txHash common.Hash,
) (*ExecutionStatus, error) {
	req := struct{ TxHash common.Hash }{txHash}
	result := new(ExecutionStatus)

	err := c.api.call(ctx, &callParams{
		module:  transactionModule,
		action:  "getstatus",
		request: req,
		result:  result,
	})

	return result, err
}

type txReceiptStatusResult struct {
	Status bool `etherscan:"status,num"`
}

func (c *TransactionsClient) GetTxReceiptStatus(
	ctx context.Context, txHash common.Hash,
) (bool, error) {
	req := struct{ TxHash common.Hash }{txHash}
	result := new(txReceiptStatusResult)
	err := c.api.call(ctx, &callParams{
		module:  transactionModule,
		action:  "gettxreceiptstatus",
		request: req,
		result:  result,
	})
	if err != nil {
		return false, err
	}

	return result.Status, nil
}
