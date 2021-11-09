package logs_test

import (
	"context"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ryanc414/etherscan-api-go"
	"github.com/ryanc414/etherscan-api-go/logs"
	"github.com/ryanc414/etherscan-api-go/testbed"
	"github.com/stretchr/testify/require"
)

func TestLogs(t *testing.T) {
	m := testbed.NewMockServer("logs", true)
	t.Cleanup(m.Close)

	u, err := m.URL()
	require.NoError(t, err)

	client := etherscan.New(&etherscan.Params{
		APIKey:  m.APIKey,
		BaseURL: u,
	})

	ctx := context.Background()
	topic0 := common.HexToHash("0xf63780e752c6a54a94fc52715dbc5518a3b4c3c2833d301a204226548a2a8545")

	t.Run("GetLogsLatest", func(t *testing.T) {
		logs, err := client.Logs.GetLogs(ctx, &logs.LogsRequest{
			FromBlock: logs.LogsBlockParam{Number: 379224},
			ToBlock:   logs.LogsBlockParam{Latest: true},
			Address:   common.HexToAddress("0x33990122638b9132ca29c723bdf037f1a891a70c"),
			Topics:    []common.Hash{topic0},
		})
		require.NoError(t, err)
		require.Len(t, logs, 2)

		cupaloy.SnapshotT(t, logs)
	})

	t.Run("GetLogsFixed", func(t *testing.T) {
		topic1 := common.HexToHash("0x72657075746174696f6e00000000000000000000000000000000000000000000")
		logs, err := client.Logs.GetLogs(ctx, &logs.LogsRequest{
			FromBlock: logs.LogsBlockParam{Number: 379224},
			ToBlock:   logs.LogsBlockParam{Number: 400000},
			Address:   common.HexToAddress("0x33990122638b9132ca29c723bdf037f1a891a70c"),
			Topics:    []common.Hash{topic0, topic1},
			Comparisons: []logs.TopicComparison{
				{
					Topics:   [2]uint8{0, 1},
					Operator: logs.ComparisonOperatorAnd,
				},
			},
		})
		require.NoError(t, err)
		require.Len(t, logs, 1)

		cupaloy.SnapshotT(t, logs)
	})
}
