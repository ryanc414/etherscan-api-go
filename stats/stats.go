//go:generate go-enum -f=$GOFILE --marshal
package stats

import (
	"context"
	"math/big"
	"time"

	ecommon "github.com/ryanc414/etherscan-api-go/common"
	"github.com/ryanc414/etherscan-api-go/httpapi"
	"github.com/ryanc414/etherscan-api-go/marshallers"
	"github.com/shopspring/decimal"
)

// StatsClient is the client for stats actions.
type StatsClient struct {
	API *httpapi.APIClient
}

// GetToalETHSupply returns the current amount of Ether in circulation.
func (c *StatsClient) GetTotalETHSupply(ctx context.Context) (*big.Int, error) {
	result := new(marshallers.BigInt)
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module: ecommon.StatsModule,
		Action: "ethsupply",
		Result: result,
	})
	if err != nil {
		return nil, err
	}

	return result.Unwrap(), nil
}

// ETHPrice describes the price of Ether at a particular time.
type ETHPrice struct {
	ETHBTC          decimal.Decimal
	ETHBTCTimestamp time.Time `etherscan:"ethbtc_timestamp"`
	ETHUSD          decimal.Decimal
	ETHUSDTimestamp time.Time `etherscan:"ethusd_timestamp"`
}

// GetLastETHPrice returns the latest price of 1 ETH.
func (c *StatsClient) GetLastETHPrice(ctx context.Context) (*ETHPrice, error) {
	result := new(ETHPrice)
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module: ecommon.StatsModule,
		Action: "ethprice",
		Result: result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// NodesSizeReq contains the request parameters for GetEthereumNodesSize.
type NodesSizeReq struct {
	StartDate  time.Time `etherscan:"startdate,date"`
	EndDate    time.Time `etherscan:"enddate,date"`
	ClientType ETHClientTypeReq
	SyncMode   NodeSyncModeReq
	Sort       ecommon.SortingPreference
}

// ETHClientTypeReq is an enumeration of ethereum client types.
// ENUM(geth,parity)
type ETHClientTypeReq int32

// ETHClientTypeResult is an enumeration of ethereum client types.
// ENUM(Geth,Parity)
type ETHClientTypeResult int32

// NodeSyncModeReq is an enumeration of ethereum node sync modes.
// ENUM(default,archive)
type NodeSyncModeReq int32

// NodeSyncModeResult is an enumeration of ethereum node sync modes.
// ENUM(Default,Archive)
type NodeSyncModeResult int32

// ETHNodeSize describes the size of the ethereum blockchain at a particular
// block number.
type ETHNodeSize struct {
	BlockNumber    uint64              `etherscan:"blockNumber"`
	ChainTimestamp time.Time           `etherscan:"chainTimeStamp,date"`
	ChainSize      *big.Int            `etherscan:"chainSize"`
	ClientType     ETHClientTypeResult `etherscan:"clientType"`
	SyncMode       NodeSyncModeResult  `etherscan:"syncMode"`
}

// GetEthereumNodesSize returns the size of the Ethereum blockchain, in bytes,
// over a date range.
func (c *StatsClient) GetEthereumNodesSize(
	ctx context.Context, req *NodesSizeReq,
) (result []ETHNodeSize, err error) {
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "chainsize",
		Request: req,
		Result:  &result,
	})
	return result, err
}

// NodeCount describes the total count of nodes on the ethereum network on a
// particular date.
type NodeCount struct {
	Date           time.Time `etherscan:"UTCDate,date"`
	TotalNodeCount uint64    `etherscan:"TotalNodeCount"`
}

// GetToalNodesCount returns the total number of discoverable Ethereum nodes.
func (c *StatsClient) GetTotalNodesCount(ctx context.Context) (*NodeCount, error) {
	result := new(NodeCount)
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module: ecommon.StatsModule,
		Action: "nodecount",
		Result: result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DailyTxFee describes the total amount of transaction fees on a particular day.
type DailyTxFee struct {
	Timestamp time.Time       `etherscan:"unixTimeStamp"`
	TxFeeETH  decimal.Decimal `etherscan:"transactionFee_Eth"`
}

// GetDailyNetworkTxFee returns the amount of transaction fees paid to miners per day.
func (c *StatsClient) GetDailyNetworkTxFee(
	ctx context.Context, dates *ecommon.DateRange,
) (result []DailyTxFee, err error) {
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "dailytxnfee",
		Request: dates,
		Result:  &result,
	})
	return result, err
}

