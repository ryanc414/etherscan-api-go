package etherscan

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (
	statsModule = "stats"
	tokenModule = "token"
)

type TokensClient struct {
	api *apiClient
}

func (c *TokensClient) GetTotalSupply(
	ctx context.Context, contractAddr common.Address,
) (*big.Int, error) {
	result := new(bigInt)
	req := struct{ ContractAddress common.Address }{contractAddr}

	err := c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "tokensupply",
		request: req,
		result:  result,
	})
	if err != nil {
		return nil, err
	}

	return result.unwrap(), nil
}

type BalanceRequest struct {
	ContractAddress common.Address
	Address         common.Address
}

func (c *TokensClient) GetAccountBalance(
	ctx context.Context, req *BalanceRequest,
) (*big.Int, error) {
	result := new(bigInt)
	err := c.api.call(ctx, &callParams{
		module:  accountModule,
		action:  "tokenbalance",
		request: req,
		result:  result,
	})
	if err != nil {
		return nil, err
	}

	return result.unwrap(), nil
}

type HistoricalSupplyRequest struct {
	ContractAddress common.Address
	BlockNo         int64
}

func (c *TokensClient) GetHistoricalSupply(
	ctx context.Context, req *HistoricalSupplyRequest,
) (*big.Int, error) {
	result := new(bigInt)
	err := c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "tokensupplyhistory",
		request: req,
		result:  result,
	})
	if err != nil {
		return nil, err
	}

	return result.unwrap(), nil
}

type HistoricalBalanceRequest struct {
	ContractAddress common.Address
	Address         common.Address
	BlockNo         int64
}

func (c *TokensClient) GetHistoricalBalance(
	ctx context.Context, req *HistoricalBalanceRequest,
) (*big.Int, error) {
	result := new(bigInt)
	err := c.api.call(ctx, &callParams{
		module:  accountModule,
		action:  "tokenbalancehistory",
		request: req,
		result:  result,
	})
	if err != nil {
		return nil, err
	}

	return result.unwrap(), nil
}

type TokenInfo struct {
	ContractAddress common.Address `etherscan:"contractAddress"`
	TokenName       string         `etherscan:"tokenName"`
	Symbol          string
	Divisor         uint
	TokenType       string   `etherscan:"tokenType"`
	TotalSupply     *big.Int `etherscan:"totalSupply"`
	BlueCheckmark   bool     `etherscan:"blueCheckmark,str"`
	Description     string
	Website         string
	Email           string
	Blog            string
	Reddit          string
	Slack           string
	Facebook        string
	Twitter         string
	BitcoinTalk     string
	Github          string
	Telegram        string
	WeChat          string
	LinkedIn        string
	Discord         string
	Whitepaper      string
	TokenPriceUSD   float64 `etherscan:"tokenPriceUSD"`
}

func (c *TokensClient) GetTokenInfo(
	ctx context.Context, contractAddr common.Address,
) (result []TokenInfo, err error) {
	req := struct{ ContractAddress common.Address }{contractAddr}
	err = c.api.call(ctx, &callParams{
		module:  tokenModule,
		action:  "tokeninfo",
		request: req,
		result:  &result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
