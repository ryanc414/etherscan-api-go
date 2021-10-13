package etherscan_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ryanc414/etherscan-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactions(t *testing.T) {
	m := newMockServer("transaction", true)
	t.Cleanup(m.close)

	u, err := url.Parse(m.srv.URL)
	require.NoError(t, err)

	client := etherscan.New(&etherscan.Params{
		APIKey:  m.apiKey,
		BaseURL: u,
	})

	ctx := context.Background()

	t.Run("GetExecutionStatus", func(t *testing.T) {
		txhash := common.HexToHash("0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a")
		status, err := client.Transactions.GetExecutionStatus(ctx, txhash)
		require.NoError(t, err)

		cupaloy.SnapshotT(t, status)
	})

	t.Run("GetTxReceiptStatus", func(t *testing.T) {
		txHash := common.HexToHash("0x513c1ba0bebf66436b5fed86ab668452b7805593c05073eb2d51d3a52f480a76")
		status, err := client.Transactions.GetTxReceiptStatus(ctx, txHash)
		require.NoError(t, err)
		assert.True(t, status)
	})
}
