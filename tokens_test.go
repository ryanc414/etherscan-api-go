package etherscan_test

import (
	"context"
	"math/big"
	"net/url"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ryanc414/etherscan-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	m := newMockServer("token", false)
	t.Cleanup(m.close)

	u, err := url.Parse(m.srv.URL)
	require.NoError(t, err)

	client := etherscan.New(&etherscan.Params{
		APIKey:  m.apiKey,
		BaseURL: u,
	})

	ctx := context.Background()

	t.Run("GetTotalSupply", func(t *testing.T) {
		result, err := client.Tokens.GetTotalSupply(
			ctx, common.HexToAddress("0x0e3a2a1f2146d86a604adc220b4967a898d7fe07"),
		)
		require.NoError(t, err)
		assert.Equal(t, 0, result.Cmp(big.NewInt(21265524714464)))
	})

	t.Run("GetAccountBalance", func(t *testing.T) {
		result, err := client.Tokens.GetAccountBalance(ctx, &etherscan.BalanceRequest{
			ContractAddress: common.HexToAddress("0x57d90b64a1a57749b0f932f1a3395792e12e7055"),
			Address:         common.HexToAddress("0xe04f27eb70e025b78871a2ad7eabe85e61212761"),
		})
		require.NoError(t, err)
		assert.Equal(t, 0, result.Cmp(big.NewInt(135499)))
	})

	t.Run("GetHistoricalSupply", func(t *testing.T) {
		result, err := client.Tokens.GetHistoricalSupply(ctx, &etherscan.HistoricalSupplyRequest{
			ContractAddress: common.HexToAddress("0x57d90b64a1a57749b0f932f1a3395792e12e7055"),
			BlockNo:         8000000,
		})
		require.NoError(t, err)
		assert.Equal(t, 0, result.Cmp(big.NewInt(21265524714464)))
	})

	t.Run("GetHistoricalBalance", func(t *testing.T) {
		result, err := client.Tokens.GetHistoricalBalance(ctx, &etherscan.HistoricalBalanceRequest{
			ContractAddress: common.HexToAddress("0x57d90b64a1a57749b0f932f1a3395792e12e7055"),
			Address:         common.HexToAddress("0xe04f27eb70e025b78871a2ad7eabe85e61212761"),
			BlockNo:         8000000,
		})
		require.NoError(t, err)
		assert.Equal(t, 0, result.Cmp(big.NewInt(135499)))
	})

	t.Run("GetTokenInfo", func(t *testing.T) {
		result, err := client.Tokens.GetTokenInfo(ctx, common.HexToAddress("0x0e3a2a1f2146d86a604adc220b4967a898d7fe07"))
		require.NoError(t, err)
		require.Len(t, result, 1)

		cupaloy.SnapshotT(t, result)
	})
}
