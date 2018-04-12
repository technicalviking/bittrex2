package bittrex

import "encoding/json"

const (
	//TickIntervalOneMin oneMin = 10 days worth of candles
	TickIntervalOneMin = "oneMin"

	//TickIntervalFiveMin fiveMin = 20 days worth of candles
	TickIntervalFiveMin = "fiveMin"

	//TickIntervalThirtyMin thirtyMin = 40 days worth of candles
	TickIntervalThirtyMin = "thirtyMin"

	//TickIntervalHour hour = 60 days worth of candles
	TickIntervalHour = "hour"

	//TickIntervalDay day = 1385 days (nearly four years)
	TickIntervalDay = "day"
)

// PubMarketGetTicks - /pub/market/getticks
// interval must be one of the TickInterval consts
func (c *Client) PubMarketGetTicks(market string, interval string) ([]Candle, error) {
	

	params := map[string]string{
		"marketName":   market,
		"tickInterval": interval,
		"useApi2":      "true",
	}

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("pub/market/getticks", params)

	if c.err != nil {
		return nil, c.err
	}

	if parsedResponse.Success != true {
		fmt.Errorf("api error - pub/market/getticks", parsedResponse.Message)
		return nil, c.err
	}

	var response []Candle

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		fmt.Errorf("api error - pub/market/getticks", err.Error())
		return nil, c.err
	}

	//clean out responses with nil values.
	var cleanedResponse []Candle
	defaultVal := Candle{}

	for _, curVal := range response {
		if curVal != defaultVal {
			cleanedResponse = append(cleanedResponse, curVal)
		}
	}

	if len(cleanedResponse) == 0 {
		fmt.Errorf("validate response - all candles had empty values")
		return nil, c.err
	}

	return cleanedResponse, nil
}

// PubMarketGetLatestTick - /pub/market/getticks
// interval must be one of the TickInterval consts
func (c *Client) PubMarketGetLatestTick(market string, interval string) (Candle, error) {
	

	params := map[string]string{
		"marketName":   market,
		"tickInterval": interval,
		"useApi2":      "true",
	}

	var parsedResponse *baseResponse

	parsedResponse = c.sendRequest("pub/market/getlatesttick", params)

	if c.err != nil {
		return Candle{}, c.err
	}

	if parsedResponse.Success != true {
		fmt.Errorf("api error - pub/market/getlatesttick", parsedResponse.Message)
		return Candle{}, c.err
	}

	var response []Candle

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		fmt.Errorf("api error - pub/market/getlatesttick", err.Error())
		return Candle{}, c.err
	}

	return response[0], nil
}
