package bittrex

import "github.com/technicalviking/bittrex2/bittrex/socketPayloads"

//SubscribeToMarketSummary retrieve a filtered list of market summary deltas by market name.
func (c *Client) SubscribeToMarketSummary(market string) chan socketPayloads.Summary {
	c.summaryDeltaMutex.Lock()
	defer c.summaryDeltaMutex.Unlock()

	if c.summaryDeltaSubscriptions == nil {
		c.summaryDeltaSubscriptions = make(map[string]chan socketPayloads.Summary)
		c.socketClient.CallHub(websocketHub, "SubscribeToSummaryDeltas")
	}

	if _, ok := c.summaryDeltaSubscriptions[market]; ok {
		return c.summaryDeltaSubscriptions[market]
	}

	c.summaryDeltaSubscriptions[market] = make(chan socketPayloads.Summary)

	return c.summaryDeltaSubscriptions[market]
}

//SubscribeToMarketSummaryLite retrieve a filtered list of market summary deltas (lite) by market name.
func (c *Client) SubscribeToMarketSummaryLite(market string) chan socketPayloads.SummaryLiteDelta {
	c.summaryLiteDeltaMutex.Lock()
	defer c.summaryLiteDeltaMutex.Unlock()

	if c.summaryLiteDeltaSubscriptions == nil {
		c.summaryLiteDeltaSubscriptions = make(map[string]chan socketPayloads.SummaryLiteDelta)
		c.socketClient.CallHub(websocketHub, "SubscribeToSummaryLiteDeltas")
	}

	if _, ok := c.summaryLiteDeltaSubscriptions[market]; ok {
		return c.summaryLiteDeltaSubscriptions[market]
	}

	c.summaryLiteDeltaSubscriptions[market] = make(chan socketPayloads.SummaryLiteDelta)

	return c.summaryLiteDeltaSubscriptions[market]
}

//SubscribeToExchange retrieve a filtered list of exchange deltas by market name.
func (c *Client) SubscribeToExchange(market string) chan socketPayloads.ExchangeDelta {
	c.exchangeDeltaMutex.Lock()
	defer c.exchangeDeltaMutex.Unlock()

	if c.exchangeDeltaSubscriptions == nil {
		c.exchangeDeltaSubscriptions = make(map[string]chan socketPayloads.ExchangeDelta)
	}

	if _, ok := c.exchangeDeltaSubscriptions[market]; ok {
		return c.exchangeDeltaSubscriptions[market]
	}

	c.socketClient.CallHub(websocketHub, "SubscribeToExchangeDeltas", market)

	c.exchangeDeltaSubscriptions[market] = make(chan socketPayloads.ExchangeDelta)

	return c.exchangeDeltaSubscriptions[market]
}
