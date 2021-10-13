package etherscan_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ryanc414/purehttp"
)

func TestMain(m *testing.M) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	os.Exit(m.Run())
}

type mockServer struct {
	apiKey      string
	checkModule bool
	module      string
	testDir     string
	responseDir string
	srv         *httptest.Server
}

func newMockServer(module string, checkModule bool) mockServer {
	m := mockServer{
		apiKey:      uuid.NewString(),
		checkModule: checkModule,
		module:      module,
		testDir:     path.Join("testData", module),
	}

	h := purehttp.NewHandler(m.handleRequest)
	m.srv = httptest.NewServer(h)

	return m
}

func (m *mockServer) close() {
	m.srv.Close()
}

func (m *mockServer) handleRequest(req *http.Request) (*purehttp.Response, error) {
	if req.URL.Path != "/api" {
		return &purehttp.Response{
			Body:       []byte("path not found\n"),
			StatusCode: http.StatusNotFound,
		}, nil
	}

	q := req.URL.Query()

	if m.checkModule {
		module := q.Get("module")
		if module != m.module {
			return &purehttp.Response{
				Body:       []byte(fmt.Sprintf("unknown module %s\n", module)),
				StatusCode: http.StatusNotFound,
			}, nil
		}
	}

	if q.Get("apikey") != m.apiKey {
		return &purehttp.Response{
			Body:       []byte("unknown API key"),
			StatusCode: http.StatusForbidden,
		}, nil
	}

	return m.handleAction(q)
}

func (m *mockServer) handleAction(q url.Values) (*purehttp.Response, error) {
	action := q.Get("action")
	if action == "" {
		return &purehttp.Response{
			Body:       []byte("action not specified"),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	params := m.filterQuery(q)

	responsePath := path.Join(m.testDir, fmt.Sprintf("%s.json", action))
	data, err := ioutil.ReadFile(responsePath)
	if err != nil {
		return nil, err
	}

	responses := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &responses); err != nil {
		return nil, err
	}

	encoded := params.Encode()
	rspData, ok := responses[encoded]
	if !ok {
		return &purehttp.Response{
			Body:       []byte(fmt.Sprintf("query path not found: %s", encoded)),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	return &purehttp.Response{
		Body:       rspData,
		StatusCode: http.StatusOK,
	}, nil
}

func (m *mockServer) filterQuery(q url.Values) url.Values {
	params := make(url.Values, len(q))

	for k, v := range q {
		if k == "apikey" || k == "action" {
			continue
		}

		if m.checkModule && k == "module" {
			continue
		}

		params[k] = v
	}

	return params
}
