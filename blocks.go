package etherscan

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const blocksModule = "block"

type BlocksClient struct {
	api *apiClient
}

type BlockRewards struct {
	BlockNumber          uint64
	Timestamp            time.Time
	BlockMiner           common.Address
	BlockReward          *big.Int
	Uncles               []UncleReward
	UncleInclusionReward *big.Int
}

type UncleReward struct {
	Miner         common.Address
	UnclePosition uint32
	BlockReward   *big.Int
}

type blockRewardsResult struct {
	BlockNumber           uintStr       `json:"blockNumber"`
	Timestamp             unixTimestamp `json:"timeStamp"`
	BlockMiner            string        `json:"blockMiner"`
	BlockReward           bigInt        `json:"blockReward"`
	Uncles                []uncleResult `json:"uncles"`
	UnclesInclusionReward bigInt        `json:"uncleInclusionReward"`
}

type uncleResult struct {
	Miner       string  `json:"miner"`
	Position    uintStr `json:"unclePosition"`
	BlockReward bigInt  `json:"blockreward"`
}

func (u *uncleResult) toUncle() *UncleReward {
	return &UncleReward{
		Miner:         common.HexToAddress(u.Miner),
		UnclePosition: uint32(u.Position.unwrap()),
		BlockReward:   u.BlockReward.unwrap(),
	}
}

func (r blockRewardsResult) toRewards() *BlockRewards {
	uncles := make([]UncleReward, len(r.Uncles))
	for i := range r.Uncles {
		uncles[i] = *r.Uncles[i].toUncle()
	}

	return &BlockRewards{
		BlockNumber:          r.BlockNumber.unwrap(),
		Timestamp:            r.Timestamp.unwrap(),
		BlockMiner:           common.HexToAddress(r.BlockMiner),
		BlockReward:          r.BlockReward.unwrap(),
		Uncles:               uncles,
		UncleInclusionReward: r.UnclesInclusionReward.unwrap(),
	}
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

	var result blockRewardsResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result.toRewards(), nil
}

type BlockCountdown struct {
	CurrentBlock      uint64
	CountdownBlock    uint64
	RemainingBlock    uint64
	EstimateTimeInSec float64
}

type blockCountdownResult struct {
	CurrentBlock      uintStr  `json:"CurrentBlock"`
	CountdownBlock    uintStr  `json:"CountdownBlock"`
	RemainingBlock    uintStr  `json:"RemainingBlock"`
	EstimateTimeInSec floatStr `json:"EstimateTimeInSec"`
}

func (r *blockCountdownResult) toCountdown() *BlockCountdown {
	return &BlockCountdown{
		CurrentBlock:      r.CurrentBlock.unwrap(),
		CountdownBlock:    r.CountdownBlock.unwrap(),
		RemainingBlock:    r.RemainingBlock.unwrap(),
		EstimateTimeInSec: r.EstimateTimeInSec.unwrap(),
	}
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

	var result blockCountdownResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	return result.toCountdown(), nil
}

type BlockNumberRequest struct {
	Timestamp time.Time
	Closest   ClosestAvailableBlock
}

type ClosestAvailableBlock int32

const (
	closestAvailableBlockUnknown = iota
	ClosestAvailableBlockBefore
	ClosestAvailableBlockAfter
)

func (c ClosestAvailableBlock) tryString() (string, error) {
	switch c {
	case ClosestAvailableBlockBefore:
		return "before", nil

	case ClosestAvailableBlockAfter:
		return "after", nil

	default:
		return "", errors.Errorf("unknown closest available block parameter %d", int32(c))
	}
}

func (req *BlockNumberRequest) toParams() (map[string]string, error) {
	closest, err := req.Closest.tryString()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"timestamp": strconv.FormatInt(req.Timestamp.Unix(), 10),
		"closest":   closest,
	}, nil
}

func (c *BlocksClient) GetBlockNumber(
	ctx context.Context, req *BlockNumberRequest,
) (uint64, error) {
	params, err := req.toParams()
	if err != nil {
		return 0, err
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "getblocknobytime",
		other:  params,
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
	StartDate time.Time
	EndDate   time.Time
	Sort      SortingPreference
}

const dateFormat = "2006-01-02"

func (d *DateRange) toParams() (map[string]string, error) {
	return map[string]string{
		"startdate": d.StartDate.Format(dateFormat),
		"enddate":   d.EndDate.Format(dateFormat),
		"sort":      d.Sort.String(),
	}, nil
}

type AverageBlockSize struct {
	Timestamp      time.Time
	BlockSizeBytes uint32
}

type avgBlockSizeResult struct {
	UTCDate   string        `json:"UTCDate"`
	Timestamp unixTimestamp `json:"unixTimestamp"`
	BlockSize uint32        `json:"blockSize_bytes"`
}

func (r *avgBlockSizeResult) toBlockSize() *AverageBlockSize {
	return &AverageBlockSize{
		Timestamp:      r.Timestamp.unwrap(),
		BlockSizeBytes: r.BlockSize,
	}
}

func (c *BlocksClient) GetDailyAverageBlockSize(
	ctx context.Context, dates *DateRange,
) ([]AverageBlockSize, error) {
	params, err := dates.toParams()
	if err != nil {
		return nil, err
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "dailyavgblocksize",
		other:  params,
	})
	if err != nil {
		return nil, err
	}

	var result []avgBlockSizeResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	avgBlockSizes := make([]AverageBlockSize, len(result))
	for i := range avgBlockSizes {
		avgBlockSizes[i] = *result[i].toBlockSize()
	}

	return avgBlockSizes, nil
}

