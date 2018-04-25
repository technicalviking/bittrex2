package bittrex

import "strconv"

//Order Types
const (
	OrderTypeLimit  = "LIMIT"
	OrderTypeMarket = "MARKET"
)

//Order Time In Effect
const (
	OrderTimeGTC = "GOOD_TIL_CANCELLED"
	OrderTimeIOC = "IMMEDIATE_OR_CANCEL"
	OrderTime    = "FILL_OR_KILL"
)

//Order Conditions
const (
	OrderConditionNone = "NONE"
	OrderConditionGT   = "GREATER_THAN"
	OrderConditionLT   = "LESS_THAN"
	OrderConditionSLF  = "STOP_LOSS_FIXED"
	OrderConditionSLP  = "STOP_LOSS_PERCENTAGE"
)

/*
KeyMarketTradeSell generate a sell order using the bittrex v2 rest api (undocumented)
values for query string come from https://github.com/ericsomdahl/python-bittrex/issues/35
timeInEffect and conditionType should use the constants found in this file.
Hardcoded to only place limit orders.
*/
func (c *Client) KeyMarketTradeSell(
	market string,
	quantity float64,
	rate float64,
	timeInEffect string,
	conditionType string,
	conditionTarget float64,
) (bool, error) {

	targetParam := "0"
	if conditionTarget != 0 {
		targetParam = strconv.FormatFloat(conditionTarget, 'f', 8, 64)
	}

	params := map[string]string{
		"useApi2":       "true",
		"marketName":    market,
		"orderType":     OrderTypeLimit,
		"quantity":      strconv.FormatFloat(quantity, 'f', 8, 64),
		"rate":          strconv.FormatFloat(rate, 'f', 8, 64),
		"timeInEffect":  timeInEffect,
		"conditionType": conditionType,
		"target":        targetParam,
	}

	//@todo haven't tested the response to this endpoint to know the format for sure,
	//so I'm ignoring it and sending the value of 'success'.
	//if you really want to know the new order id, provide an API key to the client
	//and subscribe to the orders chan.
	parseResponse, parseErr := c.sendRequest("key/market/TradeSell", params)

	if parseErr != nil {
		return false, parseErr
	}

	return parseResponse.Success, nil
}

/*
KeyMarketTradeBuy generate a buy order using the bittrex v2 rest api (undocumented)
values for query string come from https://github.com/ericsomdahl/python-bittrex/issues/35
timeInEffect and conditionType should use the constants found in this file.
Hardcoded to only place limit orders.
*/
func (c *Client) KeyMarketTradeBuy(
	market string,
	quantity float64,
	rate float64,
	timeInEffect string,
	conditionType string,
	conditionTarget float64,
) (bool, error) {

	targetParam := "0"
	if conditionTarget != 0 {
		targetParam = strconv.FormatFloat(conditionTarget, 'f', 8, 64)
	}

	params := map[string]string{
		"useApi2":       "true",
		"marketName":    market,
		"orderType":     OrderTypeLimit,
		"quantity":      strconv.FormatFloat(quantity, 'f', 8, 64),
		"rate":          strconv.FormatFloat(rate, 'f', 8, 64),
		"timeInEffect":  timeInEffect,
		"conditionType": conditionType,
		"target":        targetParam,
	}

	//@todo haven't tested the response to this endpoint to know the format for sure,
	//so I'm ignoring it and sending the value of 'success'.
	//if you really want to know the new order id, provide an API key to the client
	//and subscribe to the orders chan.
	parseResponse, parseErr := c.sendRequest("key/market/TradeBuy", params)

	if parseErr != nil {
		return false, parseErr
	}

	return parseResponse.Success, nil
}
