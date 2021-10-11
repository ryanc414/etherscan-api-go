package etherscan_test

import (
	"context"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/ryanc414/etherscan-api-go"
	"github.com/stretchr/testify/require"
)

func TestContracts(t *testing.T) {
	m := newMockServer("contract")
	t.Cleanup(m.close)

	u, err := url.Parse(m.srv.URL)
	require.NoError(t, err)

	client := etherscan.New(&etherscan.Params{
		APIKey:  m.apiKey,
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

type mockContractsAPI struct {
	apiKey                string
	getSourceCodeResponse []byte
}

const getSourceCodeFile = "getSourceCodeResponse.json"

func newMockContractsAPI() (*mockContractsAPI, error) {
	data, err := ioutil.ReadFile(getSourceCodeFile)
	if err != nil {
		return nil, err
	}

	return &mockContractsAPI{
		apiKey:                uuid.NewString(),
		getSourceCodeResponse: data,
	}, nil
}
