package etherscan_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ryanc414/etherscan-api-go"
	"github.com/stretchr/testify/require"
)

func TestLogs(t *testing.T) {
	m := newMockServer("logs")
	t.Cleanup(m.close)

	u, err := url.Parse(m.srv.URL)
	require.NoError(t, err)

	client := etherscan.New(&etherscan.Params{
		APIKey:  m.apiKey,
		BaseURL: u,
	})

	ctx := context.Background()
	topic0 := common.HexToHash("0xf63780e752c6a54a94fc52715dbc5518a3b4c3c2833d301a204226548a2a8545")

	t.Run("GetLogsLatest", func(t *testing.T) {
		logs, err := client.Logs.GetLogs(ctx, &etherscan.LogsRequest{
			FromBlock: etherscan.LogsBlockParam{Number: 379224},
			ToBlock:   etherscan.LogsBlockParam{Latest: true},
			Address:   common.HexToAddress("0x33990122638b9132ca29c723bdf037f1a891a70c"),
			Topics:    []common.Hash{topic0},
		})
		require.NoError(t, err)
		require.Len(t, logs, 2)

		cupaloy.SnapshotT(t, logs)
	})

	t.Run("GetLogsFixed", func(t *testing.T) {
		topic1 := common.HexToHash("0x72657075746174696f6e00000000000000000000000000000000000000000000")
		logs, err := client.Logs.GetLogs(ctx, &etherscan.LogsRequest{
			FromBlock: etherscan.LogsBlockParam{Number: 379224},
			ToBlock:   etherscan.LogsBlockParam{Number: 400000},
			Address:   common.HexToAddress("0x33990122638b9132ca29c723bdf037f1a891a70c"),
			Topics:    []common.Hash{topic0, topic1},
			Comparisons: []etherscan.TopicComparison{
				{
					Topics:   [2]uint8{0, 1},
					Operator: etherscan.ComparisonOperatorAnd,
				},
			},
		})
		require.NoError(t, err)
		require.Len(t, logs, 1)

		cupaloy.SnapshotT(t, logs)
	})
}
