//go:generate go-enum -f=$GOFILE
package etherscan

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const blocksModule = "block"

// BlocksClient is the client for blocks related actions.
type BlocksClient struct {
	api *apiClient
}

// BlockRewards contains information on a block's mining reward.
type BlockRewards struct {
	BlockNumber          uint64         `etherscan:"blockNumber"`
	Timestamp            time.Time      `etherscan:"timeStamp"`
	BlockMiner           common.Address `etherscan:"blockMiner"`
	BlockReward          *big.Int       `etherscan:"blockReward"`
	Uncles               []UncleReward
	UncleInclusionReward *big.Int `etherscan:"uncleInclusionReward"`
}

// UncleReward contains information on a block uncle's mining reward.
type UncleReward struct {
	Miner         common.Address
	UnclePosition uint32   `etherscan:"unclePosition"`
	BlockReward   *big.Int `etherscan:"blockreward"`
}

// GetBlockRewards returns the block reward and 'Uncle' block rewards.
func (c *BlocksClient) GetBlockRewards(
	ctx context.Context, blockNumber uint64,
) (*BlockRewards, error) {
	result := new(BlockRewards)
	req := struct{ Blockno uint64 }{blockNumber}

	err := c.api.call(ctx, &callParams{
		module:  blocksModule,
		action:  "getblockreward",
		request: req,
		result:  result,
	})

	return result, err
}

// BlockCountdown contains information on the estimated time until a block is mined.
type BlockCountdown struct {
	CurrentBlock      uint64  `etherscan:"CurrentBlock"`
	CountdownBlock    uint64  `etherscan:"CountdownBlock"`
	RemainingBlock    uint64  `etherscan:"RemainingBlock"`
	EstimateTimeInSec float64 `etherscan:"EstimateTimeInSec"`
}

// GetBlockCountdown returns the estimated time remaining, in seconds, until a certain block is mined.
func (c *BlocksClient) GetBlockCountdown(
	ctx context.Context, blockNumber uint64,
) (*BlockCountdown, error) {
	result := new(BlockCountdown)
	req := struct{ Blockno uint64 }{blockNumber}

	err := c.api.call(ctx, &callParams{
		module:  blocksModule,
		action:  "getblockcountdown",
		request: req,
		result:  result,
	})

	return result, err
}

// BlockNumberRequest contains the request parameters for GetBlockNumber.
type BlockNumberRequest struct {
	Timestamp time.Time
	Closest   ClosestAvailableBlock
}

// ClosestAvailableBlock is an enumaration of the closest available block
// parameters.
// ENUM(before,after)
type ClosestAvailableBlock int32

// GetBlockNumber returns the block number that was mined at a certain timestamp.
func (c *BlocksClient) GetBlockNumber(
	ctx context.Context, req *BlockNumberRequest,
) (uint64, error) {
	var result uintStr

	err := c.api.call(ctx, &callParams{
		module:  blocksModule,
		action:  "getblocknobytime",
		request: req,
		result:  &result,
	})
	if err != nil {
		return 0, err
	}

	return result.unwrap(), nil
}

// DateRange contains request parameters for requests that span a set of dates.
type DateRange struct {
	StartDate time.Time `etherscan:"startdate,date"`
	EndDate   time.Time `etherscan:"enddate,date"`
	Sort      SortingPreference
}

// AverageBlockSize contains information on the average size of a block on a given day.
type AverageBlockSize struct {
	Timestamp      time.Time `etherscan:"unixTimeStamp"`
	BlockSizeBytes uint32    `etherscan:"blockSize_bytes,num"`
}

// GetDailyAverageBlockSize returns the daily average block size within a date range.
func (c *BlocksClient) GetDailyAverageBlockSize(
	ctx context.Context, dates *DateRange,
) (result []AverageBlockSize, err error) {
	err = c.api.call(ctx, &callParams{
		module:  blocksModule,
		action:  "dailyavgblocksize",
		request: dates,
		result:  &result,
	})

	return result, err
}

// BlockCount contains information on the block count on a particular day.
type BlockCount struct {
	DailyBlockRewards
	BlockCount uint32 `etherscan:"blockCount,num"`
}

// GetDailyBlockCount returns the number of blocks mined daily and the amount of block rewards.
func (c *BlocksClient) GetDailyBlockCount(
	ctx context.Context, dates *DateRange,
) (result []BlockCount, err error) {
	err = c.api.call(ctx, &callParams{
		module:  blocksModule,
		action:  "dailyblkcount",
		request: dates,
		result:  &result,
	})

	return result, err
}

// DailyBlockRewards contains information on the total block rewards distributed
// to miners on a particular day.
type DailyBlockRewards struct {
	Timestamp       time.Time `etherscan:"unixTimeStamp"`
	BlockRewardsETH float64   `etherscan:"blockRewards_Eth"`
}

// GetDailyBlockRewards returns the amount of block rewards distributed to miners daily.
func (c *BlocksClient) GetDailyBlockRewards(
	ctx context.Context, dates *DateRange,
) (result []DailyBlockRewards, err error) {
	err = c.api.call(ctx, &callParams{
		module:  blocksModule,
		action:  "dailyblockrewards",
		request: dates,
		result:  &result,
	})

	return result, err
}

// DailyBlockTime contains information on the average time to mine a block on a
// particular day.
type DailyBlockTime struct {
	Timestamp        time.Time `etherscan:"unixTimeStamp"`
	BlockTimeSeconds float64   `etherscan:"blockTime_sec"`
}

// GetDailyAverageBlockTime returns the daily average of time needed for a block
// to be successfully mined.
func (c *BlocksClient) GetDailyAverageBlockTime(
	ctx context.Context, dates *DateRange,
) (result []DailyBlockTime, err error) {
	err = c.api.call(ctx, &callParams{
		module:  blocksModule,
		action:  "dailyavgblocktime",
		request: dates,
		result:  &result,
	})

	return result, err
}

// DailyUnclesCount contains information on uncle blocks mined in a particular
// day.
type DailyUnclesCount struct {
	Timestamp            time.Time `etherscan:"unixTimeStamp"`
	UncleBlockCount      uint32    `etherscan:"uncleBlockCount,num"`
	UncleBlockRewardsETH float64   `etherscan:"uncleBlockRewards_Eth"`
}

// GetDailyUnclesCount returns the number of 'Uncle' blocks mined daily and the
// amount of 'Uncle' block rewards.
func (c *BlocksClient) GetDailyUnclesCount(
	ctx context.Context, dates *DateRange,
) (result []DailyUnclesCount, err error) {
	err = c.api.call(ctx, &callParams{
		module:  blocksModule,
		action:  "dailyuncleblkcount",
		request: dates,
		result:  &result,
	})

	return result, err
}
