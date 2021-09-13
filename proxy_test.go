package etherscan_test

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/ryanc414/etherscan-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxy(t *testing.T) {
	m := newProxyAPI()

	ts := httptest.NewServer(m)
	t.Cleanup(ts.Close)

	u, err := url.Parse(ts.URL)
	require.NoError(t, err)

	client := etherscan.New(&etherscan.Params{
		APIKey:  m.apiKey,
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

		expectedBlockHash := common.HexToHash("0x7eb7c23a5ac2f2d70aa1ba4e5c56d89de5ac993590e5f6e79c394e290d998ba8")
		assert.Equal(t, expectedBlockHash, block.Hash)

		assert.Len(t, block.Transactions, 1)

		expectedTxHash := common.HexToHash("0xa442249820de6be754da81eafbd44a865773e4b23d7c0522d31fd03977823008")
		assert.Equal(t, expectedTxHash, block.Transactions[0].Hash)
	})

	t.Run("GetBlockByNumberSummary", func(t *testing.T) {
		block, err := client.Proxy.GetBlockByNumberSummary(ctx, 68943)
		require.NoError(t, err)

		expectedBlockHash := common.HexToHash("0x7eb7c23a5ac2f2d70aa1ba4e5c56d89de5ac993590e5f6e79c394e290d998ba8")
		assert.Equal(t, expectedBlockHash, block.Hash)

		assert.Len(t, block.Transactions, 1)

		expectedTxHash := common.HexToHash("0xa442249820de6be754da81eafbd44a865773e4b23d7c0522d31fd03977823008")
		assert.Equal(t, expectedTxHash, block.Transactions[0])
	})

	t.Run("GetUncleByBlockNumberAndIndex", func(t *testing.T) {
		uncle, err := client.Proxy.GetUncleByBlockNumberAndIndex(ctx, &etherscan.BlockNumberAndIndex{
			Number: 12989046,
			Index:  0,
		})
		require.NoError(t, err)

		expectedHash := common.HexToHash("0x1da88e3581315d009f1cb600bf06f509cd27a68cb3d6437bda8698d04089f14a")
		assert.Equal(t, expectedHash, uncle.Hash)
	})

	t.Run("GetBlockTransactionCount", func(t *testing.T) {
		txCount, err := client.Proxy.GetBlockTransactionCountByNumber(ctx, 1112952)
		require.NoError(t, err)
		assert.Equal(t, 3, txCount)
	})

	t.Run("GetTransactionByHash", func(t *testing.T) {
		txHash := common.HexToHash("0x1e2910a262b1008d0616a0beb24c1a491d78771baa54a33e66065e03b1f46bc1")
		txInfo, err := client.Proxy.GetTransactionByHash(ctx, txHash)
		require.NoError(t, err)

		assert.Equal(t, 0, txInfo.Value.Cmp(big.NewInt(10000000000)))
	})

	t.Run("GetTxByBlockNumberAndIndex", func(t *testing.T) {
		txInfo, err := client.Proxy.GetTransactionByBlockNumberAndIndex(ctx, &etherscan.BlockNumberAndIndex{
			Number: 12989213,
			Index:  282,
		})
		require.NoError(t, err)

		assert.Equal(t, 0, txInfo.Value.Cmp(big.NewInt(10000000000)))
	})

	t.Run("GetTransactionCount", func(t *testing.T) {
		count, err := client.Proxy.GetTransactionCount(ctx, &etherscan.TxCountRequest{
			Address: common.HexToAddress("0x4bd5900Cb274ef15b153066D736bf3e83A9ba44e"),
			Tag:     etherscan.BlockParameterLatest,
		})
		require.NoError(t, err)

		assert.Equal(t, uint64(68), count)
	})

	t.Run("SendRawTransaction", func(t *testing.T) {
		result, err := client.Proxy.SendRawTransaction(ctx, common.Hex2Bytes("0xf904808000831cfde080"))
		require.NoError(t, err)

		expectedResult := common.Hex2Bytes("0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331")
		assert.Equal(t, expectedResult, result)
	})

	t.Run("GetTransactionReceipt", func(t *testing.T) {
		txHash := common.HexToHash("0x1e2910a262b1008d0616a0beb24c1a491d78771baa54a33e66065e03b1f46bc1")
		receipt, err := client.Proxy.GetTransactionReceipt(ctx, txHash)
		require.NoError(t, err)

		address := common.HexToAddress("0xc778417e063141139fce010982780140aa0cd5ab")
		assert.Equal(t, address, receipt.To)

		require.Len(t, receipt.Logs, 1)
		require.Len(t, receipt.Logs[0].Topics, 2)
	})

	t.Run("Call", func(t *testing.T) {
		result, err := client.Proxy.Call(ctx, &etherscan.CallRequest{
			To:   common.HexToAddress("0xAEEF46DB4855E25702F8237E8f403FddcaF931C0"),
			Data: common.HexToHash("0x70a08231000000000000000000000000e16359506c028e51f16be38986ec5746251e9724"),
			Tag:  etherscan.BlockParameterLatest,
		})
		require.NoError(t, err)

		expectedResult := common.Hex2Bytes("0x00000000000000000000000000000000000000000000000000601d8888141c00")
		assert.Equal(t, expectedResult, result)
	})

	t.Run("GetCode", func(t *testing.T) {
		result, err := client.Proxy.GetCode(ctx, &etherscan.GetCodeRequest{
			Address: common.HexToAddress("0xf75e354c5edc8efed9b59ee9f67a80845ade7d0c"),
			Tag:     etherscan.BlockParameterLatest,
		})
		require.NoError(t, err)

		expectedResult := common.Hex2Bytes("0x3660008037602060003660003473273930d21e01ee25e4c219b63259d214872220a261235a5a03f21560015760206000f3")
		assert.Equal(t, expectedResult, result)
	})

	t.Run("GetStorageAt", func(t *testing.T) {
		result, err := client.Proxy.GetStorageAt(ctx, &etherscan.GetStorageRequest{
			Address:  common.HexToAddress("0x6e03d9cce9d60f3e9f2597e13cd4c54c55330cfd"),
			Position: 0,
			Tag:      etherscan.BlockParameterLatest,
		})
		require.NoError(t, err)

		expectedResult := common.Hex2Bytes("0x0000000000000000000000003d0768da09ce77d25e2d998e6a7b6ed4b9116c2d")
		assert.Equal(t, expectedResult, result)
	})

	t.Run("GetGasPrice", func(t *testing.T) {
		price, err := client.Proxy.GasPrice(ctx)
		require.NoError(t, err)

		expectedPrice := big.NewInt(18000000000)
		assert.Equal(t, 0, expectedPrice.Cmp(price))
	})

	t.Run("EstimateGas", func(t *testing.T) {
		gas, err := client.Proxy.EstimateGas(ctx, &etherscan.EstimateGasRequest{
			Data:     common.Hex2Bytes("0x4e71d92d"),
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

type mockProxyAPI struct {
	apiKey string
}

func newProxyAPI() mockProxyAPI {
	return mockProxyAPI{apiKey: uuid.NewString()}
}

func (m mockProxyAPI) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/api" {
		http.Error(w, "path not found", http.StatusNotFound)
		return
	}

	q := req.URL.Query()
	if q.Get("module") != "proxy" {
		http.Error(w, "unknown module", http.StatusNotFound)
		return
	}

	if q.Get("apikey") != m.apiKey {
		http.Error(w, "unknown API key", http.StatusForbidden)
		return
	}

	switch q.Get("action") {
	case "eth_blockNumber":
		m.handleBlockNumber(w, q)

	case "eth_getBlockByNumber":
		m.handleGetBlockByNumber(w, q)

	case "eth_getUncleByBlockNumberAndIndex":
		m.handleGetUncle(w, q)

	case "eth_getBlockTransactionCountByNumber":
		m.handleGetBlockTxCount(w, q)

	case "eth_getTransactionByHash":
		m.handleGetTxByHash(w, q)

	case "eth_getTransactionByBlockNumberAndIndex":
		m.handleGetTxByNumber(w, q)

	case "eth_getTransactionCount":
		m.handleGetTransactionCount(w, q)

	case "eth_getTransactionReceipt":
		m.handleGetTransactionReceipt(w, q)

	case "eth_call":
		m.handleCall(w, q)

	case "eth_getCode":
		m.handleGetCode(w, q)

	case "eth_getStorageAt":
		m.handleGetStorage(w, q)

	case "eth_gasPrice":
		m.handleGasPrice(w, q)

	case "eth_estimateGas":
		m.handleEstimateGas(w, q)

	default:
		http.Error(w, "unknown action", http.StatusNotFound)
	}
}

const blockNumberResponse = `{
	"jsonrpc":"2.0",
	"id":83,
	"result":"0xc36b29"
}`

func (m mockProxyAPI) handleBlockNumber(w http.ResponseWriter, _ url.Values) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(blockNumberResponse))
}

const blockByNumberFullResponse = `{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
	"difficulty": "0x1d95715bd14",
	"extraData": "0x",
	"gasLimit": "0x2fefd8",
	"gasUsed": "0x5208",
	"hash": "0x7eb7c23a5ac2f2d70aa1ba4e5c56d89de5ac993590e5f6e79c394e290d998ba8",
	"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	"miner": "0xf927a40c8b7f6e07c5af7fa2155b4864a4112b13",
	"mixHash": "0x13dd2c8aec729f75aebcd79a916ecb0f7edc6493efcc6a4da8d7b0ab3ee88444",
	"nonce": "0xc60a782e2e69ce22",
	"number": "0x10d4f",
	"parentHash": "0xf8d01370e6e274f8188954fbee435b40c35b2ad3d4ab671f6d086cd559e48f04",
	"receiptsRoot": "0x0c44b7ed0fefb613ec256341aa0ffdb643e869e3a0ebc8f58e36b4e47efedd33",
	"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
	"size": "0x275",
	"stateRoot": "0xd64a0f63e2c7f541e6e6f8548a10a5c4e49fda7ac1aa80f9dddef648c7b9e25f",
	"timestamp": "0x55c9ea07",
	"totalDifficulty": "0x120d56f6821b170",
	"transactions": [
	  {
		"blockHash": "0x7eb7c23a5ac2f2d70aa1ba4e5c56d89de5ac993590e5f6e79c394e290d998ba8",
		"blockNumber": "0x10d4f",
		"from": "0x4458f86353b4740fe9e09071c23a7437640063c9",
		"gas": "0x5208",
		"gasPrice": "0xba43b7400",
		"hash": "0xa442249820de6be754da81eafbd44a865773e4b23d7c0522d31fd03977823008",
		"input": "0x",
		"nonce": "0x1",
		"to": "0xbf3403210f9802205f426759947a80a9fda71b1e",
		"transactionIndex": "0x0",
		"value": "0xaa9f075c200000",
		"type": "0x0",
		"v": "0x1b",
		"r": "0x2c2789c6704ba2606e200e1ba4fd17ba4f0e0f94abe32a12733708c3d3442616",
		"s": "0x2946f47e3ece580b5b5ecb0f8c52604fa5f60aeb4103fc73adcbf6d620f9872b"
	  }
	],
	"transactionsRoot": "0x4a5b78c13d11559c9541576834b5172fe8b18507c0f9f76454fcdddedd8dff7a",
	"uncles": []
  }
}`

const blockByNumberSummaryResponse = `{
	"jsonrpc": "2.0",
	"id": 1,
	"result": {
	  "difficulty": "0x1d95715bd14",
	  "extraData": "0x",
	  "gasLimit": "0x2fefd8",
	  "gasUsed": "0x5208",
	  "hash": "0x7eb7c23a5ac2f2d70aa1ba4e5c56d89de5ac993590e5f6e79c394e290d998ba8",
	  "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	  "miner": "0xf927a40c8b7f6e07c5af7fa2155b4864a4112b13",
	  "mixHash": "0x13dd2c8aec729f75aebcd79a916ecb0f7edc6493efcc6a4da8d7b0ab3ee88444",
	  "nonce": "0xc60a782e2e69ce22",
	  "number": "0x10d4f",
	  "parentHash": "0xf8d01370e6e274f8188954fbee435b40c35b2ad3d4ab671f6d086cd559e48f04",
	  "receiptsRoot": "0x0c44b7ed0fefb613ec256341aa0ffdb643e869e3a0ebc8f58e36b4e47efedd33",
	  "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
	  "size": "0x275",
	  "stateRoot": "0xd64a0f63e2c7f541e6e6f8548a10a5c4e49fda7ac1aa80f9dddef648c7b9e25f",
	  "timestamp": "0x55c9ea07",
	  "totalDifficulty": "0x120d56f6821b170",
	  "transactions": [
		"0xa442249820de6be754da81eafbd44a865773e4b23d7c0522d31fd03977823008"
	  ],
	  "transactionsRoot": "0x4a5b78c13d11559c9541576834b5172fe8b18507c0f9f76454fcdddedd8dff7a",
	  "uncles": []
	}
}`

func (m mockProxyAPI) handleGetBlockByNumber(w http.ResponseWriter, q url.Values) {
	if q.Get("tag") != "0x10d4f" {
		http.Error(w, "unexpected tag", http.StatusBadRequest)
		return
	}

	switch q.Get("boolean") {
	case "true":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(blockByNumberFullResponse))

	case "false":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(blockByNumberSummaryResponse))

	default:
		http.Error(w, "unexpected boolean", http.StatusBadRequest)
	}
}

