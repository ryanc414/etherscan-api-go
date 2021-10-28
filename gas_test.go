package etherscan_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/ryanc414/etherscan-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGas(t *testing.T) {
	m := newMockServer("gastracker", false)
	t.Cleanup(m.close)

	u, err := url.Parse(m.srv.URL)
	require.NoError(t, err)

	client := etherscan.New(&etherscan.Params{
		APIKey:  m.apiKey,
		BaseURL: u,
	})

	ctx := context.Background()

	t.Run("EstimateConfirmationTime", func(t *testing.T) {
		confTime, err := client.Gas.EstimateConfirmationTime(ctx, 20)
		require.NoError(t, err)
		assert.Equal(t, uint64(9227), confTime)
	})

	t.Run("GetGasOracle", func(t *testing.T) {
		gas, err := client.Gas.GetGasOracle(ctx)
		require.NoError(t, err)
		cupaloy.SnapshotT(t, gas)
	})

	dateRange := etherscan.DateRange{
		StartDate: time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2019, 2, 28, 0, 0, 0, 0, time.UTC),
		Sort:      etherscan.SortingPreferenceAsc,
	}

	t.Run("GetDailyAvgGasLimit", func(t *testing.T) {
		avgGasLimits, err := client.Gas.GetDailyAvgGasLimit(ctx, &dateRange)
		require.NoError(t, err)
		require.Len(t, avgGasLimits, 3)
		cupaloy.SnapshotT(t, avgGasLimits)
	})

	t.Run("GetDailyTotalGasUsed", func(t *testing.T) {
		totalGasUsed, err := client.Gas.GetDailyTotalGasUsed(ctx, &dateRange)
		require.NoError(t, err)
		require.Len(t, totalGasUsed, 3)
		cupaloy.SnapshotT(t, totalGasUsed)
	})

	t.Run("GetDailyAvgGasPrice", func(t *testing.T) {
		gasPrices, err := client.Gas.GetDailyAvgGasPrice(ctx, &dateRange)
		require.NoError(t, err)
		require.Len(t, gasPrices, 3)
		cupaloy.SnapshotT(t, gasPrices)
	})
}
