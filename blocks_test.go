package etherscan_test

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/ryanc414/etherscan-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlocks(t *testing.T) {
	m := newMockBlocksAPI()

	ts := httptest.NewServer(m)
	t.Cleanup(ts.Close)

	u, err := url.Parse(ts.URL)
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

type mockBlocksAPI struct {
	apiKey string
}

func newMockBlocksAPI() mockBlocksAPI {
	return mockBlocksAPI{apiKey: uuid.NewString()}
}

func (m mockBlocksAPI) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/api" {
		http.Error(w, "path not found", http.StatusNotFound)
		return
	}

	q := req.URL.Query()
	if q.Get("module") != "block" {
		http.Error(w, "unknown module", http.StatusNotFound)
		return
	}

	if q.Get("apikey") != m.apiKey {
		http.Error(w, "unknown API key", http.StatusForbidden)
		return
	}

	switch q.Get("action") {
	case "getblockreward":
		m.handleGetBlockReward(w, q)

	case "getblockcountdown":
		m.handleGetBlockCountdown(w, q)

	case "getblocknobytime":
		m.handleGetBlockNumber(w, q)

	case "dailyavgblocksize":
		m.handleDailyAvgBlocksize(w, q)

	case "dailyblkcount":
		m.handleDailyBlockCount(w, q)

	case "dailyblockrewards":
		m.handleDailyBlockRewards(w, q)

	case "dailyavgblocktime":
		m.handleDailyAvgBlockTime(w, q)

	case "dailyuncleblkcount":
		m.handleDailyUncleBlockCount(w, q)

	default:
		http.Error(w, "unknown action", http.StatusNotFound)
	}
}

const getBlockRewardResponse = `{
	"status":"1",
	"message":"OK",
	"result":{
	   "blockNumber":"2165403",
	   "timeStamp":"1472533979",
	   "blockMiner":"0x13a06d3dfe21e0db5c016c03ea7d2509f7f8d1e3",
	   "blockReward":"5314181600000000000",
	   "uncles":[
		  {
			 "miner":"0xbcdfc35b86bedf72f0cda046a3c16829a2ef41d1",
			 "unclePosition":"0",
			 "blockreward":"3750000000000000000"
		  },
		  {
			 "miner":"0x0d0c9855c722ff0c78f21e43aa275a5b8ea60dce",
			 "unclePosition":"1",
			 "blockreward":"3750000000000000000"
		  }
	   ],
	   "uncleInclusionReward":"312500000000000000"
	}
}`

