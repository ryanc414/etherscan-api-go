// Package etherscan is a client library for the etherscan.io API.
package etherscan

import (
	"github.com/ryanc414/etherscan-api-go/accounts"
	"github.com/ryanc414/etherscan-api-go/blocks"
	"github.com/ryanc414/etherscan-api-go/contracts"
	"github.com/ryanc414/etherscan-api-go/gas"
	"github.com/ryanc414/etherscan-api-go/httpapi"
	"github.com/ryanc414/etherscan-api-go/logs"
	"github.com/ryanc414/etherscan-api-go/proxy"
	"github.com/ryanc414/etherscan-api-go/stats"
	"github.com/ryanc414/etherscan-api-go/tokens"
	"github.com/ryanc414/etherscan-api-go/transactions"
)

// Client is the main etherscan client.
type Client struct {
	Accounts     accounts.AccountsClient
	Contracts    contracts.ContractsClient
	Transactions transactions.TransactionsClient
	Blocks       blocks.BlocksClient
	Logs         logs.LogsClient
	Proxy        proxy.ProxyClient
	Tokens       tokens.TokensClient
	Gas          gas.GasClient
	Stats        stats.StatsClient
}

// Params are construction parameters for the etherscan Client.
type Params = httpapi.Params

// New constructs a new etherscan Client.
func New(params *Params) *Client {
	api := httpapi.New(params)
	return &Client{
		Accounts:     accounts.AccountsClient{API: api},
		Contracts:    contracts.ContractsClient{API: api},
		Transactions: transactions.TransactionsClient{API: api},
		Blocks:       blocks.BlocksClient{API: api},
		Logs:         logs.LogsClient{API: api},
		Proxy:        proxy.ProxyClient{API: api},
		Tokens:       tokens.TokensClient{API: api},
		Gas:          gas.GasClient{API: api},
		Stats:        stats.StatsClient{API: api},
	}
}
