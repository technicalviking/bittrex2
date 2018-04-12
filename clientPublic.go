package bittrex

import (
	"encoding/json"
	"fmt"
)

// PublicGetMarkets - public/getmarkets
func (c *Client) PublicGetMarkets() ([]MarketDescription, error) {

	parsedResponse, parseErr := c.sendRequest("public/getmarkets", nil)

	if parseErr != nil {
		return nil, parseErr
	}

	var response []MarketDescription

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return nil, fmt.Errorf("api error - public/getmarkets %s", err.Error())
	}

	//clean out responses with nil values.
	var cleanedResponse []MarketDescription
	defaultVal := MarketDescription{}

	for _, curVal := range response {
		if curVal != defaultVal {
			cleanedResponse = append(cleanedResponse, curVal)
		}
	}

	if len(cleanedResponse) == 0 {
		return nil, fmt.Errorf("validate response - all markets had empty values")
	}

	return cleanedResponse, nil
}

// PublicGetCurrencies - public/getcurrencies
func (c *Client) PublicGetCurrencies() ([]Currency, error) {

	parsedResponse, parseErr := c.sendRequest("public/getcurrencies", nil)

	if parseErr != nil {
		return nil, parseErr
	}

	var response []Currency

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return nil, fmt.Errorf("api error - public/getcurrencies %s", err.Error())
	}

	//clean out responses with nil values.
	var cleanedResponse []Currency
	defaultVal := Currency{}

	for _, curVal := range response {
		if curVal != defaultVal {
			cleanedResponse = append(cleanedResponse, curVal)
		}
	}

	if len(cleanedResponse) == 0 {
		return nil, fmt.Errorf("validate response - all currencies had empty values")
	}

	return cleanedResponse, nil
}

// PublicGetTicker - public/getticker
func (c *Client) PublicGetTicker(market string) (Ticker, error) {

	parsedResponse, parseErr := c.sendRequest("/public/getticker", map[string]string{"market": market})

	defaultValue := Ticker{}

	if parseErr != nil {
		return defaultValue, parseErr
	}

	var response Ticker

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return defaultValue, fmt.Errorf("api error - public/getticker %s", err.Error())
	}

	if response == defaultValue {
		return defaultValue, fmt.Errorf("validate response - ticker had no data")
	}

	return response, nil
}

// PublicGetMarketSummaries - public/getmarketsummaries
func (c *Client) PublicGetMarketSummaries() ([]MarketSummary, error) {

	parsedResponse, parseErr := c.sendRequest("public/getmarketsummaries", nil)

	if parseErr != nil {
		return nil, parseErr
	}

	var response []MarketSummary

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return nil, fmt.Errorf("api error - public/getmarketsummaries %s", err.Error())
	}

	//clean out responses with nil values.
	var cleanedResponse []MarketSummary
	defaultVal := MarketSummary{}

	for _, curVal := range response {
		if curVal != defaultVal {
			cleanedResponse = append(cleanedResponse, curVal)
		}
	}

	if len(cleanedResponse) == 0 {
		return nil, fmt.Errorf("validate response - all market summaries had empty values")
	}

	return cleanedResponse, nil
}

// PublicGetMarketSummary - public/getmarketsummary
func (c *Client) PublicGetMarketSummary(market string) (MarketSummary, error) {

	parsedResponse, parseErr := c.sendRequest("public/getmarketsummary", map[string]string{"market": market})

	if parseErr != nil {
		return MarketSummary{}, parseErr
	}

	defaultValue := MarketSummary{}

	var response []MarketSummary

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return defaultValue, fmt.Errorf("api error - public/getmarketsummary %s", err.Error())
	}

	if response[0] == defaultValue {
		return defaultValue, fmt.Errorf("validate response - market summary had no data")
	}

	return response[0], nil
}

// PublicGetOrderBook - public/getorderbook
func (c *Client) PublicGetOrderBook(market string, orderType string) (OrderBook, error) {

	parsedResponse, parseErr := c.sendRequest("/public/getorderbook", map[string]string{"market": market, "type": orderType})
	defaultValue := OrderBook{}

	if parseErr != nil {
		return defaultValue, parseErr
	}

	var response OrderBook

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return defaultValue, fmt.Errorf("api error - public/getorderbook %s", err.Error())
	}

	if (response.Buy == nil && response.Sell == nil) || (len(response.Buy) == 0 && len(response.Sell) == 0) {
		return defaultValue, fmt.Errorf("validate response - OrderBook had no data")
	}

	return response, nil
}

// PublicGetMarketHistory - public/getmarkethistory
func (c *Client) PublicGetMarketHistory(market string) ([]Trade, error) {

	parsedResponse, parseErr := c.sendRequest("public/getmarkethistory", map[string]string{"market": market})

	if parseErr != nil {
		return nil, parseErr
	}

	var response []Trade

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return nil, fmt.Errorf("api error - public/getmarkethistory %s", err.Error())
	}

	//clean out responses with nil values.
	var cleanedResponse []Trade
	defaultVal := Trade{}

	for _, curVal := range response {
		if curVal != defaultVal {
			cleanedResponse = append(cleanedResponse, curVal)
		}
	}

	if len(cleanedResponse) == 0 {
		return nil, fmt.Errorf("validate response - all markets had empty values")
	}

	return cleanedResponse, nil
}
