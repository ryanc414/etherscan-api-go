package contracts_test

import (
	"context"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ryanc414/etherscan-api-go"
	"github.com/ryanc414/etherscan-api-go/testbed"
	"github.com/stretchr/testify/require"
)

func TestContracts(t *testing.T) {
	m := testbed.NewMockServer("contract", true)
	t.Cleanup(m.Close)

	u, err := m.URL()
	require.NoError(t, err)

	client := etherscan.New(&etherscan.Params{
		APIKey:  m.APIKey,
		BaseURL: u,
	})

	ctx := context.Background()

	t.Run("GetContractABI", func(t *testing.T) {
		abi, err := client.Contracts.GetContractABI(
			ctx, common.HexToAddress("0xBB9bc244D798123fDe783fCc1C72d3Bb8C189413"),
		)
		require.NoError(t, err)
		cupaloy.SnapshotT(t, abi)
	})

	t.Run("GetContractSourceCode", func(t *testing.T) {
		info, err := client.Contracts.GetContractSourceCode(
			ctx,
			common.HexToAddress("0xBB9bc244D798123fDe783fCc1C72d3Bb8C189413"),
		)
		require.NoError(t, err)
		require.Len(t, info, 1)

		cupaloy.SnapshotT(t, info)
	})
}