func (m mockBlocksAPI) handleGetBlockReward(w http.ResponseWriter, q url.Values) {
	if q.Get("blockno") != "2165403" {
		http.Error(w, "unexpected block number", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(getBlockRewardResponse))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const getBlockCountdownResponse = `{
	"status":"1",
	"message":"OK",
	"result":{
	   "CurrentBlock":"12715477",
	   "CountdownBlock":"16701588",
	   "RemainingBlock":"3986111",
	   "EstimateTimeInSec":"52616680.2"
	}
}`

func (m mockBlocksAPI) handleGetBlockCountdown(w http.ResponseWriter, q url.Values) {
	if q.Get("blockno") != "16701588" {
		http.Error(w, "unexpected block number", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(getBlockCountdownResponse))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const getBlockNumberResponse = `{
	"status":"1",
	"message":"OK",
	"result":"12712551"
}`

func (m mockBlocksAPI) handleGetBlockNumber(w http.ResponseWriter, q url.Values) {
	if q.Get("timestamp") != "1578638524" {
		http.Error(w, "unexpected timestamp", http.StatusBadRequest)
		return
	}

	if q.Get("closest") != "before" {
		http.Error(w, "unexpected closest value", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(getBlockNumberResponse))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const dailyAverageBlockSize = `{
	"status":"1",
	"message":"OK",
	"result":[
	   {
		  "UTCDate":"2019-02-01",
		  "unixTimeStamp":"1548979200",
		  "blockSize_bytes":20373
	   },
	   {
		  "UTCDate":"2019-02-28",
		  "unixTimeStamp":"1551312000",
		  "blockSize_bytes":25117
	   }
	]
}`

func (m mockBlocksAPI) handleDailyAvgBlocksize(w http.ResponseWriter, q url.Values) {
	if q.Get("startdate") != "2019-02-01" {
		http.Error(w, "unexpected start date", http.StatusBadRequest)
		return
	}

	if q.Get("enddate") != "2019-02-28" {
		http.Error(w, "unexpected end date", http.StatusBadRequest)
		return
	}

	if q.Get("sort") != "asc" {
		http.Error(w, "unexpected sort param", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(dailyAverageBlockSize))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const dailyBlockCountResponse = `{
	"status":"1",
	"message":"OK",
	"result":[
	   {
		  "UTCDate":"2019-02-01",
		  "unixTimeStamp":"1548979200",
		  "blockCount":4848,
		  "blockRewards_Eth":"14929.464690870590355682"
	   },
	   {
		  "UTCDate":"2019-02-28",
		  "unixTimeStamp":"1551312000",
		  "blockCount":4366,
		  "blockRewards_Eth":"12808.485512162356907132"
	   }
	]
}`

func (m mockBlocksAPI) handleDailyBlockCount(w http.ResponseWriter, q url.Values) {
	if q.Get("startdate") != "2019-02-01" {
		http.Error(w, "unexpected start date", http.StatusBadRequest)
		return
	}

	if q.Get("enddate") != "2019-02-28" {
		http.Error(w, "unexpected end date", http.StatusBadRequest)
		return
	}

	if q.Get("sort") != "asc" {
		http.Error(w, "unexpected sort param", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(dailyBlockCountResponse))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const dailyBlockRewardsResponse = `{
	"status":"1",
	"message":"OK",
	"result":[
	   {
		  "UTCDate":"2019-02-01",
		  "unixTimeStamp":"1548979200",
		  "blockRewards_Eth":"15300.65625"
	   },
	   {
		  "UTCDate":"2019-02-28",
		  "unixTimeStamp":"1551312000",
		  "blockRewards_Eth":"12954.84375"
	   }
	]
}`

func (m mockBlocksAPI) handleDailyBlockRewards(w http.ResponseWriter, q url.Values) {
	if q.Get("startdate") != "2019-02-01" {
		http.Error(w, "unexpected start date", http.StatusBadRequest)
		return
	}

	if q.Get("enddate") != "2019-02-28" {
		http.Error(w, "unexpected end date", http.StatusBadRequest)
		return
	}

	if q.Get("sort") != "asc" {
		http.Error(w, "unexpected sort param", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(dailyBlockRewardsResponse))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const dailyAvgBlocktimeResponse = `{
	"status":"1",
	"message":"OK",
	"result":[
	   {
		  "UTCDate":"2019-02-01",
		  "unixTimeStamp":"1548979200",
		  "blockTime_sec":"17.67"
	   },
	   {
		  "UTCDate":"2019-02-28",
		  "unixTimeStamp":"1551312000",
		  "blockTime_sec":"19.61"
	   }
	]
}`

func (m mockBlocksAPI) handleDailyAvgBlockTime(w http.ResponseWriter, q url.Values) {
	if q.Get("startdate") != "2019-02-01" {
		http.Error(w, "unexpected start date", http.StatusBadRequest)
		return
	}

	if q.Get("enddate") != "2019-02-28" {
		http.Error(w, "unexpected end date", http.StatusBadRequest)
		return
	}

	if q.Get("sort") != "asc" {
		http.Error(w, "unexpected sort param", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(dailyAvgBlocktimeResponse))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const dailyUncleBlockCountResponse = `{
	"status":"1",
	"message":"OK",
	"result":[
	   {
		  "UTCDate":"2019-02-01",
		  "unixTimeStamp":"1548979200",
		  "uncleBlockCount":287,
		  "uncleBlockRewards_Eth":"729.75"
	   },
	   {
		  "UTCDate":"2019-02-28",
		  "unixTimeStamp":"1551312000",
		  "uncleBlockCount":288,
		  "uncleBlockRewards_Eth":"691.5"
	   }
	]
}`

func (m mockBlocksAPI) handleDailyUncleBlockCount(w http.ResponseWriter, q url.Values) {
	if q.Get("startdate") != "2019-02-01" {
		http.Error(w, "unexpected start date", http.StatusBadRequest)
		return
	}

	if q.Get("enddate") != "2019-02-28" {
		http.Error(w, "unexpected end date", http.StatusBadRequest)
		return
	}

	if q.Get("sort") != "asc" {
		http.Error(w, "unexpected sort param", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(dailyUncleBlockCountResponse))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
