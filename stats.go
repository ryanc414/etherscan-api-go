package etherscan

import (
	"context"
	"math/big"
	"time"

	"github.com/pkg/errors"
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
	ClientType ETHClientType
	SyncMode   NodeSyncMode
	Sort       SortingPreference
}

type ETHClientType int32

const (
	ETHClientTypeGeth ETHClientType = iota
	ETHClientTypeParity
)

func (t *ETHClientType) UnmarshalJSON(data []byte) error {
	switch s := string(data); s {
	case "\"Geth\"":
		*t = ETHClientTypeGeth
		return nil

	case "\"Parity\"":
		*t = ETHClientTypeParity
		return nil

	default:
		return errors.Errorf("unknown ETH client type %s", s)
	}
}

func (t ETHClientType) String() string {
	switch t {
	case ETHClientTypeGeth:
		return "geth"

	case ETHClientTypeParity:
		return "parity"

	default:
		panic("unknown ETH client type")
	}
}

type NodeSyncMode int32

const (
	NodeSyncModeDefault NodeSyncMode = iota
	NodeSyncModeArchive
)

func (m *NodeSyncMode) UnmarshalJSON(data []byte) error {
	switch s := string(data); s {
	case "\"Default\"":
		*m = NodeSyncModeDefault
		return nil

	case "\"Archive\"":
		*m = NodeSyncModeArchive
		return nil

	default:
		return errors.Errorf("unknown node sync mode %s", s)
	}
}

func (m NodeSyncMode) String() string {
	switch m {
	case NodeSyncModeDefault:
		return "default"

	case NodeSyncModeArchive:
		return "archive"

	default:
		panic("unknown node sync mode")
	}
}

type ETHNodeSize struct {
	BlockNumber    uint64        `etherscan:"blockNumber"`
	ChainTimestamp time.Time     `etherscan:"chainTimeStamp,date"`
	ChainSize      *big.Int      `etherscan:"chainSize"`
	ClientType     ETHClientType `etherscan:"clientType"`
	SyncMode       NodeSyncMode  `etherscan:"syncMode"`
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
	Timestamp time.Time       `etherscan:"unixTimestamp"`
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
