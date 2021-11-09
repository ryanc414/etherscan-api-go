package tokens

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ecommon "github.com/ryanc414/etherscan-api-go/common"
	"github.com/ryanc414/etherscan-api-go/httpapi"
	"github.com/ryanc414/etherscan-api-go/marshallers"
)

// TokensClient is the client for tokens actions.
type TokensClient struct {
	API *httpapi.APIClient
}

// GetTotalSupply returns the current amount of an ERC-20 token in circulation.
func (c *TokensClient) GetTotalSupply(
	ctx context.Context, contractAddr common.Address,
) (*big.Int, error) {
	result := new(marshallers.BigInt)
	req := struct{ ContractAddress common.Address }{contractAddr}

	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "tokensupply",
		Request: req,
		Result:  result,
	})
	if err != nil {
		return nil, err
	}

	return result.Unwrap(), nil
}

// BalanceRequest contains the request parameters for GetAccountBalance.
type BalanceRequest struct {
	ContractAddress common.Address
	Address         common.Address
}

// GetAccountBalance returns the current balance of an ERC-20 token of an address.
func (c *TokensClient) GetAccountBalance(
	ctx context.Context, req *BalanceRequest,
) (*big.Int, error) {
	result := new(marshallers.BigInt)
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.AccountsModule,
		Action:  "tokenbalance",
		Request: req,
		Result:  result,
	})
	if err != nil {
		return nil, err
	}

	return result.Unwrap(), nil
}

// HistoricalSupplyRequest contains the request parameters for GetHistoricalSupply.
type HistoricalSupplyRequest struct {
	ContractAddress common.Address
	BlockNo         int64
}

// GetHistoricalSupply returns the amount of an ERC-20 token in circulation at a certain block height.
func (c *TokensClient) GetHistoricalSupply(
	ctx context.Context, req *HistoricalSupplyRequest,
) (*big.Int, error) {
	result := new(marshallers.BigInt)
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "tokensupplyhistory",
		Request: req,
		Result:  result,
	})
	if err != nil {
		return nil, err
	}

	return result.Unwrap(), nil
}

// HistoricalBalanceRequest contains the request parameters for GetHistoricalBalance.
type HistoricalBalanceRequest struct {
	ContractAddress common.Address
	Address         common.Address
	BlockNo         int64
}

// GetHistoricalBalance returns the balance of an ERC-20 token of an address at a certain block height.
func (c *TokensClient) GetHistoricalBalance(
	ctx context.Context, req *HistoricalBalanceRequest,
) (*big.Int, error) {
	result := new(marshallers.BigInt)
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.AccountsModule,
		Action:  "tokenbalancehistory",
		Request: req,
		Result:  result,
	})
	if err != nil {
		return nil, err
	}

	return result.Unwrap(), nil
}

// TokenInfo describes an ERC20/ERC721 token.
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

// GetTokenInfo returns project information and social media links of an ERC-20/ERC-721 token.
func (c *TokensClient) GetTokenInfo(
	ctx context.Context, contractAddr common.Address,
) (result []TokenInfo, err error) {
	req := struct{ ContractAddress common.Address }{contractAddr}
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.TokenModule,
		Action:  "tokeninfo",
		Request: req,
		Result:  &result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
