package bittrex

import (
	"encoding/json"

	"github.com/technicalviking/bittrex2/socketPayloads"
)

//QueryExchangeState https://github.com/Bittrex/beta#queryexchangestate
func (c *Client) QueryExchangeState(market string) (*socketPayloads.ExchangeState, error) {
	resp, err := c.socketClient.CallHub(websocketHub, "QueryExchangeState", market)

	if err != nil {
		return nil, err
	}

	decoded, decodeErr := socketPayloads.Parse(resp)
	if decodeErr != nil {
		return nil, decodeErr
	}

	var state socketPayloads.ExchangeState
	parseErr := json.Unmarshal(decoded, &state)

	if parseErr != nil {
		return nil, parseErr
	}

	return &state, nil
}

//QuerySummaryState https://github.com/Bittrex/beta#querysummarystate
func (c *Client) QuerySummaryState() (*socketPayloads.SummaryQueryResponse, error) {
	resp, err := c.socketClient.CallHub(websocketHub, "QuerySummaryState")

	if err != nil {
		return nil, err
	}

	decoded, decodeErr := socketPayloads.Parse(resp)

	if decodeErr != nil {
		return nil, decodeErr
	}

	var state socketPayloads.SummaryQueryResponse
	parseErr := json.Unmarshal(decoded, &state)

	if parseErr != nil {
		return nil, parseErr
	}

	return &state, nil
}