const uncleByBlockNumberAndIndexResponse = `{
	"jsonrpc":"2.0",
	"id":1,
	"result":{
	   "baseFeePerGas":"0x65a42b13c",
	   "difficulty":"0x1b1457a8247bbb",
	   "extraData":"0x486976656f6e2063612d68656176792059476f6e",
	   "gasLimit":"0x1ca359a",
	   "gasUsed":"0xb48fe1",
	   "hash":"0x1da88e3581315d009f1cb600bf06f509cd27a68cb3d6437bda8698d04089f14a",
	   "logsBloom":"0xf1a360ca505cdda510d810c1c81a03b51a8a508ed601811084833072945290235c8721e012182e40d57df552cf00f1f01bc498018da19e008681832b43762a30c26e11709948a9b96883a42ad02568e3fcc3000004ee12813e4296498261619992c40e22e60bd95107c5bd8462fcca570a0095d52a4c24720b00f13a2c3d62aca81e852017470c109643b15041fd69742406083d67654fc841a18b405ab380e06a8c14c0138b6602ea8f48b2cd90ac88c3478212011136802900264718a085047810221225080dfb2c214010091a6f233883bb0084fa1c197330a10bb0006686e678b80e50e4328000041c218d1458880181281765d28d51066058f3f80a7822",
	   "miner":"0x1ad91ee08f21be3de0ba2ba6918e714da6b45836",
	   "mixHash":"0xa8e1dbbf073614c7ed05f44b9e92fbdb3e1d52575ed8167fa57f934210bbb0a2",
	   "nonce":"0x28cc3e5b7bee9866",
	   "number":"0xc63274",
	   "parentHash":"0x496dae3e722efdd9ee1eb69499bdc7ed0dca54e13cd1157a42811c442f01941f",
	   "receiptsRoot":"0x9c9a7a99b4af7607691a7f2a50d474290385c0a6f39c391131ea0c67307213f4",
	   "sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
	   "size":"0x224",
	   "stateRoot":"0xde9a11f0ee321390c1a7843cab7b9ffd3779d438bc8f77de4361dfe2807d7dee",
	   "timestamp":"0x6110bd1a",
	   "transactionsRoot":"0xa04a79e531db3ec373cb63e9ebfbc9c95525de6347958918a273675d4f221575",
	   "uncles":[

	   ]
	}
}`

