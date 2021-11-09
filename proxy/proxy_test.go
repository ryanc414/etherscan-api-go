package proxy_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ryanc414/etherscan-api-go"
	ecommon "github.com/ryanc414/etherscan-api-go/common"
	"github.com/ryanc414/etherscan-api-go/proxy"
	"github.com/ryanc414/etherscan-api-go/testbed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxy(t *testing.T) {
	m := testbed.NewMockServer("proxy", true)
	t.Cleanup(m.Close)

	u, err := m.URL()
	require.NoError(t, err)

	client := etherscan.New(&etherscan.Params{
		APIKey:  m.APIKey,
		BaseURL: u,
	})

	ctx := context.Background()

	t.Run("BlockNumber", func(t *testing.T) {
		num, err := client.Proxy.BlockNumber(ctx)
		require.NoError(t, err)
		assert.Equal(t, uint64(12806953), num)
	})

	t.Run("GetBlockByNumberFull", func(t *testing.T) {
		block, err := client.Proxy.GetBlockByNumberFull(ctx, 68943)
		require.NoError(t, err)
		cupaloy.SnapshotT(t, block)
	})

	t.Run("GetBlockByNumberSummary", func(t *testing.T) {
		block, err := client.Proxy.GetBlockByNumberSummary(ctx, 68943)
		require.NoError(t, err)
		cupaloy.SnapshotT(t, block)
	})

	t.Run("GetUncleByBlockNumberAndIndex", func(t *testing.T) {
		uncle, err := client.Proxy.GetUncleByBlockNumberAndIndex(ctx, &proxy.BlockNumberAndIndex{
			Number: 12989046,
			Index:  0,
		})
		require.NoError(t, err)

		cupaloy.SnapshotT(t, uncle)
	})

	t.Run("GetBlockTransactionCount", func(t *testing.T) {
		txCount, err := client.Proxy.GetBlockTransactionCountByNumber(ctx, 1112952)
		require.NoError(t, err)
		assert.Equal(t, uint32(3), txCount)
	})

	t.Run("GetTransactionByHash", func(t *testing.T) {
		txHash := common.HexToHash("0x1e2910a262b1008d0616a0beb24c1a491d78771baa54a33e66065e03b1f46bc1")
		txInfo, err := client.Proxy.GetTransactionByHash(ctx, txHash)
		require.NoError(t, err)

		cupaloy.SnapshotT(t, txInfo)
	})

	t.Run("GetTxByBlockNumberAndIndex", func(t *testing.T) {
		txInfo, err := client.Proxy.GetTransactionByBlockNumberAndIndex(ctx, &proxy.BlockNumberAndIndex{
			Number: 12989213,
			Index:  282,
		})
		require.NoError(t, err)

		cupaloy.SnapshotT(t, txInfo)
	})

	t.Run("GetTransactionCount", func(t *testing.T) {
		count, err := client.Proxy.GetTransactionCount(ctx, &proxy.TxCountRequest{
			Address: common.HexToAddress("0x4bd5900Cb274ef15b153066D736bf3e83A9ba44e"),
			Tag:     ecommon.BlockParameterLatest,
		})
		require.NoError(t, err)

		assert.Equal(t, uint64(68), count)
	})

	t.Run("SendRawTransaction", func(t *testing.T) {
		result, err := client.Proxy.SendRawTransaction(ctx, hexutil.MustDecode("0xf904808000831cfde080"))
		require.NoError(t, err)

		expectedResult := common.HexToHash("0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331")
		assert.Equal(t, expectedResult, result)
	})

	t.Run("GetTransactionReceipt", func(t *testing.T) {
		txHash := common.HexToHash("0x1e2910a262b1008d0616a0beb24c1a491d78771baa54a33e66065e03b1f46bc1")
		receipt, err := client.Proxy.GetTransactionReceipt(ctx, txHash)
		require.NoError(t, err)

		cupaloy.SnapshotT(t, receipt)
	})

	t.Run("Call", func(t *testing.T) {
		result, err := client.Proxy.Call(ctx, &proxy.CallRequest{
			To:   common.HexToAddress("0xAEEF46DB4855E25702F8237E8f403FddcaF931C0"),
			Data: hexutil.MustDecode("0x70a08231000000000000000000000000e16359506c028e51f16be38986ec5746251e9724"),
			Tag:  ecommon.BlockParameterLatest,
		})
		require.NoError(t, err)

		expectedResult := hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000601d8888141c00")
		assert.Equal(t, expectedResult, result)
	})

	t.Run("GetCode", func(t *testing.T) {
		result, err := client.Proxy.GetCode(ctx, &proxy.GetCodeRequest{
			Address: common.HexToAddress("0xf75e354c5edc8efed9b59ee9f67a80845ade7d0c"),
			Tag:     ecommon.BlockParameterLatest,
		})
		require.NoError(t, err)

		expectedResult := hexutil.MustDecode("0x3660008037602060003660003473273930d21e01ee25e4c219b63259d214872220a261235a5a03f21560015760206000f3")
		assert.Equal(t, expectedResult, result)
	})

	t.Run("GetStorageAt", func(t *testing.T) {
		result, err := client.Proxy.GetStorageAt(ctx, &proxy.GetStorageRequest{
			Address:  common.HexToAddress("0x6e03d9cce9d60f3e9f2597e13cd4c54c55330cfd"),
			Position: 0,
			Tag:      ecommon.BlockParameterLatest,
		})
		require.NoError(t, err)

		expectedResult := hexutil.MustDecode("0x0000000000000000000000003d0768da09ce77d25e2d998e6a7b6ed4b9116c2d")
		assert.Equal(t, expectedResult, result)
	})

	t.Run("GetGasPrice", func(t *testing.T) {
		price, err := client.Proxy.GasPrice(ctx)
		require.NoError(t, err)

		expectedPrice := big.NewInt(18000000000)
		assert.Equal(t, 0, expectedPrice.Cmp(price))
	})

	t.Run("EstimateGas", func(t *testing.T) {
		gas, err := client.Proxy.EstimateGas(ctx, &proxy.EstimateGasRequest{
			Data:     hexutil.MustDecode("0x4e71d92d"),
			To:       common.HexToAddress("0xf0160428a8552ac9bb7e050d90eeade4ddd52843"),
			Value:    big.NewInt(65314),
			GasPrice: big.NewInt(21971876044),
			Gas:      big.NewInt(99999999),
		})
		require.NoError(t, err)

		expectedGas := big.NewInt(25942)
		require.Equal(t, 0, gas.Cmp(expectedGas))
	})
}
