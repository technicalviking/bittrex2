package bittrex

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// AccountGetBalances - /account/getbalances
func (c *Client) AccountGetBalances() ([]AccountBalance, error) {

	params := map[string]string{
		"apikey": c.apiKey,
	}

	//var parsedResponse *baseResponse

	parsedResponse, parseErr := c.sendRequest("account/getbalances", params)

	if parseErr != nil {
		return nil, parseErr
	}

	var response []AccountBalance

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return nil, fmt.Errorf("api error - account/getbalances %s", err.Error())
	}

	//clean out responses with nil values.
	var cleanedResponse []AccountBalance
	defaultAB := AccountBalance{}

	for _, curBalance := range response {
		if curBalance != defaultAB {
			cleanedResponse = append(cleanedResponse, curBalance)
		}
	}

	if len(cleanedResponse) == 0 && len(response) != 0 {

		return nil, fmt.Errorf("validate response - all account balances had empty values")
	}

	return cleanedResponse, nil
}

// AccountGetBalance - /account/getbalance
func (c *Client) AccountGetBalance(currency string) (AccountBalance, error) {

	params := map[string]string{
		"apikey":   c.apiKey,
		"currency": currency,
	}

	parsedResponse, parseErr := c.sendRequest("account/getbalance", params)

	if parseErr != nil {
		return AccountBalance{}, parseErr
	}

	var response AccountBalance

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {

		return AccountBalance{}, fmt.Errorf("api error - account/getbalance %s", err.Error())
	}

	if response == (AccountBalance{}) {
		return AccountBalance{}, fmt.Errorf("validate response - account balance had empty values")
	}

	return response, nil
}

// AccountGetDepositAddress - /account/getdepositaddress
func (c *Client) AccountGetDepositAddress(currency string) (WalletAddress, error) {

	params := map[string]string{
		"apikey":   c.apiKey,
		"currency": currency,
	}

	parsedResponse, parseErr := c.sendRequest("account/getdepositaddress", params)

	if parseErr != nil {
		return WalletAddress{}, parseErr
	}

	var response WalletAddress
	defaultVal := WalletAddress{}

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return defaultVal, fmt.Errorf("api error - account/getdepositaddress %s", err.Error())
	}

	if response == defaultVal {
		return defaultVal, fmt.Errorf("validate response - deposit address empty")
	}

	return response, nil
}

/*
AccountWithdraw - /account/withdraw
paymentId field is optional for the api (used as a memo field for other services
such as CryptoNotes, BitShareX, Nxt).  Set it to empty string to exclude it from
api call
*/
func (c *Client) AccountWithdraw(currency string, quantity decimal, address string, paymentID string) (TransactionID, error) {

	params := map[string]string{
		"apikey":   c.apiKey,
		"currency": currency,
		"quantity": strconv.FormatFloat(quantity, 'f', 8, 64),
		"address":  address,
	}

	if paymentID != "" {
		params["paymentid"] = paymentID
	}

	parsedResponse, parseErr := c.sendRequest("account/withdraw", params)

	if parseErr != nil {
		return TransactionID{}, parseErr
	}

	var response TransactionID
	defaultVal := TransactionID{}

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {

		return defaultVal, fmt.Errorf("api error - account/withdraw %s", err.Error())
	}

	if response == defaultVal {
		fmt.Errorf("validate response nil vals in withdraw response")
	}

	return response, nil
}

// AccountGetOrder - /account/getorder
func (c *Client) AccountGetOrder(orderID string) (AccountOrderDescription, error) {

	params := map[string]string{
		"apikey": c.apiKey,
		"uuid":   orderID,
	}

	parsedResponse, parseErr := c.sendRequest("account/getorder", params)

	if parseErr != nil {
		return AccountOrderDescription{}, parseErr
	}

	defaultVal := AccountOrderDescription{}

	var response AccountOrderDescription

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return defaultVal, fmt.Errorf("api error - account/getorder %s", err.Error())
	}

	if response == defaultVal {
		return defaultVal, fmt.Errorf("validate response - nil vals in get order response")
	}

	return response, nil
}

/*
AccountGetOrderHistory - /account/getorderhistory
market is optional param.  set it to empty strinng to get all markets.
*/
func (c *Client) AccountGetOrderHistory(market string) ([]AccountOrderHistoryDescription, error) {

	params := map[string]string{
		"apikey": c.apiKey,
	}

	if market != "" {
		params["market"] = market
	}

	parsedResponse, parseErr := c.sendRequest("account/getorderhistory", params)

	if parseErr != nil {
		return nil, parseErr
	}

	var response []AccountOrderHistoryDescription

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return nil, fmt.Errorf("api error - account/getorderhistory %s", err.Error())
	}

	//clean out responses with nil values.
	var cleanedResponse []AccountOrderHistoryDescription
	defaultVal := AccountOrderHistoryDescription{}

	for _, curVal := range response {
		if curVal != defaultVal {
			cleanedResponse = append(cleanedResponse, curVal)
		}
	}

	if len(cleanedResponse) == 0 && len(response) != 0 {
		return nil, fmt.Errorf("validate response - all historical orders had empty values")
	}

	return cleanedResponse, nil
}

/*
AccountGetWithdrawalHistory - /account/getwithdrawalhistory
setting currency to empty string will get all currencies.
*/
func (c *Client) AccountGetWithdrawalHistory(currency string) ([]TransactionHistoryDescription, error) {

	params := map[string]string{
		"apikey": c.apiKey,
	}

	if currency != "" {
		params["currency"] = currency
	}

	parsedResponse, parseErr := c.sendRequest("account/getwithdrawalhistory", params)

	if parseErr != nil {
		return nil, parseErr
	}

	if parsedResponse.Success != true {
		return nil, fmt.Errorf("api error - account/getwithdrawalhistory %s", parsedResponse.Message)
	}

	var response []TransactionHistoryDescription

	//clean out responses with nil values.
	var cleanedResponse []TransactionHistoryDescription
	defaultVal := TransactionHistoryDescription{}

	for _, curVal := range response {
		if curVal != defaultVal {
			cleanedResponse = append(cleanedResponse, curVal)
		}
	}

	if len(cleanedResponse) == 0 && len(response) != 0 {
		return nil, fmt.Errorf("validate response - all historical withdrawals had empty values")
	}

	return cleanedResponse, nil
}

/*
AccountGetDepositHistory - /account/getdeposithistory
setting currency to empty string will get all currencies.
*/
func (c *Client) AccountGetDepositHistory(currency string) ([]TransactionHistoryDescription, error) {

	params := map[string]string{
		"apikey": c.apiKey,
	}

	if currency != "" {
		params["currency"] = currency
	}

	parsedResponse, parseErr := c.sendRequest("account/getdeposithistory", params)

	if parseErr != nil {
		return nil, parseErr
	}

	var response []TransactionHistoryDescription

	if err := json.Unmarshal(parsedResponse.Result, &response); err != nil {
		return nil, fmt.Errorf("api error - account/getdeposithistory %s", err.Error())
	}

	//clean out responses with nil values.
	var cleanedResponse []TransactionHistoryDescription
	defaultVal := TransactionHistoryDescription{}

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
