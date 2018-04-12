package bittrex

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// MarketBuyLimit - market/buylimit
func (c *Client) MarketBuyLimit(market string, quantity decimal, rate decimal) (TransactionID, error) {

	params := map[string]string{
		"apikey":   c.apiKey,
		"market":   market,
		"quantity": strconv.FormatFloat(quantity, 'f', 8, 64),
		"rate":     strconv.FormatFloat(rate, 'f', 8, 64),
	}

	parsedResponse, parseErr := c.sendRequest("market/buylimit", params)

	if parseErr != nil {
		return TransactionID{}, parseErr
	}

	var response TransactionID

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return TransactionID{}, fmt.Errorf("api error - market/buylimit: %s", err.Error())
	}

	return response, nil
}

// MarketSellLimit - market/selllimit
func (c *Client) MarketSellLimit(market string, quantity decimal, rate decimal) (TransactionID, error) {

	params := map[string]string{
		"apikey":   c.apiKey,
		"market":   market,
		"quantity": strconv.FormatFloat(quantity, 'f', 8, 64),
		"rate":     strconv.FormatFloat(rate, 'f', 8, 64),
	}

	parsedResponse, parseErr := c.sendRequest("market/selllimit", params)

	if parseErr != nil {
		return TransactionID{}, parseErr
	}

	var response TransactionID

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return TransactionID{}, fmt.Errorf("api error - market/selllimit %s", err.Error())
	}

	return response, nil
}

// MarketCancel - market/cancel
func (c *Client) MarketCancel(uuid string) (bool, error) {

	params := map[string]string{
		"apikey": c.apiKey,
		"uuid":   uuid,
	}

	_, parseErr := c.sendRequest("market/cancel", params)

	if parseErr != nil {
		return false, parseErr
	}

	return true, nil
}

// MarketGetOpenOrders - market/getopenorders
func (c *Client) MarketGetOpenOrders(market string) ([]OrderDescription, error) {

	params := map[string]string{
		"market": market,
		"apikey": c.apiKey,
	}

	parsedResponse, parseErr := c.sendRequest("market/getopenorders", params)

	if parseErr != nil {
		return nil, parseErr
	}

	var response []OrderDescription

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return nil, fmt.Errorf("api error - market/getopenorders %s", err.Error())
	}

	//clean out responses with nil values.
	var cleanedResponse []OrderDescription
	defaultVal := OrderDescription{}

	for _, curVal := range response {
		if curVal != defaultVal {
			cleanedResponse = append(cleanedResponse, curVal)
		}
	}

	if len(cleanedResponse) == 0 && len(response) != 0 {
		return nil, fmt.Errorf("validate response - all historical deposits had empty values")
	}

	return cleanedResponse, nil
}
