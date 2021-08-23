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

	t.Run("TextMultiGetETHBalance", func(t *testing.T) {
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

	case "balancemulti":
		m.handleGetMultiBalance(w, q)

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

func (m mockAccountsAPI) handleGetMultiBalance(w http.ResponseWriter, q url.Values) {
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
