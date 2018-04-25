/*
	As of End of March 2018, Bittrex imposed limitations on their API to 60 calls / second,
	and started internally caching results to these calls such that AT BEST the results are only
	3 minutes old. This means that PubMarketGetLatestTick is pretty much useless.

	If you need candles, I'd recommend using the PubMarketGetTicks call once a minute until you have 'fresh' results,
	then swapping over to the socketSubscription for exchange deltas.  --DM
*/

package bittrex

import (
	"encoding/json"
	"fmt"
)

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

	parsedResponse, parseErr := c.sendRequest("pub/market/getticks", params)

	if parseErr != nil {
		return nil, parseErr
	}

	var response []Candle

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return nil, fmt.Errorf("api error - pub/market/getticks %s", err.Error())
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
		return nil, fmt.Errorf("validate response - all candles had empty values")
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

	parsedResponse, parseErr := c.sendRequest("pub/market/getlatesttick", params)

	if parseErr != nil {
		return Candle{}, parseErr
	}

	var response []Candle

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return Candle{}, fmt.Errorf("api error - pub/market/getlatesttick %s", err.Error())
	}

	return response[0], nil
}
