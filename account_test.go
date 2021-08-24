package etherscan_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/ryanc414/etherscan-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	m := &mockAccountsAPI{apiKey: uuid.NewString()}
	ts := httptest.NewServer(m)
	t.Cleanup(ts.Close)

	u, err := url.Parse(ts.URL)
	require.NoError(t, err)

	client := etherscan.New(&etherscan.Params{
		APIKey:  m.apiKey,
		BaseURL: u,
	})

	ctx := context.Background()

	t.Run("TestGetETHBalance", func(t *testing.T) {
		bal, err := client.Accounts.GetETHBalance(ctx, &etherscan.ETHBalanceRequest{
			Address: common.HexToAddress("0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae"),
			Tag:     etherscan.BlockParameterLatest,
		})
		require.NoError(t, err)
		assert.Equal(t, "40891626854930000000000", bal.String())
	})

	t.Run("TestMultiGetETHBalance", func(t *testing.T) {
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

	t.Run("TestListNormalTxs", func(t *testing.T) {
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
}

type mockAccountsAPI struct {
	apiKey string
}

func (m *mockAccountsAPI) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/api" {
		http.Error(w, "path not found", http.StatusNotFound)
		return
	}

	q := req.URL.Query()
	if q.Get("module") != "account" {
		http.Error(w, "unknown model", http.StatusNotFound)
		return
	}

	if q.Get("apikey") != m.apiKey {
		http.Error(w, "unknown API key", http.StatusForbidden)
		return
	}

	switch q.Get("action") {
	case "balance":
		m.handleGetBalance(w, q)

	case "balancemulti":
		m.handleGetMultiBalance(w, q)

	case "txlist":
		m.handleTxList(w, q)

	case "txlistinternal":
		m.handleInternalTxList(w, q)

	default:
		http.Error(w, "unknown action", http.StatusNotFound)
	}
}

const getBalanceResponse = `{
	"status":"1",
	"message":"OK",
	"result":"40891626854930000000000"
}`

func (m *mockAccountsAPI) handleGetBalance(w http.ResponseWriter, q url.Values) {
	address := common.HexToAddress(q.Get("address"))
	if address != common.HexToAddress("0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae") {
		http.Error(w, "unknown address", http.StatusBadRequest)
		return
	}

	if q.Get("tag") != "latest" {
		http.Error(w, "unknown tag", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(getBalanceResponse))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const getMultiBalanceResponse = `{
	"status":"1",
	"message":"OK",
	"result":[
	   {
		  "account":"0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a",
		  "balance":"40891626854930000000000"
	   },
	   {
		  "account":"0x63a9975ba31b0b9626b34300f7f627147df1f526",
		  "balance":"332567136222827062478"
	   },
	   {
		  "account":"0x198ef1ec325a96cc354c7266a038be8b5c558f67",
		  "balance":"0"
	   }
	]
}`

var multiETHBalAddrs = []common.Address{
	common.HexToAddress("0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a"),
	common.HexToAddress("0x63a9975ba31b0b9626b34300f7f627147df1f526"),
	common.HexToAddress("0x198ef1ec325a96cc354c7266a038be8b5c558f67"),
}

func (m *mockAccountsAPI) handleGetMultiBalance(w http.ResponseWriter, q url.Values) {
	addresses := strings.Split(q.Get("address"), ",")
	if len(addresses) != 3 {
		http.Error(w, "unexpected number of addresses", http.StatusBadRequest)
		return
	}

	for i := range addresses {
		if common.HexToAddress(addresses[i]) != multiETHBalAddrs[i] {
			http.Error(w, "unknown address", http.StatusBadRequest)
			return
		}
	}

	if q.Get("tag") != "latest" {
		http.Error(w, "unknown tag", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(getMultiBalanceResponse))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const listNormalTxResponse = `{
	"status":"1",
	"message":"OK",
	"result":[
	   {
		  "blockNumber":"0",
		  "timeStamp":"1438269973",
		  "hash":"GENESIS_ddbd2b932c763ba5b1b7ae3b362eac3e8d40121a",
		  "nonce":"",
		  "blockHash":"",
		  "transactionIndex":"0",
		  "from":"GENESIS",
		  "to":"0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a",
		  "value":"10000000000000000000000",
		  "gas":"0",
		  "gasPrice":"0",
		  "isError":"0",
		  "txreceipt_status":"",
		  "input":"",
		  "contractAddress":"",
		  "cumulativeGasUsed":"0",
		  "gasUsed":"0",
		  "confirmations":"12698061"
	   },
	   {
		  "blockNumber":"47884",
		  "timeStamp":"1438947953",
		  "hash":"0xad1c27dd8d0329dbc400021d7477b34ac41e84365bd54b45a4019a15deb10c0d",
		  "nonce":"0",
		  "blockHash":"0xf2988b9870e092f2898662ccdbc06e0e320a08139e9c6be98d0ce372f8611f22",
		  "transactionIndex":"0",
		  "from":"0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a",
		  "to":"0x2910543af39aba0cd09dbb2d50200b3e800a63d2",
		  "value":"5000000000000000000",
		  "gas":"23000",
		  "gasPrice":"400000000000",
		  "isError":"0",
		  "txreceipt_status":"",
		  "input":"0x454e34354139455138",
		  "contractAddress":"",
		  "cumulativeGasUsed":"21612",
		  "gasUsed":"21612",
		  "confirmations":"12650177"
	   }
	]
}`

func (m *mockAccountsAPI) handleTxList(w http.ResponseWriter, q url.Values) {
	address := common.HexToAddress(q.Get("address"))
	if address != common.HexToAddress("0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a") {
		http.Error(w, "unknown address", http.StatusBadRequest)
		return
	}

	if q.Get("startblock") != "0" || q.Get("endblock") != "99999999" {
		http.Error(w, "unexpected block params", http.StatusBadRequest)
		return
	}

	if q.Get("sort") != "asc" {
		http.Error(w, "unexpected sort param", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(listNormalTxResponse))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (m *mockAccountsAPI) handleInternalTxList(w http.ResponseWriter, q url.Values) {
	if q.Get("address") != "" {
		m.handleInternalTxListAddress(w, q)
		return
	}

	if q.Get("txhash") != "" {
		m.handleInternalTxListHash(w, q)
		return
	}

	m.handleInternalTxListBlockRange(w, q)
}

const listInternalTxByAddressResponse = `{
	"status":"1",
	"message":"OK",
	"result":[
	   {
		  "blockNumber":"2535368",
		  "timeStamp":"1477837690",
		  "hash":"0x8a1a9989bda84f80143181a68bc137ecefa64d0d4ebde45dd94fc0cf49e70cb6",
		  "from":"0x20d42f2e99a421147acf198d775395cac2e8b03d",
		  "to":"",
		  "value":"0",
		  "contractAddress":"0x2c1ba59d6f58433fb1eaee7d20b26ed83bda51a3",
		  "input":"",
		  "type":"create",
		  "gas":"254791",
		  "gasUsed":"46750",
		  "traceId":"0",
		  "isError":"0",
		  "errCode":""
	   },
	   {
		  "blockNumber":"2535479",
		  "timeStamp":"1477839134",
		  "hash":"0x1a50f1dc0bc912745f7d09b988669f71d199719e2fb7592c2074ede9578032d0",
		  "from":"0x2c1ba59d6f58433fb1eaee7d20b26ed83bda51a3",
		  "to":"0x20d42f2e99a421147acf198d775395cac2e8b03d",
		  "value":"100000000000000000",
		  "contractAddress":"",
		  "input":"",
		  "type":"call",
		  "gas":"235231",
		  "gasUsed":"0",
		  "traceId":"0",
		  "isError":"0",
		  "errCode":""
	   }
	]
}`

func (m *mockAccountsAPI) handleInternalTxListAddress(w http.ResponseWriter, q url.Values) {
	address := common.HexToAddress(q.Get("address"))
	if address != common.HexToAddress("0x2c1ba59d6f58433fb1eaee7d20b26ed83bda51a3") {
		http.Error(w, "unknown address", http.StatusBadRequest)
		return
	}

	if q.Get("startblock") != "0" || q.Get("endblock") != "99999999" {
		http.Error(w, "unexpected block params", http.StatusBadRequest)
		return
	}

	if q.Get("sort") != "asc" {
		http.Error(w, "unexpected sort param", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(listInternalTxByAddressResponse))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const listInternalTxByHashResponse = `{
	"status":"1",
	"message":"OK",
	"result":[
	   {
		  "blockNumber":"1743059",
		  "timeStamp":"1466489498",
		  "from":"0x2cac6e4b11d6b58f6d3c1c9d5fe8faa89f60e5a2",
		  "to":"0x66a1c3eaf0f1ffc28d209c0763ed0ca614f3b002",
		  "value":"7106740000000000",
		  "contractAddress":"",
		  "input":"",
		  "type":"call",
		  "gas":"2300",
		  "gasUsed":"0",
		  "isError":"0",
		  "errCode":""
	   }
	]
}`

func (m *mockAccountsAPI) handleInternalTxListHash(w http.ResponseWriter, q url.Values) {
	hash := common.HexToHash(q.Get("txhash"))
	if hash != common.HexToHash("0x40eb908387324f2b575b4879cd9d7188f69c8fc9d87c901b9e2daaea4b442170") {
		http.Error(w, "unknown tx hash", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(listInternalTxByHashResponse))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const listInternalTxBlockRangeResponse = `{
	"status":"1",
	"message":"OK",
	"result":[
	   {
		  "blockNumber":"50107",
		  "timeStamp":"1438984016",
		  "hash":"0x3f97c969ddf71f515ce5373b1f8e76e9fd7016611d8ce455881009414301789e",
		  "from":"0x109c4f2ccc82c4d77bde15f306707320294aea3f",
		  "to":"0x881b0a4e9c55d08e31d8d3c022144d75a454211c",
		  "value":"1000000000000000000",
		  "contractAddress":"",
		  "input":"",
		  "type":"call",
		  "gas":"2300",
		  "gasUsed":"0",
		  "traceId":"0",
		  "isError":"1",
		  "errCode":""
	   },
	   {
		  "blockNumber":"50111",
		  "timeStamp":"1438984075",
		  "hash":"0x893c428fed019404f704cf4d9be977ed9ca01050ed93dccdd6c169422155586f",
		  "from":"0x109c4f2ccc82c4d77bde15f306707320294aea3f",
		  "to":"0x881b0a4e9c55d08e31d8d3c022144d75a454211c",
		  "value":"1000000000000000000",
		  "contractAddress":"",
		  "input":"",
		  "type":"call",
		  "gas":"2300",
		  "gasUsed":"0",
		  "traceId":"0",
		  "isError":"0",
		  "errCode":""
	   }
	]
}`

func (m *mockAccountsAPI) handleInternalTxListBlockRange(w http.ResponseWriter, q url.Values) {
	if q.Get("startblock") != "0" || q.Get("endblock") != "2702578" {
		http.Error(w, "unexpected block parameters", http.StatusBadRequest)
		return
	}

	if q.Get("sort") != "asc" {
		http.Error(w, "unexpected sort param", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(listInternalTxBlockRangeResponse))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
