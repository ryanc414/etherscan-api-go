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

func TestTransactions(t *testing.T) {
	m := newMockTransactionsAPI()

	ts := httptest.NewServer(m)
	t.Cleanup(ts.Close)

	u, err := url.Parse(ts.URL)
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

		assert.True(t, status.IsError)
		assert.Equal(t, "Bad jump destination", status.ErrDescription)
	})

	t.Run("GetTxReceiptStatus", func(t *testing.T) {
		txHash := common.HexToHash("0x513c1ba0bebf66436b5fed86ab668452b7805593c05073eb2d51d3a52f480a76")
		status, err := client.Transactions.GetTxReceiptStatus(ctx, txHash)
		require.NoError(t, err)
		assert.True(t, status)
	})
}

type mockTransactionsAPI struct {
	apiKey string
}

func newMockTransactionsAPI() mockTransactionsAPI {
	return mockTransactionsAPI{apiKey: uuid.NewString()}
}

func (m mockTransactionsAPI) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/api" {
		http.Error(w, "path not found", http.StatusNotFound)
		return
	}

	q := req.URL.Query()
	if q.Get("module") != "transaction" {
		http.Error(w, "unknown model", http.StatusNotFound)
		return
	}

	if q.Get("apikey") != m.apiKey {
		http.Error(w, "unknown API key", http.StatusForbidden)
		return
	}

	switch q.Get("action") {
	case "getstatus":
		m.handleGetStatus(w, q)

	case "gettxreceiptstatus":
		m.handleReceiptStatus(w, q)

	default:
		http.Error(w, "unknown action", http.StatusNotFound)
	}
}

const getStatusResponse = `{
	"status":"1",
	"message":"OK",
	"result":{
	   "isError":"1",
	   "errDescription":"Bad jump destination"
	}
}`

func (m mockTransactionsAPI) handleGetStatus(w http.ResponseWriter, q url.Values) {
	txHash := common.HexToHash(q.Get("txhash"))
	if txHash != common.HexToHash("0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a") {
		http.Error(w, "unknown address", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(getStatusResponse)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const getReceiptResponse = `{
	"status":"1",
	"message":"OK",
	"result":{
	   "status":"1"
	}
}`

func (m mockTransactionsAPI) handleReceiptStatus(w http.ResponseWriter, q url.Values) {
	txHash := common.HexToHash(q.Get("txhash"))
	if txHash != common.HexToHash("0x513c1ba0bebf66436b5fed86ab668452b7805593c05073eb2d51d3a52f480a76") {
		http.Error(w, "unknown address", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(getReceiptResponse)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
