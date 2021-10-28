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

func TestBlocks(t *testing.T) {
	m := newMockServer("block", true)
	t.Cleanup(m.close)

	u, err := url.Parse(m.srv.URL)
	require.NoError(t, err)

	client := etherscan.New(&etherscan.Params{
		APIKey:  m.apiKey,
		BaseURL: u,
	})

	ctx := context.Background()

	t.Run("GetBlockRewards", func(t *testing.T) {
		rewards, err := client.Blocks.GetBlockRewards(ctx, 2165403)
		require.NoError(t, err)
		cupaloy.SnapshotT(t, rewards)
	})

	t.Run("GetBlockCountdown", func(t *testing.T) {
		countdown, err := client.Blocks.GetBlockCountdown(ctx, 16701588)
		require.NoError(t, err)
		cupaloy.SnapshotT(t, countdown)
	})

	t.Run("GetBlockNumber", func(t *testing.T) {
		number, err := client.Blocks.GetBlockNumber(ctx, &etherscan.BlockNumberRequest{
			Timestamp: time.Unix(1578638524, 0),
			Closest:   etherscan.ClosestAvailableBlockBefore,
		})
		require.NoError(t, err)
		assert.Equal(t, uint64(12712551), number)
	})

	dates := etherscan.DateRange{
		StartDate: time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2019, 2, 28, 0, 0, 0, 0, time.UTC),
		Sort:      etherscan.SortingPreferenceAsc,
	}

	t.Run("GetDailyAverageBlockSize", func(t *testing.T) {
		blockSizes, err := client.Blocks.GetDailyAverageBlockSize(ctx, &dates)
		require.NoError(t, err)

		require.Len(t, blockSizes, 2)
		cupaloy.SnapshotT(t, blockSizes)
	})

	t.Run("GetDailyBlockCount", func(t *testing.T) {
		blockCounts, err := client.Blocks.GetDailyBlockCount(ctx, &dates)
		require.NoError(t, err)
		require.Len(t, blockCounts, 2)

		cupaloy.SnapshotT(t, blockCounts)
	})

	t.Run("GetDailyBlockRewards", func(t *testing.T) {
		blockRewards, err := client.Blocks.GetDailyBlockRewards(ctx, &dates)
		require.NoError(t, err)
		require.Len(t, blockRewards, 2)

		cupaloy.SnapshotT(t, blockRewards)
	})

	t.Run("GetDailyAverageBlockTime", func(t *testing.T) {
		blockTimes, err := client.Blocks.GetDailyAverageBlockTime(ctx, &dates)
		require.NoError(t, err)
		require.Len(t, blockTimes, 2)

		cupaloy.SnapshotT(t, blockTimes)
	})

	t.Run("GetDailyUncleCount", func(t *testing.T) {
		unclesCounts, err := client.Blocks.GetDailyUnclesCount(ctx, &dates)
		require.NoError(t, err)
		require.Len(t, unclesCounts, 2)

		cupaloy.SnapshotT(t, unclesCounts)
	})
}
