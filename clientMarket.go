package bittrex

import (
	"encoding/json"
)

// MarketBuyLimit - market/buylimit
func (c *Client) MarketBuyLimit(market string, quantity decimal, rate decimal) (TransactionID, error) {
	

	params := map[string]string{
		"apikey":   c.apiKey,
		"market":   market,
		"quantity": strconv.FormatFloat(quantity, 'f', 8, 64),
		"rate":     strconv.FormatFloat(rate, 'f', 8, 64),
	}

	parsedResponse, parseErr = c.sendRequest("market/buylimit", params)

	if parseErr != nil {
		return TransactionID{}, parseErr
	}

	if parsedResponse.Success != true {
		return TransactionID{}, fmt.Errorf("api error - market/buylimit %s", parsedResponse.Message)
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
		"quantity": quantity.String(),
		"rate":     rate.String(),
	}

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("market/selllimit", params)

	if c.err != nil {
		return TransactionID{}, c.err
	}

	if parsedResponse.Success != true {
		fmt.Errorf("api error - market/selllimit", parsedResponse.Message)
		return TransactionID{}, c.err
	}

	var response TransactionID

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		fmt.Errorf("api error - market/selllimit", err.Error())
		return TransactionID{}, c.err
	}

	return response, nil
}

// MarketBuyMarket - market/buymarket - EXPERIMENTAL
func (c *Client) MarketBuyMarket(market string, quantity decimal, rate decimal) (TransactionID, error) {
	

	params := map[string]string{
		"apikey":   c.apiKey,
		"market":   market,
		"quantity": quantity.String(),
		"rate":     rate.String(),
	}

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("market/buymarket", params)

	if c.err != nil {
		return TransactionID{}, c.err
	}

	defaultValue := TransactionID{}

	if parsedResponse.Success != true {
		fmt.Errorf("api error - /market/buymarket", parsedResponse.Message)
		return defaultValue, c.err
	}

	var response TransactionID

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		fmt.Errorf("api error - market/buymarket", err.Error())
		return defaultValue, c.err
	}

	if response == defaultValue {
		fmt.Errorf("validate response - buy limit response had no data")
		return defaultValue, c.err
	}

	return response, nil
}

// MarketSellMarket - market/sellmarket - EXPERIMENTAL
func (c *Client) MarketSellMarket(market string, quantity decimal, rate decimal) (TransactionID, error) {
	

	params := map[string]string{
		"apikey":   c.apiKey,
		"market":   market,
		"quantity": quantity.String(),
		"rate":     rate.String(),
	}

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("market/sellmarket", params)

	if c.err != nil {
		return TransactionID{}, c.err
	}

	defaultValue := TransactionID{}

	if parsedResponse.Success != true {
		fmt.Errorf("api error - /market/selllimit", parsedResponse.Message)
		return defaultValue, c.err
	}

	var response TransactionID

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		fmt.Errorf("api error - market/sellmarket", err.Error())
		return defaultValue, c.err
	}

	if response == defaultValue {
		fmt.Errorf("validate response - sell limit response had no data")
		return defaultValue, c.err
	}

	return response, nil
}

// MarketCancel - market/cancel
func (c *Client) MarketCancel(uuid string) (bool, error) {
	

	params := map[string]string{
		"apikey": c.apiKey,
		"uuid":   uuid,
	}

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("market/cancel", params)

	if c.err != nil {
		return false, c.err
	}

	if parsedResponse.Success != true {
		fmt.Errorf("api error - market/cancel", parsedResponse.Message)
		return false, c.err
	}

	return true, nil
}

// MarketGetOpenOrders - market/getopenorders
func (c *Client) MarketGetOpenOrders(market string) ([]OrderDescription, error) {
	

	params := map[string]string{
		"market": market,
		"apikey": c.apiKey,
	}

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("market/getopenorders", params)

	if c.err != nil {
		return nil, c.err
	}

	if parsedResponse.Success != true {
		fmt.Errorf("api error - market/getopenorders", parsedResponse.Message)
		return nil, c.err
	}

	var response []OrderDescription

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		fmt.Errorf("api error - market/getopenorders", err.Error())
		return nil, c.err
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
		fmt.Errorf("validate response - all historical deposits had empty values")
		return nil, c.err
	}

	return cleanedResponse, nil
}
