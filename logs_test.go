package etherscan_test

import (
	"context"
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

func TestLogs(t *testing.T) {
	m := newLogsAPI()

	ts := httptest.NewServer(m)
	t.Cleanup(ts.Close)

	u, err := url.Parse(ts.URL)
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

		assert.Equal(
			t,
			common.HexToHash("0x0b03498648ae2da924f961dda00dc6bb0a8df15519262b7e012b7d67f4bb7e83"),
			logs[0].TransactionHash,
		)

		assert.Equal(
			t,
			common.HexToHash("0x8c72ea19b48947c4339077bd9c9c09a780dfbdb1cafe68db4d29cdf2754adc11"),
			logs[1].TransactionHash,
		)

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

			assert.Equal(
				t,
				common.HexToHash("0x0b03498648ae2da924f961dda00dc6bb0a8df15519262b7e012b7d67f4bb7e83"),
				logs[0].TransactionHash,
			)
		})
	})
}

type mockLogsAPI struct {
	apiKey string
}

func newLogsAPI() mockLogsAPI {
	return mockLogsAPI{apiKey: uuid.NewString()}
}

func (m mockLogsAPI) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/api" {
		http.Error(w, "path not found", http.StatusNotFound)
		return
	}

	q := req.URL.Query()
	if q.Get("module") != "logs" {
		http.Error(w, "unknown module", http.StatusNotFound)
		return
	}

	if q.Get("apikey") != m.apiKey {
		http.Error(w, "unknown API key", http.StatusForbidden)
		return
	}

	if q.Get("action") != "getLogs" {
		http.Error(w, "unknown action", http.StatusNotFound)
		return
	}

	m.handleGetLogs(w, q)
}

func (m mockLogsAPI) handleGetLogs(w http.ResponseWriter, q url.Values) {
	if q.Get("fromBlock") != "379224" {
		http.Error(w, "unexpected fromBlock", http.StatusBadRequest)
		return
	}

	address := common.HexToAddress(q.Get("address"))
	if address != common.HexToAddress("0x33990122638b9132ca29c723bdf037f1a891a70c") {
		http.Error(w, "unexpected address", http.StatusBadRequest)
		return
	}

	if q.Get("topic0") != "0xf63780e752c6a54a94fc52715dbc5518a3b4c3c2833d301a204226548a2a8545" {
		http.Error(w, "unexpected topic 0", http.StatusBadRequest)
		return
	}

	switch q.Get("toBlock") {
	case "latest":
		m.handleGetLogsToLatest(w, q)

	case "400000":
		m.handleGetLogsToFixed(w, q)

	default:
		http.Error(w, "unexpected toBlock", http.StatusBadRequest)
	}
}

const getLogsLatestResponse = `{
  "message": "OK",
  "status": "1",
  "result": [
	{
	  "address": "0x33990122638b9132ca29c723bdf037f1a891a70c",
	  "blockNumber": "0x5c958",
	  "data": "0x",
	  "gasPrice": "0xba43b7400",
	  "gasUsed": "0x10682",
	  "logIndex": "0x",
	  "timeStamp": "0x561d688c",
	  "topics": [
		"0xf63780e752c6a54a94fc52715dbc5518a3b4c3c2833d301a204226548a2a8545",
		"0x72657075746174696f6e00000000000000000000000000000000000000000000",
		"0x000000000000000000000000d9b2f59f3b5c7b3c67047d2f03c3e8052470be92"
	  ],
	  "transactionHash": "0x0b03498648ae2da924f961dda00dc6bb0a8df15519262b7e012b7d67f4bb7e83",
	  "transactionIndex": "0x"
	},
	{
	  "address": "0x33990122638b9132ca29c723bdf037f1a891a70c",
	  "blockNumber": "0x5c965",
	  "data": "0x",
	  "gasPrice": "0xba43b7400",
	  "gasUsed": "0x105c2",
	  "logIndex": "0x",
	  "timeStamp": "0x561d6930",
	  "topics": [
		  "0xf63780e752c6a54a94fc52715dbc5518a3b4c3c2833d301a204226548a2a8545",
		  "0x6c6f747465727900000000000000000000000000000000000000000000000000",
		  "0x0000000000000000000000001f6cc3f7c927e1196c03ac49c5aff0d39c9d103d"
	  ],
	  "transactionHash": "0x8c72ea19b48947c4339077bd9c9c09a780dfbdb1cafe68db4d29cdf2754adc11",
	  "transactionIndex": "0x"
	}
  ]
}`

func (m mockLogsAPI) handleGetLogsToLatest(w http.ResponseWriter, q url.Values) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(getLogsLatestResponse)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const getLogsToFixedResponse = `{
  "status": "1",
  "message": "OK",
  "result": [
	{
	  "address": "0x33990122638b9132ca29c723bdf037f1a891a70c",
	  "topics": [
		"0xf63780e752c6a54a94fc52715dbc5518a3b4c3c2833d301a204226548a2a8545",
		"0x72657075746174696f6e00000000000000000000000000000000000000000000",
		"0x000000000000000000000000d9b2f59f3b5c7b3c67047d2f03c3e8052470be92"
	  ],
	  "data": "0x",
	  "blockNumber": "0x5c958",
	  "timeStamp": "0x561d688c",
	  "gasPrice": "0xba43b7400",
	  "gasUsed": "0x10682",
	  "logIndex": "0x",
	  "transactionHash": "0x0b03498648ae2da924f961dda00dc6bb0a8df15519262b7e012b7d67f4bb7e83",
	  "transactionIndex": "0x"
	}
  ]
}`

func (m mockLogsAPI) handleGetLogsToFixed(w http.ResponseWriter, q url.Values) {
	if q.Get("topic0_1_opr") != "and" {
		http.Error(w, "unexpected topic0_1_opr", http.StatusBadRequest)
		return
	}

	if q.Get("topic1") != "0x72657075746174696f6e00000000000000000000000000000000000000000000" {
		http.Error(w, "unexpected topic1", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(getLogsToFixedResponse)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
