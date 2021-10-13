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

func TestAccount(t *testing.T) {
	m := newMockServer("account", true)
	t.Cleanup(m.close)

	u, err := url.Parse(m.srv.URL)
	require.NoError(t, err)

	client := etherscan.New(&etherscan.Params{
		APIKey:  m.apiKey,
		BaseURL: u,
	})

	ctx := context.Background()

	t.Run("GetETHBalance", func(t *testing.T) {
		bal, err := client.Accounts.GetETHBalance(ctx, &etherscan.ETHBalanceRequest{
			Address: common.HexToAddress("0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae"),
			Tag:     etherscan.BlockParameterLatest,
		})
		require.NoError(t, err)
		assert.Equal(t, "40891626854930000000000", bal.String())
	})

	var multiETHBalAddrs = []common.Address{
		common.HexToAddress("0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a"),
		common.HexToAddress("0x63a9975ba31b0b9626b34300f7f627147df1f526"),
		common.HexToAddress("0x198ef1ec325a96cc354c7266a038be8b5c558f67"),
	}

	t.Run("MultiGetETHBalance", func(t *testing.T) {
		bals, err := client.Accounts.GetMultiETHBalances(ctx, &etherscan.MultiETHBalancesRequest{
			Addresses: multiETHBalAddrs,
			Tag:       etherscan.BlockParameterLatest,
		})
		require.NoError(t, err)
		require.Len(t, bals, 3)

		cupaloy.SnapshotT(t, bals)
	})

	t.Run("ListNormalTxs", func(t *testing.T) {
		txs, err := client.Accounts.ListNormalTransactions(ctx, &etherscan.ListTxRequest{
			Address:    common.HexToAddress("0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a"),
			StartBlock: 0,
			EndBlock:   99999999,
			Sort:       etherscan.SortingPreferenceAscending,
		})
		require.NoError(t, err)
		require.Len(t, txs, 2)

		cupaloy.SnapshotT(t, txs)
	})

	t.Run("ListInternalTxs", func(t *testing.T) {
		txs, err := client.Accounts.ListInternalTransactions(ctx, &etherscan.ListTxRequest{
			Address:    common.HexToAddress("0x2c1ba59d6f58433fb1eaee7d20b26ed83bda51a3"),
			StartBlock: 0,
			EndBlock:   99999999,
			Sort:       etherscan.SortingPreferenceAscending,
		})
		require.NoError(t, err)
		require.Len(t, txs, 2)

		cupaloy.SnapshotT(t, txs)
	})

	t.Run("GetInternalTxsByHash", func(t *testing.T) {
		txs, err := client.Accounts.GetInternalTxsByHash(
			ctx,
			common.HexToHash("0x40eb908387324f2b575b4879cd9d7188f69c8fc9d87c901b9e2daaea4b442170"),
		)
		require.NoError(t, err)
		require.Len(t, txs, 1)

		cupaloy.SnapshotT(t, txs)
	})

	t.Run("GetInternalTxsBlockRange", func(t *testing.T) {
		txs, err := client.Accounts.GetInternalTxsByBlockRange(ctx, &etherscan.BlockRangeRequest{
			StartBlock: 0,
			EndBlock:   2702578,
			Sort:       etherscan.SortingPreferenceAscending,
		})
		require.NoError(t, err)
		require.Len(t, txs, 2)

		cupaloy.SnapshotT(t, txs)
	})

	t.Run("ListTokenTransfers", func(t *testing.T) {
		txs, err := client.Accounts.ListTokenTransfers(ctx, &etherscan.TokenTransfersRequest{
			Address:         common.HexToAddress("0x4e83362442b8d1bec281594cea3050c8eb01311c"),
			ContractAddress: common.HexToAddress("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2"),
			Sort:            etherscan.SortingPreferenceAscending,
		})
		require.NoError(t, err)
		require.Len(t, txs, 2)

		cupaloy.SnapshotT(t, txs)
	})

	t.Run("ListNFTTransfers", func(t *testing.T) {
		address := common.HexToAddress("0x6975be450864c02b4613023c2152ee0743572325")
		contractAddress := common.HexToAddress("0x06012c8cf97bead5deae237070f9587f8e7a266d")

		txs, err := client.Accounts.ListNFTTransfers(ctx, &etherscan.ListNFTTransferRequest{
			Address:         &address,
			ContractAddress: &contractAddress,
			Sort:            etherscan.SortingPreferenceAscending,
		})
		require.NoError(t, err)
		require.Len(t, txs, 2)

		cupaloy.SnapshotT(t, txs)
	})

	t.Run("ListMinedBlocks", func(t *testing.T) {
		blocks, err := client.Accounts.ListBlocksMined(ctx, &etherscan.ListBlocksRequest{
			Address: common.HexToAddress("0x9dd134d14d1e65f84b706d6f205cd5b1cd03a46b"),
			Type:    etherscan.BlockTypeBlocks,
		})
		require.NoError(t, err)
		require.Len(t, blocks, 3)

		cupaloy.SnapshotT(t, blocks)
	})
}