func (m mockProxyAPI) handleGetUncle(w http.ResponseWriter, q url.Values) {
	if q.Get("tag") != "0xC63276" {
		http.Error(w, "unexpected tag", http.StatusBadRequest)
		return
	}

	if q.Get("index") != "0x0" {
		http.Error(w, "unexpected index", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(uncleByBlockNumberAndIndexResponse))
}

const blockTransactionCountResponse = `{
	"jsonrpc":"2.0",
	"id":1,
	"result":"0x3"
}`

func (m mockProxyAPI) handleGetBlockTxCount(w http.ResponseWriter, q url.Values) {
	if q.Get("tag") != "0x10FB78" {
		http.Error(w, "unexpected tag", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(blockTransactionCountResponse))
}

const txByHashResponse = `{
	"jsonrpc":"2.0",
	"result":{
	   "accessList":[

	   ],
	   "blockHash":"0xdce94191f861842c2786e3594da0c0109707fd78409cab5f38e10eb87d0f301c",
	   "blockNumber":"0xa36e44",
	   "chainId":"0x3",
	   "condition":null,
	   "creates":null,
	   "from":"0xb910ae1db14a9fbc64ce175bdca6d3a743f690ab",
	   "gas":"0x186a0",
	   "gasPrice":"0x3b9aca09",
	   "hash":"0xf96ff62ba5aaf46cd824b6766f7fa6f6b9595b1dd4ef1d31bcf1f765047c2835",
	   "input":"0xd0e30db0",
	   "maxFeePerGas":"0x3b9aca12",
	   "maxPriorityFeePerGas":"0x3b9aca00",
	   "nonce":"0xc6",
	   "publicKey":"0x6dbf7068e19de8457c426a758a92ea54827ebd5b8467c3a1a5c4ac19bc7570457738fe496a40ea4e1f59d39d89636a430afdec0bf2a8060c6bf7d612bfe90ad3",
	   "r":"0xdecdc48821a06bf116e82b355d520dc5a44d6df98234e5344c16565b0b3dfdba",
	   "raw":"0x02f8750381c6843b9aca00843b9aca12830186a094c778417e063141139fce010982780140aa0cd5ab8502540be40084d0e30db0c001a0decdc48821a06bf116e82b355d520dc5a44d6df98234e5344c16565b0b3dfdbaa06b85bb6fd8153e86b50f0011787585e8c709a2a25e7ee3c2579572f07acfd42e",
	   "s":"0x6b85bb6fd8153e86b50f0011787585e8c709a2a25e7ee3c2579572f07acfd42e",
	   "to":"0xc778417e063141139fce010982780140aa0cd5ab",
	   "transactionIndex":"0xd",
	   "type":"0x2",
	   "v":"0x1",
	   "value":"0x2540be400"
	},
	"id":1
}`

func (m mockProxyAPI) handleGetTxByHash(w http.ResponseWriter, q url.Values) {
	txHash := common.HexToHash(q.Get("txhash"))
	if txHash != common.HexToHash("0x1e2910a262b1008d0616a0beb24c1a491d78771baa54a33e66065e03b1f46bc1") {
		http.Error(w, "unexpected txhash", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(txByHashResponse))
}

const txByNumberResponse = `{
	"jsonrpc":"2.0",
	"result":{
	   "accessList":[

	   ],
	   "blockHash":"0xdce94191f861842c2786e3594da0c0109707fd78409cab5f38e10eb87d0f301c",
	   "blockNumber":"0xa36e44",
	   "chainId":"0x3",
	   "condition":null,
	   "creates":null,
	   "from":"0xb910ae1db14a9fbc64ce175bdca6d3a743f690ab",
	   "gas":"0x186a0",
	   "gasPrice":"0x3b9aca09",
	   "hash":"0xf96ff62ba5aaf46cd824b6766f7fa6f6b9595b1dd4ef1d31bcf1f765047c2835",
	   "input":"0xd0e30db0",
	   "maxFeePerGas":"0x3b9aca12",
	   "maxPriorityFeePerGas":"0x3b9aca00",
	   "nonce":"0xc6",
	   "publicKey":"0x6dbf7068e19de8457c426a758a92ea54827ebd5b8467c3a1a5c4ac19bc7570457738fe496a40ea4e1f59d39d89636a430afdec0bf2a8060c6bf7d612bfe90ad3",
	   "r":"0xdecdc48821a06bf116e82b355d520dc5a44d6df98234e5344c16565b0b3dfdba",
	   "raw":"0x02f8750381c6843b9aca00843b9aca12830186a094c778417e063141139fce010982780140aa0cd5ab8502540be40084d0e30db0c001a0decdc48821a06bf116e82b355d520dc5a44d6df98234e5344c16565b0b3dfdbaa06b85bb6fd8153e86b50f0011787585e8c709a2a25e7ee3c2579572f07acfd42e",
	   "s":"0x6b85bb6fd8153e86b50f0011787585e8c709a2a25e7ee3c2579572f07acfd42e",
	   "to":"0xc778417e063141139fce010982780140aa0cd5ab",
	   "transactionIndex":"0xd",
	   "type":"0x2",
	   "v":"0x1",
	   "value":"0x2540be400"
	},
	"id":1
}`

func (m mockProxyAPI) handleGetTxByNumber(w http.ResponseWriter, q url.Values) {
	if q.Get("tag") != "0xC6331D" {
		http.Error(w, "unexpected tag", http.StatusBadRequest)
		return
	}

	if q.Get("index") != "0x11A" {
		http.Error(w, "unexpected index", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(txByNumberResponse))
}

const getTxCountResponse = `{
	"jsonrpc":"2.0",
	"id":1,
	"result":"0x44"
}`

func (m mockProxyAPI) handleGetTransactionCount(w http.ResponseWriter, q url.Values) {
	address := common.HexToAddress(q.Get("address"))
	if address != common.HexToAddress("0x4bd5900Cb274ef15b153066D736bf3e83A9ba44e") {
		http.Error(w, "unexpected address", http.StatusBadRequest)
		return
	}

	if q.Get("tag") != "latest" {
		http.Error(w, "unexpected tag", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(getTxCountResponse))
}

const getTxReceiptResponse = `{
	"jsonrpc":"2.0",
	"result":{
	   "blockHash":"0xdce94191f861842c2786e3594da0c0109707fd78409cab5f38e10eb87d0f301c",
	   "blockNumber":"0xa36e44",
	   "contractAddress":null,
	   "cumulativeGasUsed":"0x23e989",
	   "effectiveGasPrice":"0x3b9aca09",
	   "from":"0xb910ae1db14a9fbc64ce175bdca6d3a743f690ab",
	   "gasUsed":"0x6d22",
	   "logs":[
		  {
			 "address":"0xc778417e063141139fce010982780140aa0cd5ab",
			 "blockHash":"0xdce94191f861842c2786e3594da0c0109707fd78409cab5f38e10eb87d0f301c",
			 "blockNumber":"0xa36e44",
			 "data":"0x00000000000000000000000000000000000000000000000000000002540be400",
			 "logIndex":"0xd",
			 "removed":false,
			 "topics":[
				"0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c",
				"0x000000000000000000000000b910ae1db14a9fbc64ce175bdca6d3a743f690ab"
			 ],
			 "transactionHash":"0xf96ff62ba5aaf46cd824b6766f7fa6f6b9595b1dd4ef1d31bcf1f765047c2835",
			 "transactionIndex":"0xd",
			 "transactionLogIndex":"0x0",
			 "type":"mined"
		  }
	   ],
	   "logsBloom":"0x00000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001400000000000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000400000000000000000",
	   "status":"0x1",
	   "to":"0xc778417e063141139fce010982780140aa0cd5ab",
	   "transactionHash":"0xf96ff62ba5aaf46cd824b6766f7fa6f6b9595b1dd4ef1d31bcf1f765047c2835",
	   "transactionIndex":"0xd",
	   "type":"0x2"
	},
	"id":1
}`

func (m mockProxyAPI) handleGetTransactionReceipt(w http.ResponseWriter, q url.Values) {
	txHash := common.HexToHash(q.Get("txhash"))
	if txHash != common.HexToHash("0x1e2910a262b1008d0616a0beb24c1a491d78771baa54a33e66065e03b1f46bc1") {
		http.Error(w, "unexpected txhash", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(getTxReceiptResponse))
}

const callResponse = `{
	"jsonrpc":"2.0",
	"id":1,
	"result":"0x00000000000000000000000000000000000000000000000000601d8888141c00"
}`

func (m mockProxyAPI) handleCall(w http.ResponseWriter, q url.Values) {
	toAddr := common.HexToAddress(q.Get("to"))
	if toAddr != common.HexToAddress("0xAEEF46DB4855E25702F8237E8f403FddcaF931C0") {
		http.Error(w, "unexpected to address", http.StatusBadRequest)
		return
	}

	if q.Get("data") != "0x70a08231000000000000000000000000e16359506c028e51f16be38986ec5746251e9724" {
		http.Error(w, "unexpected data", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(callResponse))
}

const getCodeResponse = `{
	"jsonrpc":"2.0",
	"id":1,
	"result":"0x3660008037602060003660003473273930d21e01ee25e4c219b63259d214872220a261235a5a03f21560015760206000f3"
}`

func (m mockProxyAPI) handleGetCode(w http.ResponseWriter, q url.Values) {
	address := common.HexToAddress(q.Get("address"))
	if address != common.HexToAddress("0xf75e354c5edc8efed9b59ee9f67a80845ade7d0c") {
		http.Error(w, "unexpected address", http.StatusBadRequest)
		return
	}

	if q.Get("tag") != "latest" {
		http.Error(w, "unexpected tag", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(getCodeResponse))
}

const getStorageResponse = `{
	"jsonrpc":"2.0",
	"id":1,
	"result":"0x0000000000000000000000003d0768da09ce77d25e2d998e6a7b6ed4b9116c2d"
}`

func (m mockProxyAPI) handleGetStorage(w http.ResponseWriter, q url.Values) {
	address := common.HexToAddress(q.Get("address"))
	if address != common.HexToAddress("0x6e03d9cce9d60f3e9f2597e13cd4c54c55330cfd") {
		http.Error(w, "unexpected address", http.StatusBadRequest)
		return
	}

	if q.Get("position") != "0x0" {
		http.Error(w, "unexpected position", http.StatusBadRequest)
		return
	}

	if q.Get("tag") != "latest" {
		http.Error(w, "unexpected tag", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(getStorageResponse))
}

const gasPriceResponse = `{
	"jsonrpc":"2.0",
	"id":73,
	"result":"0x430e23400"
}`

func (m mockProxyAPI) handleGasPrice(w http.ResponseWriter, q url.Values) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(gasPriceResponse))
}

const estimateGasResponse = `{
	"jsonrpc":"2.0",
	"id":1,
	"result":"0x6556"
}`

func (m mockProxyAPI) handleEstimateGas(w http.ResponseWriter, q url.Values) {
	if q.Get("data") != "0x4e71d92d" {
		http.Error(w, "unexpected data", http.StatusBadRequest)
		return
	}

	address := common.HexToAddress(q.Get("address"))
	if address != common.HexToAddress("0xf0160428a8552ac9bb7e050d90eeade4ddd52843") {
		http.Error(w, "unexpected address", http.StatusBadRequest)
		return
	}

	if q.Get("value") != "0xff22" {
		http.Error(w, "unexpected value", http.StatusBadRequest)
		return
	}

	if q.Get("gasPrice") != "0x51da038cc" {
		http.Error(w, "unexpected gas price", http.StatusBadRequest)
		return
	}

	if q.Get("gas") != "0x5f5e0ff" {
		http.Error(w, "unexpected gas", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(estimateGasResponse))
}