// DailyNewAddrCount describes the number of new Ethereum addresses created on
// a particular day.
type DailyNewAddrCount struct {
	Timestamp    time.Time `etherscan:"unixTimeStamp"`
	NewAddrCount uint64    `etherscan:"newAddressCount,num"`
}

// GetDailyNewAddrCount returns the number of new Ethereum addresses created per day.
func (c *StatsClient) GetDailyNewAddrCount(
	ctx context.Context, dates *ecommon.DateRange,
) (result []DailyNewAddrCount, err error) {
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "dailynewaddress",
		Request: dates,
		Result:  &result,
	})
	return result, err

}

// NetworkUtil describes the ethereum network utilization on a particular day.
type NetworkUtil struct {
	NetworkUtil decimal.Decimal `etherscan:"networkUtilization"`
	Timestamp   time.Time       `etherscan:"unixTimeStamp"`
}

// GetDailyNetworkUtil returns the daily average gas used over gas limit, in percentage.
func (c *StatsClient) GetDailyNetworkUtil(
	ctx context.Context, dates *ecommon.DateRange,
) (result []NetworkUtil, err error) {
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "dailynetutilization",
		Request: dates,
		Result:  &result,
	})
	return result, err
}

// NetworkHashRate describes the total processing power of the ethereum network
// on a particular day.
type NetworkHashRate struct {
	NetworkHashRate decimal.Decimal `etherscan:"networkHashRate"`
	Timestamp       time.Time       `etherscan:"unixTimeStamp"`
}

// GetDailyAvgHashRate returns the historical measure of processing power of the Ethereum network.
func (c *StatsClient) GetDailyAvgHashRate(
	ctx context.Context, dates *ecommon.DateRange,
) (result []NetworkHashRate, err error) {
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "dailyavghashrate",
		Request: dates,
		Result:  &result,
	})
	return result, err
}

// TxCount describes the total transaction count on a particular day.
type TxCount struct {
	Timestamp time.Time `etherscan:"unixTimeStamp"`
	TxCount   *big.Int  `etherscan:"transactionCount,num"`
}

// GetDailyTxCount returns the number of transactions performed on the Ethereum blockchain per day.
func (c *StatsClient) GetDailyTxCount(
	ctx context.Context, dates *ecommon.DateRange,
) (result []TxCount, err error) {
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "dailytx",
		Request: dates,
		Result:  &result,
	})
	return result, err
}

// NetDifficulty describes the mining difficulty on a particular day.
type NetDifficulty struct {
	Difficulty decimal.Decimal `etherscan:"networkDifficulty,comma"`
	Timestamp  time.Time       `etherscan:"unixTimeStamp"`
}

// GetDailyAvgNetDifficulty returns the historical mining difficulty of the Ethereum network.
func (c *StatsClient) GetDailyAvgNetDifficulty(
	ctx context.Context, dates *ecommon.DateRange,
) (result []NetDifficulty, err error) {
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "dailyavgnetdifficulty",
		Request: dates,
		Result:  &result,
	})
	return result, err
}

// HistoricalMarketCap describes the market cap of Ether on a particular day.
type HistoricalMarketCap struct {
	Timestamp time.Time `etherscan:"unixTimeStamp"`
	Supply    decimal.Decimal
	MarketCap decimal.Decimal `etherscan:"marketCap"`
	Price     decimal.Decimal
}

// GetETHHistoricalDailyMarketCap returns the historical Ether daily market capitalization.
func (c *StatsClient) GetETHHistoricalDailyMarketCap(
	ctx context.Context, dates *ecommon.DateRange,
) (result []HistoricalMarketCap, err error) {
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "ethdailymarketcap",
		Request: dates,
		Result:  &result,
	})
	return result, err
}

// HistoricalETHPrice describes the price of Ether on a particular day.
type HistoricalETHPrice struct {
	Timestamp time.Time       `etherscan:"unixTimeStamp"`
	USDValue  decimal.Decimal `etherscan:"value"`
}

// GetETHHistoricalPrice returns the historical price of 1 ETH.
func (c *StatsClient) GetETHHistoricalPrice(
	ctx context.Context, dates *ecommon.DateRange,
) (result []HistoricalETHPrice, err error) {
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "ethdailyprice",
		Request: dates,
		Result:  &result,
	})
	return result, err
}
