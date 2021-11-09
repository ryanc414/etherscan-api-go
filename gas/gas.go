package gas

import (
	"context"
	"math/big"
	"time"

	ecommon "github.com/ryanc414/etherscan-api-go/common"
	"github.com/ryanc414/etherscan-api-go/httpapi"
	"github.com/ryanc414/etherscan-api-go/marshallers"
)

// GasClient is the client for gas actions.
type GasClient struct {
	API *httpapi.APIClient
}

// EstimateConfirmationTime returns the estimated time, in seconds, for a
// transaction to be confirmed on the blockchain.
func (c *GasClient) EstimateConfirmationTime(
	ctx context.Context, gasPriceGwei int64,
) (uint64, error) {
	// 1 gwei = 10^9 wei
	gasPriceWei := new(big.Int).Mul(
		big.NewInt(gasPriceGwei),
		new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil),
	)

	req := struct{ GasPrice *big.Int }{gasPriceWei}
	var result marshallers.UintStr

	err := c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.GasModule,
		Action:  "gasestimate",
		Request: &req,
		Result:  &result,
	})
	if err != nil {
		return 0, err
	}

	return result.Unwrap(), nil
}

// GasPrices describes the current recommended gas prices.
type GasPrices struct {
	LastBlock       uint64    `etherscan:"LastBlock"`
	SafeGasPrice    uint64    `etherscan:"SafeGasPrice"`
	ProposeGasPrice uint64    `etherscan:"ProposeGasPrice"`
	FastGasPrice    uint64    `etherscan:"FastGasPrice"`
	SuggestBaseFee  float64   `etherscan:"suggestBaseFee"`
	GasUsedRatio    []float64 `etherscan:"gasUsedRatio,sep"`
}

// GetGasOracle returns the current Safe, Proposed and Fast gas prices.
func (c *GasClient) GetGasOracle(ctx context.Context) (*GasPrices, error) {
	result := new(GasPrices)
	err := c.API.Call(ctx, &httpapi.CallParams{
		Module: ecommon.GasModule,
		Action: "gasoracle",
		Result: result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// AvgGasLimit describes the average gas limit on a particular day.
type AvgGasLimit struct {
	Timestamp time.Time `etherscan:"unixTimeStamp"`
	GasLimit  uint64    `etherscan:"gasLimit"`
}

// GetDailyAvgGasLimit returns the historical daily average gas limit of the Ethereum network.
func (c *GasClient) GetDailyAvgGasLimit(
	ctx context.Context, req *ecommon.DateRange,
) (result []AvgGasLimit, err error) {
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "dailyavggaslimit",
		Request: req,
		Result:  &result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GasUsed describes the total amount of gas used on a particular day.
type GasUsed struct {
	Timestamp time.Time `etherscan:"unixTimeStamp"`
	GasUsed   *big.Int  `etherscan:"gasUsed"`
}

// GetDailyTotalGasUsed returns the total amount of gas used daily for
// transctions on the Ethereum network.
func (c *GasClient) GetDailyTotalGasUsed(
	ctx context.Context, req *ecommon.DateRange,
) (result []GasUsed, err error) {
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "dailygasused",
		Request: req,
		Result:  &result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// AvgGasPrice describes the average gas prices on a particular day.
type AvgGasPrice struct {
	Timestamp   time.Time `etherscan:"unixTimeStamp"`
	MaxGasPrice *big.Int  `etherscan:"maxGasPrice_Wei"`
	MinGasPrice *big.Int  `etherscan:"minGasPrice_Wei"`
	AvgGasPrice *big.Int  `etherscan:"avgGasPrice_Wei"`
}

// GetDailyAvgGasPrice returns the daily average gas price used on the Ethereum network.
func (c *GasClient) GetDailyAvgGasPrice(
	ctx context.Context, req *ecommon.DateRange,
) (result []AvgGasPrice, err error) {
	err = c.API.Call(ctx, &httpapi.CallParams{
		Module:  ecommon.StatsModule,
		Action:  "dailyavggasprice",
		Request: req,
		Result:  &result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
