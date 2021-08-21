package etherscan

import (
	"net/http"
	"net/url"
)

type Client struct {
	Accounts     AccountsClient
	Contracts    ContractsClient
	Transactions TransactionsClient
	Blocks       BlocksClient
	Logs         LogsClient
	Proxy        ProxyClient
	Tokens       TokensClient
	GasTracker   GasClient
	Stats        StatsClient
}

type Params struct {
	APIKey  string
	BaseURL *url.URL
	HTTP    *http.Client
}

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
		GasTracker:   GasClient{api},
		Stats:        StatsClient{api},
	}
}
