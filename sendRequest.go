package bittrex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type queryParams = map[string]string

func (c *Client) sendRequest(endpoint string, params queryParams) (*baseResponse, error) {
	fullURI := c.getFullURI(endpoint, params)

	sign := c.sign(fullURI)

	var request *http.Request
	var reqErr error

	if request, reqErr = http.NewRequest("GET", fullURI, nil); reqErr != nil {
		return nil, fmt.Errorf("sendRequest - make request: %s", reqErr.Error())
	}

	request.Header.Add("apisign", sign)

	var resp *http.Response
	var respErr error

	done := make(chan error, 1)

	clientTimer := time.NewTimer(c.timeout)

	go func() {
		httpClient := &http.Client{}
		if resp, respErr = httpClient.Do(request); respErr != nil {
			done <- fmt.Errorf("sendRequest - do request: %s", respErr.Error())
		}

		done <- nil
	}()

	select {
	case e := <-done:
		if e != nil {
			return nil, e
		}
	case <-clientTimer.C:
		return nil, fmt.Errorf("sendRequest - do request %s",
			fmt.Sprintf(
				"BittrexAPI request timeout at %d seconds",
				c.timeout/time.Second,
			),
		)
	}

	defer resp.Body.Close()

	var rawBody []byte
	var readErr error

	if rawBody, readErr = ioutil.ReadAll(resp.Body); readErr != nil {
		return nil, fmt.Errorf("sendRequest - read response %s", readErr.Error())
	}

	var response baseResponse

	if rawBody == nil || len(rawBody) == 0 {
		response = baseResponse{
			Success: false,
			Message: fmt.Sprintf("Response from API endpoint %s was nil or empty", endpoint),
			Result:  rawBody,
		}
	} else if parseBaseResponseErr := json.Unmarshal(rawBody, &response); parseBaseResponseErr != nil {
		return nil, fmt.Errorf("parseBaseResponseErr for endpoint %s, %+v", endpoint, parseBaseResponseErr.Error())
	}

	if response.Success == false {
		return nil, fmt.Errorf("Send Request Endpoint - %s: %s", endpoint, response.Message)
	}

	return &response, nil
}

func (c *Client) getFullURI(endpoint string, params queryParams) string {

	apiURI := v1APIURL
	if params["useApi2"] != "" {
		apiURI = v2APIURL
		delete(params, "useApi2")
	}

	fullURI := strings.Join([]string{apiURI, endpoint}, "/")

	u, _ := url.Parse(fullURI)

	query := u.Query()

	query.Set("nonce", fmt.Sprintf("%d", time.Now().Unix()))
	query.Set("apikey", c.apiKey)

	//prevent 304 responses.
	query.Set("_", fmt.Sprintf("%d", time.Now().Unix()))

	for param, value := range params {
		query.Set(param, value)
	}

	u.RawQuery = query.Encode()

	return u.String()
}
