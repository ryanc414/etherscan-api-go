//go:generate go-enum -f=$GOFILE --marshal
package etherscan

import (
	"context"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

type StatsClient struct {
	api *apiClient
}

const statsModule = "stats"

func (c *StatsClient) GetTotalETHSupply(ctx context.Context) (*big.Int, error) {
	result := new(bigInt)
	err := c.api.call(ctx, &callParams{
		module: statsModule,
		action: "ethsupply",
		result: result,
	})
	if err != nil {
		return nil, err
	}

	return result.unwrap(), nil
}

type ETHPrice struct {
	ETHBTC          decimal.Decimal
	ETHBTCTimestamp time.Time `etherscan:"ethbtc_timestamp"`
	ETHUSD          decimal.Decimal
	ETHUSDTimestamp time.Time `etherscan:"ethusd_timestamp"`
}

func (c *StatsClient) GetLastETHPrice(ctx context.Context) (*ETHPrice, error) {
	result := new(ETHPrice)
	err := c.api.call(ctx, &callParams{
		module: statsModule,
		action: "ethprice",
		result: result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

type NodesSizeReq struct {
	StartDate  time.Time `etherscan:"startdate,date"`
	EndDate    time.Time `etherscan:"enddate,date"`
	ClientType ETHClientTypeReq
	SyncMode   NodeSyncModeReq
	Sort       SortingPreference
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

type ETHNodeSize struct {
	BlockNumber    uint64              `etherscan:"blockNumber"`
	ChainTimestamp time.Time           `etherscan:"chainTimeStamp,date"`
	ChainSize      *big.Int            `etherscan:"chainSize"`
	ClientType     ETHClientTypeResult `etherscan:"clientType"`
	SyncMode       NodeSyncModeResult  `etherscan:"syncMode"`
}

func (c *StatsClient) GetEthereumNodesSize(
	ctx context.Context, req *NodesSizeReq,
) (result []ETHNodeSize, err error) {
	err = c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "chainsize",
		request: req,
		result:  &result,
	})
	return result, err
}

type NodeCount struct {
	Date           time.Time `etherscan:"UTCDate,date"`
	TotalNodeCount uint64    `etherscan:"TotalNodeCount"`
}

func (c *StatsClient) GetTotalNodesCount(ctx context.Context) (*NodeCount, error) {
	result := new(NodeCount)
	err := c.api.call(ctx, &callParams{
		module: statsModule,
		action: "nodecount",
		result: result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

type DailyTxFee struct {
	Timestamp time.Time       `etherscan:"unixTimeStamp"`
	TxFeeETH  decimal.Decimal `etherscan:"transactionFee_Eth"`
}

func (c *StatsClient) GetDailyNetworkTxFee(
	ctx context.Context, dates *DateRange,
) (result []DailyTxFee, err error) {
	err = c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "dailytxnfee",
		request: dates,
		result:  &result,
	})
	return result, err
}

type DailyNewAddrCount struct {
	Timestamp    time.Time `etherscan:"unixTimeStamp"`
	NewAddrCount uint64    `etherscan:"newAddressCount,num"`
}

func (c *StatsClient) GetDailyNewAddrCount(
	ctx context.Context, dates *DateRange,
) (result []DailyNewAddrCount, err error) {
	err = c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "dailynewaddress",
		request: dates,
		result:  &result,
	})
	return result, err

}

type NetworkUtil struct {
	NetworkUtil decimal.Decimal `etherscan:"networkUtilization"`
	Timestamp   time.Time       `etherscan:"unixTimeStamp"`
}

func (c *StatsClient) GetDailyNetworkUtil(
	ctx context.Context, dates *DateRange,
) (result []NetworkUtil, err error) {
	err = c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "dailynetutilization",
		request: dates,
		result:  &result,
	})
	return result, err
}

type NetworkHashRate struct {
	NetworkHashRate decimal.Decimal `etherscan:"networkHashRate"`
	Timestamp       time.Time       `etherscan:"unixTimeStamp"`
}

func (c *StatsClient) GetDailyAvgHashRate(
	ctx context.Context, dates *DateRange,
) (result []NetworkHashRate, err error) {
	err = c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "dailyavghashrate",
		request: dates,
		result:  &result,
	})
	return result, err
}

type TxCount struct {
	Timestamp time.Time `etherscan:"unixTimeStamp"`
	TxCount   *big.Int  `etherscan:"transactionCount,num"`
}

func (c *StatsClient) GetDailyTxCount(
	ctx context.Context, dates *DateRange,
) (result []TxCount, err error) {
	err = c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "dailytx",
		request: dates,
		result:  &result,
	})
	return result, err
}

type NetDifficulty struct {
	Difficulty decimal.Decimal `etherscan:"networkDifficulty,comma"`
	Timestamp  time.Time       `etherscan:"unixTimeStamp"`
}

func (c *StatsClient) GetDailyAvgNetDifficulty(
	ctx context.Context, dates *DateRange,
) (result []NetDifficulty, err error) {
	err = c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "dailyavgnetdifficulty",
		request: dates,
		result:  &result,
	})
	return result, err
}

type HistoricalMarketCap struct {
	Timestamp time.Time `etherscan:"unixTimeStamp"`
	Supply    decimal.Decimal
	MarketCap decimal.Decimal `etherscan:"marketCap"`
	Price     decimal.Decimal
}

func (c *StatsClient) GetETHHistoricalDailyMarketCap(
	ctx context.Context, dates *DateRange,
) (result []HistoricalMarketCap, err error) {
	err = c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "ethdailymarketcap",
		request: dates,
		result:  &result,
	})
	return result, err
}

type HistoricalETHPrice struct {
	Timestamp time.Time       `etherscan:"unixTimeStamp"`
	USDValue  decimal.Decimal `etherscan:"value"`
}

func (c *StatsClient) GetETHHistoricalPrice(
	ctx context.Context, dates *DateRange,
) (result []HistoricalETHPrice, err error) {
	err = c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "ethdailyprice",
		request: dates,
		result:  &result,
	})
	return result, err
}
