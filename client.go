// etherscan is a client library for the etherscan.io API.
package etherscan

import (
	"net/http"
	"net/url"
)

// Client is the main etherscan client.
type Client struct {
	Accounts     AccountsClient
	Contracts    ContractsClient
	Transactions TransactionsClient
	Blocks       BlocksClient
	Logs         LogsClient
	Proxy        ProxyClient
	Tokens       TokensClient
	Gas          GasClient
	Stats        StatsClient
}

// Params are construction parameters for the etherscan Client.
type Params struct {
	APIKey  string
	BaseURL *url.URL
	HTTP    *http.Client
}

// New constructs a new etherscan Client.
func New(params *Params) *Client {
	api := newAPIClient(params)
	return &Client{
		Accounts:     AccountsClient{api},
		Contracts:    ContractsClient{api},
		Transactions: TransactionsClient{api},
		Blocks:       BlocksClient{api},
		Logs:         LogsClient{api},
		Proxy:        ProxyClient{api},
		Tokens:       TokensClient{api},
		Gas:          GasClient{api},
		Stats:        StatsClient{api},
	}
}
