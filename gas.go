package etherscan

import (
	"context"
	"math/big"
	"time"
)

const gasModule = "gastracker"

type GasClient struct {
	api *apiClient
}

func (c *GasClient) EstimateConfirmationTime(
	ctx context.Context, gasPriceGwei int64,
) (uint64, error) {
	// 1 gwei = 10^9 wei
	gasPriceWei := new(big.Int).Mul(
		big.NewInt(gasPriceGwei),
		new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil),
	)

	req := struct{ GasPrice *big.Int }{gasPriceWei}
	var result uintStr

	err := c.api.call(ctx, &callParams{
		module:  gasModule,
		action:  "gasestimate",
		request: &req,
		result:  &result,
	})
	if err != nil {
		return 0, err
	}

	return result.unwrap(), nil
}

type GasPrices struct {
	LastBlock       uint64    `etherscan:"LastBlock"`
	SafeGasPrice    uint64    `etherscan:"SafeGasPrice"`
	ProposeGasPrice uint64    `etherscan:"ProposeGasPrice"`
	FastGasPrice    uint64    `etherscan:"FastGasPrice"`
	SuggestBaseFee  float64   `etherscan:"suggestBaseFee"`
	GasUsedRatio    []float64 `etherscan:"gasUsedRatio,sep"`
}

func (c *GasClient) GetGasOracle(ctx context.Context) (*GasPrices, error) {
	result := new(GasPrices)
	err := c.api.call(ctx, &callParams{
		module: gasModule,
		action: "gasoracle",
		result: result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

type AvgGasLimit struct {
	Timestamp time.Time `etherscan:"unixTimeStamp"`
	GasLimit  uint64    `etherscan:"gasLimit"`
}

func (c *GasClient) GetDailyAvgGasLimit(
	ctx context.Context, req *DateRange,
) (result []AvgGasLimit, err error) {
	err = c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "dailyavggaslimit",
		request: req,
		result:  &result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

type GasUsed struct {
	Timestamp time.Time `etherscan:"unixTimeStamp"`
	GasUsed   *big.Int  `etherscan:"gasUsed"`
}

func (c *GasClient) GetDailyTotalGasUsed(
	ctx context.Context, req *DateRange,
) (result []GasUsed, err error) {
	err = c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "dailygasused",
		request: req,
		result:  &result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

type AvgGasPrice struct {
	Timestamp   time.Time `etherscan:"unixTimeStamp"`
	MaxGasPrice *big.Int  `etherscan:"maxGasPrice_Wei"`
	MinGasPrice *big.Int  `etherscan:"minGasPrice_Wei"`
	AvgGasPrice *big.Int  `etherscan:"avgGasPrice_Wei"`
}

func (c *GasClient) GetDailyAvgGasPrice(
	ctx context.Context, req *DateRange,
) (result []AvgGasPrice, err error) {
	err = c.api.call(ctx, &callParams{
		module:  statsModule,
		action:  "dailyavggasprice",
		request: req,
		result:  &result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
