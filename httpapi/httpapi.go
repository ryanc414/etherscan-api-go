package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/ryanc414/etherscan-api-go/marshallers"
)

const rspStatusOK = "1"

// Params are construction parameters for the API client.
type Params struct {
	APIKey  string
	BaseURL *url.URL
	HTTP    *http.Client
}

type APIClient struct {
	apiURL url.URL
	http   *http.Client
}

func New(params *Params) *APIClient {
	apiURL := *params.BaseURL
	apiURL.Path = path.Join(apiURL.Path, "api")

	q := apiURL.Query()
	q.Set("apikey", params.APIKey)
	apiURL.RawQuery = q.Encode()

	httpClient := params.HTTP
	if httpClient == nil {
		httpClient = new(http.Client)
	}

	return &APIClient{
		apiURL: apiURL,
		http:   httpClient,
	}
}

type CallParams struct {
	Module  string
	Action  string
	Request interface{}
	Result  interface{}
}

func (r APIClient) Call(
	ctx context.Context, params *CallParams,
) error {
	rspData, err := r.Get(ctx, &RequestParams{
		Module: params.Module,
		Action: params.Action,
		Other:  marshallers.MarshalRequest(params.Request),
	})
	if err != nil {
		return err
	}

	return marshallers.UnmarshalResponse(rspData, params.Result)
}

type RequestParams struct {
	Module string
	Action string
	Other  map[string]string
}

type apiResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

func (r APIClient) Get(ctx context.Context, params *RequestParams) (json.RawMessage, error) {
	u := r.apiURL
	q := u.Query()
	q.Set("module", params.Module)
	q.Set("action", params.Action)

	for k, v := range params.Other {
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

	if rspBody.Status != "" && rspBody.Status != rspStatusOK {
		return nil, newResponseErr(&rspBody)
	}

	return rspBody.Result, nil
}

func (r APIClient) makeRequest(ctx context.Context, urlStr, method string) ([]byte, error) {
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
