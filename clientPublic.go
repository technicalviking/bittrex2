package bittrex

import "encoding/json"

// PublicGetMarkets - public/getmarkets
func (c *Client) PublicGetMarkets() ([]MarketDescription, error) {
	

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("public/getmarkets", nil)

	if c.err != nil {
		return nil, c.err
	}

	if parsedResponse.Success != true {
		fmt.Errorf("api error - public/getmarkets", parsedResponse.Message)
		return nil, c.err
	}

	var response []MarketDescription

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		fmt.Errorf("api error - public/getmarkets", err.Error())
		return nil, c.err
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
		fmt.Errorf("validate response - all markets had empty values.")
		return nil, c.err
	}

	return cleanedResponse, nil
}

// PublicGetCurrencies - public/getcurrencies
func (c *Client) PublicGetCurrencies() ([]Currency, error) {
	

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("public/getcurrencies", nil)

	if c.err != nil {
		return nil, c.err
	}

	if parsedResponse.Success != true {
		fmt.Errorf("api error - public/getcurrencies", parsedResponse.Message)
		return nil, c.err
	}

	var response []Currency

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		fmt.Errorf("api error - public/getcurrencies", err.Error())
		return nil, c.err
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
		fmt.Errorf("validate response - all markets had empty values.")
		return nil, c.err
	}

	return cleanedResponse, nil
}

// PublicGetTicker - public/getticker
func (c *Client) PublicGetTicker(market string) (Ticker, error) {
	

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("/public/getticker", map[string]string{"market": market})
	defaultValue := Ticker{}

	if parsedResponse.Success != true {
		fmt.Errorf("api error - /public/getticker", parsedResponse.Message)
		return defaultValue, c.err
	}

	var response Ticker

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		fmt.Errorf("api error - public/getticker", err.Error())
		return defaultValue, c.err
	}

	if response == defaultValue {
		fmt.Errorf("validate response - ticker had no data.")
		return defaultValue, c.err
	}

	return response, nil
}

// PublicGetMarketSummaries - public/getmarketsummaries
func (c *Client) PublicGetMarketSummaries() ([]MarketSummary, error) {
	

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("public/getmarketsummaries", nil)

	if c.err != nil {
		return nil, c.err
	}

	if parsedResponse.Success != true {
		fmt.Errorf("api error - public/getmarketsummaries", parsedResponse.Message)
		return nil, c.err
	}

	var response []MarketSummary

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		fmt.Errorf("api error - public/getmarketsummaries", err.Error())
		return nil, c.err
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
		fmt.Errorf("validate response - all markets had empty values.")
		return nil, c.err
	}

	return cleanedResponse, nil
}

// PublicGetMarketSummary - public/getmarketsummary
func (c *Client) PublicGetMarketSummary(market string) (MarketSummary, error) {
	

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("public/getmarketsummary", map[string]string{"market": market})

	if c.err != nil {
		return MarketSummary{}, c.err
	}

	defaultValue := MarketSummary{}

	if parsedResponse.Success != true {
		fmt.Errorf("api error - /public/getmarketsummary", parsedResponse.Message)
		return defaultValue, c.err
	}

	var response []MarketSummary

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		fmt.Errorf("api error - public/getmarketsummary", err.Error())
		return defaultValue, c.err
	}

	if response[0] == defaultValue {
		fmt.Errorf("validate response - marketsummary had no data.")
		return defaultValue, c.err
	}

	return response[0], nil
}

// PublicGetOrderBook - public/getorderbook
func (c *Client) PublicGetOrderBook(market string, orderType string) (OrderBook, error) {
	

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("/public/getorderbook", map[string]string{"market": market, "type": orderType})
	defaultValue := OrderBook{}
	if parsedResponse.Success != true {
		fmt.Errorf("api error - /public/getorderbook", parsedResponse.Message)
		return defaultValue, c.err
	}

	var response OrderBook

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		fmt.Errorf("api error - public/getorderbook", err.Error())
		return defaultValue, c.err
	}

	if (response.Buy == nil && response.Sell == nil) || (len(response.Buy) == 0 && len(response.Sell) == 0) {
		fmt.Errorf("validate response - OrderBook had no data.")
		return defaultValue, c.err
	}

	return response, nil
}

// PublicGetMarketHistory - public/getmarkethistory
func (c *Client) PublicGetMarketHistory(market string) ([]Trade, error) {
	

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("public/getmarkethistory", map[string]string{"market": market})

	if c.err != nil {
		return nil, c.err
	}

	if parsedResponse.Success != true {
		fmt.Errorf("api error - public/getmarkethistory", parsedResponse.Message)
		return nil, c.err
	}

	var response []Trade

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		fmt.Errorf("api error - public/getmarkethistory", err.Error())
		return nil, c.err
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
		fmt.Errorf("validate response - all markets had empty values.")
		return nil, c.err
	}

	return cleanedResponse, nil
}