type BlockCount struct {
	DailyBlockRewards
	BlockCount uint32
}

type blockCountResult struct {
	blockRewardResult
	BlockCount uint32 `json:"blockCount"`
}

func (r *blockCountResult) toCount() *BlockCount {
	return &BlockCount{
		DailyBlockRewards: *r.blockRewardResult.toRewards(),
		BlockCount:        r.BlockCount,
	}
}

func (c *BlocksClient) GetDailyBlockCount(
	ctx context.Context, dates *DateRange,
) ([]BlockCount, error) {
	params, err := dates.toParams()
	if err != nil {
		return nil, err
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "dailyblkcount",
		other:  params,
	})
	if err != nil {
		return nil, err
	}

	var result []blockCountResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	counts := make([]BlockCount, len(result))
	for i := range result {
		counts[i] = *result[i].toCount()
	}

	return counts, nil
}

type DailyBlockRewards struct {
	Timestamp       time.Time
	BlockRewardsETH float64
}

type blockRewardResult struct {
	UTCDate      string        `json:"UTCDate"`
	Timestamp    unixTimestamp `json:"unixTimestamp"`
	BlockRewards floatStr      `json:"blockRewards_Eth"`
}

func (r *blockRewardResult) toRewards() *DailyBlockRewards {
	return &DailyBlockRewards{
		Timestamp:       r.Timestamp.unwrap(),
		BlockRewardsETH: r.BlockRewards.unwrap(),
	}
}

func (c *BlocksClient) GetDailyBlockRewards(
	ctx context.Context, dates *DateRange,
) ([]DailyBlockRewards, error) {
	params, err := dates.toParams()
	if err != nil {
		return nil, err
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "dailyblockrewards",
		other:  params,
	})
	if err != nil {
		return nil, err
	}

	var result []blockRewardResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	counts := make([]DailyBlockRewards, len(result))
	for i := range result {
		counts[i] = *result[i].toRewards()
	}

	return counts, nil
}

type DailyBlockTime struct {
	Timestamp        time.Time
	BlockTimeSeconds float64
}

type dailyBlockTimeResult struct {
	UTCDate       string        `json:"UTCDate"`
	Timestamp     unixTimestamp `json:"unixTimestamp"`
	BlockTimeSecs floatStr      `json:"blockTime_sec"`
}

func (r *dailyBlockTimeResult) toBlockTime() *DailyBlockTime {
	return &DailyBlockTime{
		Timestamp:        r.Timestamp.unwrap(),
		BlockTimeSeconds: r.BlockTimeSecs.unwrap(),
	}
}

func (c *BlocksClient) GetDailyAverageBlockTime(
	ctx context.Context, dates *DateRange,
) ([]DailyBlockTime, error) {
	params, err := dates.toParams()
	if err != nil {
		return nil, err
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "dailyavgblocktime",
		other:  params,
	})
	if err != nil {
		return nil, err
	}

	var result []dailyBlockTimeResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	counts := make([]DailyBlockTime, len(result))
	for i := range result {
		counts[i] = *result[i].toBlockTime()
	}

	return counts, nil
}

type DailyUnclesCount struct {
	Timestamp            time.Time
	UncleBlockCount      uint32
	UncleBlockRewardsETH float64
}

type dailyUnclesResult struct {
	UTCDate           string        `json:"UTCDate"`
	Timestamp         unixTimestamp `json:"unixTimestamp"`
	UncleBlockCount   uint32        `json:"uncleBlockCount"`
	UncleBlockRewards floatStr      `json:"uncleBlockRewards_Eth"`
}

func (r *dailyUnclesResult) toUnclesCount() *DailyUnclesCount {
	return &DailyUnclesCount{
		Timestamp:            r.Timestamp.unwrap(),
		UncleBlockCount:      r.UncleBlockCount,
		UncleBlockRewardsETH: r.UncleBlockRewards.unwrap(),
	}
}

func (c *BlocksClient) GetDailyUnclesCount(
	ctx context.Context, dates *DateRange,
) ([]DailyUnclesCount, error) {
	params, err := dates.toParams()
	if err != nil {
		return nil, err
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: blocksModule,
		action: "dailyuncleblkcount",
		other:  params,
	})
	if err != nil {
		return nil, err
	}

	var result []dailyUnclesResult
	if err := json.Unmarshal(rspData, &result); err != nil {
		return nil, err
	}

	counts := make([]DailyUnclesCount, len(result))
	for i := range result {
		counts[i] = *result[i].toUnclesCount()
	}

	return counts, nil
}
