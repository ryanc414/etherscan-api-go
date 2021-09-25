package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

type LogsClient struct {
	api *apiClient
}

const logsModule = "logs"

type LogsRequest struct {
	FromBlock   LogsBlockParam
	ToBlock     LogsBlockParam
	Address     common.Address
	Topics      []common.Hash
	Comparisons []TopicComparison
}

func (req *LogsRequest) toParams() (map[string]string, error) {
	params := make(map[string]string)

	fromBlock, err := req.FromBlock.toParam()
	if err != nil {
		return nil, err
	}
	params["fromBlock"] = fromBlock

	toBlock, err := req.ToBlock.toParam()
	if err != nil {
		return nil, err
	}
	params["toBlock"] = toBlock

	params["address"] = req.Address.String()

	if err := req.addTopicParams(params); err != nil {
		return nil, err
	}

	if err := req.addCompParams(params); err != nil {
		return nil, err
	}

	return params, nil
}

func (req *LogsRequest) addTopicParams(params map[string]string) error {
	if len(req.Topics) > 4 {
		return errors.New("a maximum of 4 topics is allowed")
	}

	for i := range req.Topics {
		key := fmt.Sprintf("topic%d", i)
		params[key] = req.Topics[i].String()
	}

	return nil
}

func (req *LogsRequest) addCompParams(params map[string]string) error {
	for i := range req.Comparisons {
		k, v, err := req.Comparisons[i].toParam()
		if err != nil {
			return err
		}

		params[k] = v
	}

	return nil
}

type LogsBlockParam struct {
	Number uint64
	Latest bool
}

func (b LogsBlockParam) toParam() (string, error) {
	if b.Latest {
		if b.Number != 0 {
			return "", errors.New("number must not be specified when latest is true for block")
		}

		return "latest", nil
	}

	return strconv.FormatUint(b.Number, 10), nil
}

type TopicComparison struct {
	Topics   [2]uint8
	Operator ComparisonOperator
}

func (c *TopicComparison) toParam() (string, string, error) {
	if c.Topics[1] <= c.Topics[0] {
		return "", "", errors.New("second topic must be greater than first")
	}

	key := fmt.Sprintf("topic%d_%d_opr", c.Topics[0], c.Topics[1])
	val := c.Operator.String()

	return key, val, nil
}

type ComparisonOperator int32

const (
	ComparisonOperatorAnd = iota
	ComparisonOperatorOr
)

func (op ComparisonOperator) String() string {
	switch op {
	case ComparisonOperatorAnd:
		return "and"

	case ComparisonOperatorOr:
		return "or"

	default:
		panic(fmt.Sprintf("unexpected comparison operator %d", int32(op)))
	}
}

type LogResponse struct {
	Address          common.Address
	BlockNumber      uint64
	Data             []byte
	GasPrice         *big.Int
	GasUsed          *big.Int
	LogIndex         uint32
	Timestamp        time.Time
	Topics           []common.Hash
	TransactionHash  common.Hash
	TransactionIndex uint32
}

type logResult struct {
	Address          string         `json:"address"`
	BlockNumber      hexutil.Uint64 `json:"blockNumber"`
	Data             string         `json:"data"`
	GasPrice         *hexutil.Big   `json:"gasPrice"`
	GasUsed          *hexutil.Big   `json:"gasUsed"`
	LogIndex         hexUint        `json:"logIndex"`
	Timestamp        hexTimestamp   `json:"timeStamp"`
	Topics           []string       `json:"topics"`
	TransactionHash  string         `json:"transactionHash"`
	TransactionIndex hexUint        `json:"transactionIndex"`
}

func (res *logResult) toLog() LogResponse {
	topics := make([]common.Hash, len(res.Topics))
	for i := range res.Topics {
		topics[i] = common.HexToHash(res.Topics[i])
	}

	return LogResponse{
		Address:          common.HexToAddress(res.Address),
		BlockNumber:      uint64(res.BlockNumber),
		Data:             common.Hex2Bytes(res.Data),
		GasPrice:         res.GasPrice.ToInt(),
		GasUsed:          res.GasUsed.ToInt(),
		LogIndex:         uint32(res.LogIndex),
		Timestamp:        res.Timestamp.unwrap(),
		Topics:           topics,
		TransactionHash:  common.HexToHash(res.TransactionHash),
		TransactionIndex: uint32(res.TransactionIndex),
	}
}

func (c *LogsClient) GetLogs(ctx context.Context, req *LogsRequest) ([]LogResponse, error) {
	params, err := req.toParams()
	if err != nil {
		return nil, err
	}

	rspData, err := c.api.get(ctx, &requestParams{
		module: logsModule,
		action: "getLogs",
		other:  params,
	})
	if err != nil {
		return nil, err
	}

	var logResult []logResult
	if err := json.Unmarshal(rspData, &logResult); err != nil {
		return nil, err
	}

	logs := make([]LogResponse, len(logResult))
	for i := range logResult {
		logs[i] = logResult[i].toLog()
	}

	return logs, nil
}
