package etherscan

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const blocksModule = "block"

type BlocksClient struct {
	api *apiClient
}

type BlockRewards struct {
	BlockNumber          uint64         `etherscan:"blockNumber"`
	Timestamp            time.Time      `etherscan:"timeStamp"`
	BlockMiner           common.Address `etherscan:"blockMiner"`
	BlockReward          *big.Int       `etherscan:"blockReward"`
	Uncles               []UncleReward
	UncleInclusionReward *big.Int `etherscan:"uncleInclusionReward"`
}

type UncleReward struct {
	Miner         common.Address
	UnclePosition uint32   `etherscan:"unclePosition"`
	BlockReward   *big.Int `etherscan:"blockReward"`
}

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

type BlockCountdown struct {
	CurrentBlock      uint64  `etherscan:"CurrentBlock"`
	CountdownBlock    uint64  `etherscan:"CountdownBlock"`
	RemainingBlock    uint64  `etherscan:"RemainingBlock"`
	EstimateTimeInSec float64 `etherscan:"EstimateTimeInSec"`
}

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

type BlockNumberRequest struct {
	Timestamp time.Time
	Closest   ClosestAvailableBlock
}

type ClosestAvailableBlock int32

const (
	ClosestAvailableBlockBefore = iota
	ClosestAvailableBlockAfter
)

func (c ClosestAvailableBlock) String() string {
	switch c {
	case ClosestAvailableBlockBefore:
		return "before"

	case ClosestAvailableBlockAfter:
		return "after"

	default:
		panic(fmt.Sprintf("unknown closest available block parameter %d", int32(c)))
	}
}

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

type DateRange struct {
	StartDate time.Time `etherscan:"startdate,date"`
	EndDate   time.Time `etherscan:"enddate,date"`
	Sort      SortingPreference
}

type AverageBlockSize struct {
	Timestamp      time.Time `etherscan:"unixTimeStamp"`
	BlockSizeBytes uint32    `etherscan:"blockSize_bytes,num"`
}

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

type BlockCount struct {
	DailyBlockRewards
	BlockCount uint32 `etherscan:"blockCount,num"`
}

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

type DailyBlockRewards struct {
	Timestamp       time.Time `etherscan:"unixTimeStamp"`
	BlockRewardsETH float64   `etherscan:"blockRewards_Eth"`
}

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

type DailyBlockTime struct {
	Timestamp        time.Time `etherscan:"unixTimeStamp"`
	BlockTimeSeconds float64   `etherscan:"blockTime_sec"`
}

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

type DailyUnclesCount struct {
	Timestamp            time.Time `etherscan:"unixTimeStamp"`
	UncleBlockCount      uint32    `etherscan:"uncleBlockCount,num"`
	UncleBlockRewardsETH float64   `etherscan:"uncleBlockRewards_Eth"`
}

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
