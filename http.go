package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/rs/zerolog/log"
)

const rspStatusOK = "1"

type apiClient struct {
	apiURL url.URL
	http   *http.Client
}

func newAPIClient(params *Params) *apiClient {
	apiURL := *params.BaseURL
	apiURL.Path = path.Join(apiURL.Path, "api")

	q := apiURL.Query()
	q.Set("apikey", params.APIKey)
	apiURL.RawQuery = q.Encode()

	httpClient := params.HTTP
	if httpClient == nil {
		httpClient = new(http.Client)
	}

	return &apiClient{
		apiURL: apiURL,
		http:   httpClient,
	}
}

type requestParams struct {
	module string
	action string
	other  map[string]string
}

type apiResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

func (r apiClient) get(ctx context.Context, params *requestParams) (json.RawMessage, error) {
	u := r.apiURL
	q := u.Query()
	q.Set("module", params.module)
	q.Set("action", params.action)

	for k, v := range params.other {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()

	bodyData, err := r.makeRequest(ctx, u.String(), http.MethodGet)
	if err != nil {
		return nil, err
	}

	var rspBody apiResponse
	if err := json.Unmarshal(bodyData, &rspBody); err != nil {
		return nil, err
	}

	if rspBody.Status != rspStatusOK {
		return nil, newResponseErr(&rspBody)
	}

	return rspBody.Result, nil
}

func (r apiClient) makeRequest(ctx context.Context, urlStr, method string) ([]byte, error) {
	log.Debug().Str("url", urlStr).Str("method", method).Msg("making HTTP request")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	rsp, err := r.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, newHTTPErr(rsp)
	}

	bodyData, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return bodyData, nil
}

type httpError struct {
	status string
	body   []byte
}

func (err *httpError) Error() string {
	if len(err.body) == 0 {
		return err.status
	}

	return fmt.Sprintf("%s %s", err.status, string(err.body))
}

func newHTTPErr(rsp *http.Response) *httpError {
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Error().Err(err).Msg("error while reading HTTP response body")
		body = nil
	}

	return &httpError{
		status: rsp.Status,
		body:   body,
	}
}

type responseError struct {
	rsp apiResponse
}

func (err responseError) Error() string {
	if len(err.rsp.Result) == 0 {
		return fmt.Sprintf("API error - Status: %s, Message: %s", err.rsp.Status, err.rsp.Message)
	}

	return fmt.Sprintf(
		"API error - Status: %s, Message: %s, Result: %s",
		err.rsp.Status,
		err.rsp.Message,
		string(err.rsp.Result),
	)
}

func newResponseErr(rsp *apiResponse) responseError {
	return responseError{rsp: *rsp}
}
