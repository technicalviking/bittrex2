package bittrex

import (
	"fmt"

	"github.com/technicalviking/bittrex2/socketPayloads"
)

//SubscribeToMarketSummary retrieve a filtered list of market summary deltas by market name.
func (c *Client) SubscribeToMarketSummary(market string) (chan socketPayloads.Summary, error) {
	if !c.isSubbedToSummaryDelta {
		if _, callErr := c.socketClient.CallHub(websocketHub, "SubscribeToSummaryDeltas"); callErr != nil {
			return nil, callErr
		}

		c.isSubbedToSummaryDelta = true
	}

	if ch := c.getSummaryDeltaChan(market); ch != nil {
		return ch, nil
	}

	newChan := make(chan socketPayloads.Summary)
	c.setSummaryDeltaChan(market, newChan)
	return newChan, nil
}

func (c *Client) getSummaryDeltaChan(market string) chan socketPayloads.Summary {
	c.summaryDeltaMutex.RLock()
	defer c.summaryDeltaMutex.RUnlock()

	if ch, ok := c.summaryDeltaSubscriptions[market]; ok {
		return ch
	}

	return nil
}

func (c *Client) setSummaryDeltaChan(market string, newChan chan socketPayloads.Summary) {
	c.summaryDeltaMutex.Lock()
	c.summaryDeltaSubscriptions[market] = newChan
	c.summaryDeltaMutex.Unlock()
}

//SubscribeToMarketSummaryLite retrieve a filtered list of market summary deltas (lite) by market name.
func (c *Client) SubscribeToMarketSummaryLite(market string) (chan socketPayloads.SummaryLiteDelta, error) {

	if !c.isSubbedToSummaryLiteDelta {
		if _, callErr := c.socketClient.CallHub(websocketHub, "SubscribeToSummaryLiteDeltas"); callErr != nil {
			return nil, callErr
		}

		c.isSubbedToSummaryLiteDelta = true
	}

	if ch := c.getSummaryLiteDeltaChan(market); ch != nil {
		return ch, nil
	}

	newChan := make(chan socketPayloads.SummaryLiteDelta)

	c.setSummaryLiteDeltaChan(market, newChan)
	return newChan, nil
}

func (c *Client) getSummaryLiteDeltaChan(market string) chan socketPayloads.SummaryLiteDelta {
	c.summaryDeltaMutex.RLock()
	defer c.summaryDeltaMutex.RUnlock()

	if ch, ok := c.summaryLiteDeltaSubscriptions[market]; ok {
		return ch
	}

	return nil
}

func (c *Client) setSummaryLiteDeltaChan(market string, newChan chan socketPayloads.SummaryLiteDelta) {
	c.summaryDeltaMutex.Lock()
	c.summaryLiteDeltaSubscriptions[market] = newChan
	c.summaryDeltaMutex.Unlock()
}

//SubscribeToExchange retrieve a filtered list of exchange deltas by market name.
func (c *Client) SubscribeToExchange(market string) (chan socketPayloads.ExchangeDelta, error) {
	if ch := c.getExchangeDeltaChan(market); ch != nil {
		return ch, nil
	}

	resp, callErr := c.socketClient.CallHub(websocketHub, "SubscribeToExchangeDeltas", market)
	if callErr != nil {
		return nil, callErr
	}

	if string(resp) != "true" {
		return nil, fmt.Errorf("unsuccessful subscription to %s", market)
	}

	newChan := make(chan socketPayloads.ExchangeDelta)
	c.setExchangeDeltaChan(market, newChan)
	return newChan, nil
}

func (c *Client) getExchangeDeltaChan(market string) chan socketPayloads.ExchangeDelta {
	c.exchangeDeltaMutex.RLock()
	defer c.exchangeDeltaMutex.RUnlock()

	if ch, ok := c.exchangeDeltaSubscriptions[market]; ok {
		return ch
	}

	return nil
}

func (c *Client) setExchangeDeltaChan(market string, newChan chan socketPayloads.ExchangeDelta) {
	c.exchangeDeltaMutex.Lock()
	c.exchangeDeltaSubscriptions[market] = newChan
	c.exchangeDeltaMutex.Unlock()
}
