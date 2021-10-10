package etherscan_test

import (
	"context"
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ryanc414/etherscan-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlocks(t *testing.T) {
	m := newMockServer("block")
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

		assert.Equal(t, uint64(2165403), rewards.BlockNumber)
		assert.Equal(t, time.Unix(1472533979, 0), rewards.Timestamp)
		assert.Equal(
			t,
			common.HexToAddress("0x13a06d3dfe21e0db5c016c03ea7d2509f7f8d1e3"),
			rewards.BlockMiner,
		)

		expectedBlockReward, ok := new(big.Int).SetString("5314181600000000000", 10)
		require.True(t, ok)
		assert.Equal(t, 0, expectedBlockReward.Cmp(rewards.BlockReward))

		expectedUnclesReward, ok := new(big.Int).SetString("312500000000000000", 10)
		require.True(t, ok)
		assert.Equal(t, 0, expectedUnclesReward.Cmp(rewards.UncleInclusionReward))

		require.Len(t, rewards.Uncles, 2)
	})

	t.Run("GetBlockCountdown", func(t *testing.T) {
		countdown, err := client.Blocks.GetBlockCountdown(ctx, 16701588)
		require.NoError(t, err)

		assert.Equal(t, uint64(12715477), countdown.CurrentBlock)
		assert.Equal(t, uint64(16701588), countdown.CountdownBlock)
		assert.Equal(t, uint64(3986111), countdown.RemainingBlock)
		assert.Equal(t, float64(52616680.2), countdown.EstimateTimeInSec)
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
		Sort:      etherscan.SortingPreferenceAscending,
	}

	t.Run("GetDailyAverageBlockSize", func(t *testing.T) {
		blockSizes, err := client.Blocks.GetDailyAverageBlockSize(ctx, &dates)
		require.NoError(t, err)

		require.Len(t, blockSizes, 2)

		assert.Equal(t, time.Unix(1548979200, 0), blockSizes[0].Timestamp)
		assert.Equal(t, uint32(20373), blockSizes[0].BlockSizeBytes)

		assert.Equal(t, time.Unix(1551312000, 0), blockSizes[1].Timestamp)
		assert.Equal(t, uint32(25117), blockSizes[1].BlockSizeBytes)
	})

	t.Run("GetDailyBlockCount", func(t *testing.T) {
		blockCounts, err := client.Blocks.GetDailyBlockCount(ctx, &dates)
		require.NoError(t, err)
		require.Len(t, blockCounts, 2)

		assert.Equal(t, time.Unix(1548979200, 0), blockCounts[0].Timestamp)
		assert.Equal(t, uint32(4848), blockCounts[0].BlockCount)
		assert.Equal(t, float64(14929.464690870590355682), blockCounts[0].BlockRewardsETH)

		assert.Equal(t, time.Unix(1551312000, 0), blockCounts[1].Timestamp)
		assert.Equal(t, uint32(4366), blockCounts[1].BlockCount)
		assert.Equal(t, float64(12808.485512162356907132), blockCounts[1].BlockRewardsETH)
	})

	t.Run("GetDailyBlockRewards", func(t *testing.T) {
		blockRewards, err := client.Blocks.GetDailyBlockRewards(ctx, &dates)
		require.NoError(t, err)

		require.Len(t, blockRewards, 2)

		assert.Equal(t, time.Unix(1548979200, 0), blockRewards[0].Timestamp)
		assert.Equal(t, float64(15300.65625), blockRewards[0].BlockRewardsETH)

		assert.Equal(t, time.Unix(1551312000, 0), blockRewards[1].Timestamp)
		assert.Equal(t, float64(12954.84375), blockRewards[1].BlockRewardsETH)
	})

	t.Run("GetDailyAverageBlockTime", func(t *testing.T) {
		blockTimes, err := client.Blocks.GetDailyAverageBlockTime(ctx, &dates)
		require.NoError(t, err)

		require.Len(t, blockTimes, 2)

		assert.Equal(t, time.Unix(1548979200, 0), blockTimes[0].Timestamp)
		assert.Equal(t, float64(17.67), blockTimes[0].BlockTimeSeconds)

		assert.Equal(t, time.Unix(1551312000, 0), blockTimes[1].Timestamp)
		assert.Equal(t, float64(19.61), blockTimes[1].BlockTimeSeconds)
	})

	t.Run("GetDailyUncleCount", func(t *testing.T) {
		unclesCounts, err := client.Blocks.GetDailyUnclesCount(ctx, &dates)
		require.NoError(t, err)

		require.Len(t, unclesCounts, 2)

		assert.Equal(t, time.Unix(1548979200, 0), unclesCounts[0].Timestamp)
		assert.Equal(t, uint32(287), unclesCounts[0].UncleBlockCount)
		assert.Equal(t, float64(729.75), unclesCounts[0].UncleBlockRewardsETH)

		assert.Equal(t, time.Unix(1551312000, 0), unclesCounts[1].Timestamp)
		assert.Equal(t, uint32(288), unclesCounts[1].UncleBlockCount)
		assert.Equal(t, float64(691.5), unclesCounts[1].UncleBlockRewardsETH)
	})
}
