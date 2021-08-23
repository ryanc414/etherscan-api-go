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

func TestAccount(t *testing.T) {
	m := mockAccountsAPI{apiKey: uuid.NewString()}
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
}

type mockAccountsAPI struct {
	apiKey string
}

func (m mockAccountsAPI) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

	default:
		http.Error(w, "unknown action", http.StatusNotFound)
	}
}

const getBalanceResponse = `{
	"status":"1",
	"message":"OK",
	"result":"40891626854930000000000"
}`

func (m mockAccountsAPI) handleGetBalance(w http.ResponseWriter, q url.Values) {
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
