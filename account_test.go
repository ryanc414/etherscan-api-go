package etherscan_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ryanc414/etherscan-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	m := newMockServer("account")
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

		expectedBals := []string{
			"40891626854930000000000",
			"332567136222827062478",
			"0",
		}
		for i := range bals {
			assert.Equal(t, multiETHBalAddrs[i], bals[i].Account)
			assert.Equal(t, expectedBals[i], bals[i].Balance.String())
		}
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

		assert.Equal(t, uint64(0), txs[0].BlockNumber)
		assert.Equal(t, uint64(47884), txs[1].BlockNumber)
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

		assert.Equal(
			t,
			common.HexToHash("0x8a1a9989bda84f80143181a68bc137ecefa64d0d4ebde45dd94fc0cf49e70cb6"),
			txs[0].Hash,
		)

		assert.Equal(
			t,
			common.HexToHash("0x1a50f1dc0bc912745f7d09b988669f71d199719e2fb7592c2074ede9578032d0"),
			txs[1].Hash,
		)
	})

	t.Run("GetInternalTxsByHash", func(t *testing.T) {
		txs, err := client.Accounts.GetInternalTxsByHash(
			ctx,
			common.HexToHash("0x40eb908387324f2b575b4879cd9d7188f69c8fc9d87c901b9e2daaea4b442170"),
		)
		require.NoError(t, err)
		require.Len(t, txs, 1)

		assert.Equal(t, uint64(1743059), txs[0].BlockNumber)
	})

	t.Run("GetInternalTxsBlockRange", func(t *testing.T) {
		txs, err := client.Accounts.GetInternalTxsByBlockRange(ctx, &etherscan.BlockRangeRequest{
			StartBlock: 0,
			EndBlock:   2702578,
			Sort:       etherscan.SortingPreferenceAscending,
		})
		require.NoError(t, err)
		require.Len(t, txs, 2)

		assert.Equal(
			t,
			common.HexToHash("0x3f97c969ddf71f515ce5373b1f8e76e9fd7016611d8ce455881009414301789e"),
			txs[0].Hash,
		)

		assert.Equal(
			t,
			common.HexToHash("0x893c428fed019404f704cf4d9be977ed9ca01050ed93dccdd6c169422155586f"),
			txs[1].Hash,
		)
	})

	t.Run("ListTokenTransfers", func(t *testing.T) {
		txs, err := client.Accounts.ListTokenTransfers(ctx, &etherscan.TokenTransfersRequest{
			Address:         common.HexToAddress("0x4e83362442b8d1bec281594cea3050c8eb01311c"),
			ContractAddress: common.HexToAddress("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2"),
			Sort:            etherscan.SortingPreferenceAscending,
		})
		require.NoError(t, err)
		require.Len(t, txs, 2)

		assert.Equal(
			t,
			common.HexToHash("0xe8c208398bd5ae8e4c237658580db56a2a94dfa0ca382c99b776fa6e7d31d5b4"),
			txs[0].Hash,
		)

		assert.Equal(
			t,
			common.HexToHash("0x9c82e89b7f6a4405d11c361adb6d808d27bcd9db3b04b3fb3bc05d182bbc5d6f"),
			txs[1].Hash,
		)
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
		assert.Equal(t, "202106", txs[0].TokenID)
		assert.Equal(t, "147739", txs[1].TokenID)
	})

	t.Run("ListMinedBlocks", func(t *testing.T) {
		blocks, err := client.Accounts.ListBlocksMined(ctx, &etherscan.ListBlocksRequest{
			Address: common.HexToAddress("0x9dd134d14d1e65f84b706d6f205cd5b1cd03a46b"),
			Type:    etherscan.BlockTypeBlocks,
		})
		require.NoError(t, err)
		require.Len(t, blocks, 3)

		assert.Equal(t, uint64(3462296), blocks[0].BlockNumber)
		assert.Equal(t, uint64(2691400), blocks[1].BlockNumber)
		assert.Equal(t, uint64(2687700), blocks[2].BlockNumber)
	})
}
