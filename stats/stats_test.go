package stats_test

import (
	"context"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/ryanc414/etherscan-api-go"
	ecommon "github.com/ryanc414/etherscan-api-go/common"
	"github.com/ryanc414/etherscan-api-go/stats"
	"github.com/ryanc414/etherscan-api-go/testbed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStats(t *testing.T) {
	m := testbed.NewMockServer("stats", true)
	t.Cleanup(m.Close)

	u, err := m.URL()
	require.NoError(t, err)

	client := etherscan.New(&etherscan.Params{
		APIKey:  m.APIKey,
		BaseURL: u,
	})

	ctx := context.Background()

	t.Run("GetTotalETHSupply", func(t *testing.T) {
		supply, err := client.Stats.GetTotalETHSupply(ctx)
		require.NoError(t, err)
		assert.Equal(t, "116487067186500000000000000", supply.String())
	})

	t.Run("GetLastETHPrice", func(t *testing.T) {
		price, err := client.Stats.GetLastETHPrice(ctx)
		require.NoError(t, err)
		cupaloy.SnapshotT(t, price)
	})

	t.Run("GetEthereumNodesSize", func(t *testing.T) {
		nodes, err := client.Stats.GetEthereumNodesSize(ctx, &stats.NodesSizeReq{
			StartDate:  time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC),
			EndDate:    time.Date(2019, 2, 28, 0, 0, 0, 0, time.UTC),
			ClientType: stats.ETHClientTypeReqGeth,
			SyncMode:   stats.NodeSyncModeReqDefault,
			Sort:       ecommon.SortingPreferenceAsc,
		})
		require.NoError(t, err)
		require.Len(t, nodes, 2)
		cupaloy.SnapshotT(t, nodes)
	})

	t.Run("GetToalNodesCount", func(t *testing.T) {
		count, err := client.Stats.GetTotalNodesCount(ctx)
		require.NoError(t, err)
		cupaloy.SnapshotT(t, count)
	})

	dates := ecommon.DateRange{
		StartDate: time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2019, 2, 28, 0, 0, 0, 0, time.UTC),
		Sort:      ecommon.SortingPreferenceAsc,
	}

	t.Run("GetDailyNetworkTxFee", func(t *testing.T) {
		fees, err := client.Stats.GetDailyNetworkTxFee(ctx, &dates)
		require.NoError(t, err)
		require.Len(t, fees, 2)
		cupaloy.SnapshotT(t, fees)
	})

	t.Run("GetDailyNewAddrCount", func(t *testing.T) {
		count, err := client.Stats.GetDailyNewAddrCount(ctx, &dates)
		require.NoError(t, err)
		require.Len(t, count, 2)
		cupaloy.SnapshotT(t, count)
	})

	t.Run("GetDailyNetworkUtil", func(t *testing.T) {
		util, err := client.Stats.GetDailyNetworkUtil(ctx, &dates)
		require.NoError(t, err)
		require.Len(t, util, 2)
		cupaloy.SnapshotT(t, util)
	})

	t.Run("GetDailyAvgNetHashRate", func(t *testing.T) {
		hashRate, err := client.Stats.GetDailyAvgHashRate(ctx, &dates)
		require.NoError(t, err)
		require.Len(t, hashRate, 2)
		cupaloy.SnapshotT(t, hashRate)
	})

	t.Run("GetDailyTxCount", func(t *testing.T) {
		count, err := client.Stats.GetDailyTxCount(ctx, &dates)
		require.NoError(t, err)
		require.Len(t, count, 2)
		cupaloy.SnapshotT(t, count)
	})

	t.Run("GetDailyAvgNetDifficulty", func(t *testing.T) {
		difficulty, err := client.Stats.GetDailyAvgNetDifficulty(ctx, &dates)
		require.NoError(t, err)
		require.Len(t, difficulty, 2)
		cupaloy.SnapshotT(t, difficulty)
	})

	t.Run("GetETHHistoricalDailyMarketCap", func(t *testing.T) {
		marketCap, err := client.Stats.GetETHHistoricalDailyMarketCap(ctx, &dates)
		require.NoError(t, err)
		require.Len(t, marketCap, 2)
		cupaloy.SnapshotT(t, marketCap)
	})

	t.Run("GetETHHistoricalPrice", func(t *testing.T) {
		price, err := client.Stats.GetETHHistoricalPrice(ctx, &dates)
		require.NoError(t, err)
		require.Len(t, price, 2)
		cupaloy.SnapshotT(t, price)
	})
}
