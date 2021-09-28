package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
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
	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "getblockreward",
		other:  map[string]string{"blockno": strconv.FormatUint(blockNumber, 10)},
	})
	if err != nil {
		return nil, err
	}

	result := new(BlockRewards)
	if err := unmarshalResponse(rspData, result); err != nil {
		return nil, err
	}

	return result, nil
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
	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "getblockcountdown",
		other:  map[string]string{"blockno": strconv.FormatUint(blockNumber, 10)},
	})
	if err != nil {
		return nil, err
	}

	result := new(BlockCountdown)
	if err := unmarshalResponse(rspData, result); err != nil {
		return nil, err
	}

	return result, nil
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
	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "getblocknobytime",
		other:  marshalRequest(req),
	})
	if err != nil {
		return 0, err
	}

	var result uintStr
	if err := json.Unmarshal(rspData, &result); err != nil {
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
) ([]AverageBlockSize, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "dailyavgblocksize",
		other:  marshalRequest(dates),
	})
	if err != nil {
		return nil, err
	}

	var result []AverageBlockSize
	if err := unmarshalResponse(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

type BlockCount struct {
	DailyBlockRewards
	BlockCount uint32 `etherscan:"blockCount,num"`
}

func (c *BlocksClient) GetDailyBlockCount(
	ctx context.Context, dates *DateRange,
) ([]BlockCount, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "dailyblkcount",
		other:  marshalRequest(dates),
	})
	if err != nil {
		return nil, err
	}

	var result []BlockCount
	if err := unmarshalResponse(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

type DailyBlockRewards struct {
	Timestamp       time.Time `etherscan:"unixTimeStamp"`
	BlockRewardsETH float64   `etherscan:"blockRewards_Eth"`
}

func (c *BlocksClient) GetDailyBlockRewards(
	ctx context.Context, dates *DateRange,
) ([]DailyBlockRewards, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "dailyblockrewards",
		other:  marshalRequest(dates),
	})
	if err != nil {
		return nil, err
	}

	var result []DailyBlockRewards
	if err := unmarshalResponse(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

type DailyBlockTime struct {
	Timestamp        time.Time `etherscan:"unixTimeStamp"`
	BlockTimeSeconds float64   `etherscan:"blockTime_sec"`
}

func (c *BlocksClient) GetDailyAverageBlockTime(
	ctx context.Context, dates *DateRange,
) ([]DailyBlockTime, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "dailyavgblocktime",
		other:  marshalRequest(dates),
	})
	if err != nil {
		return nil, err
	}

	var result []DailyBlockTime
	if err := unmarshalResponse(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

type DailyUnclesCount struct {
	Timestamp            time.Time `etherscan:"unixTimeStamp"`
	UncleBlockCount      uint32    `etherscan:"uncleBlockCount,num"`
	UncleBlockRewardsETH float64   `etherscan:"uncleBlockRewards_Eth"`
}

func (c *BlocksClient) GetDailyUnclesCount(
	ctx context.Context, dates *DateRange,
) ([]DailyUnclesCount, error) {
	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "dailyuncleblkcount",
		other:  marshalRequest(dates),
	})
	if err != nil {
		return nil, err
	}

	var result []DailyUnclesCount
	if err := unmarshalResponse(rspData, &result); err != nil {
		return nil, err
	}

	return result, nil
}
