package bittrex

import (
	"encoding/json"

	"github.com/technicalviking/bittrex2/socketPayloads"
)

func (c *Client) pipeEventOrderDelta(args json.RawMessage) {
	var order socketPayloads.OrderResponse
	parseErr := json.Unmarshal(args, &order)

	if parseErr != nil {
		//@todo pipe this to the client's error chan, when it exists.
		panic(parseErr)
	}

	c.orderSubscription <- order
}

func (c *Client) pipeBalanceDelta(args json.RawMessage) {
	var balance socketPayloads.Balance
	parseErr := json.Unmarshal(args, &balance)

	if parseErr != nil {
		//@todo pipe this to the client's error chan, when it exists.
		panic(parseErr)
	}

	c.balanceSubscription <- balance.BalanceDelta
}

func (c *Client) pipeMarketExchangeDelta(args json.RawMessage) {
	var exchangeDelta socketPayloads.ExchangeDelta
	parseErr := json.Unmarshal(args, &exchangeDelta)

	if parseErr != nil {
		//@todo pipe this to the client's error chan, when it exists.
		panic(parseErr)
	}

	c.exchangeDeltaMutex.Lock()
	defer c.exchangeDeltaMutex.Unlock()

	if _, ok := c.exchangeDeltaSubscriptions[exchangeDelta.MarketName]; ok {
		c.exchangeDeltaSubscriptions[exchangeDelta.MarketName] <- exchangeDelta
	}
}

func (c *Client) pipeEventSummaryDelta(args json.RawMessage) {
	var summary socketPayloads.SummaryDeltaResponse
	parseErr := json.Unmarshal(args, &summary)

	if parseErr != nil {
		//@todo pipe this to the client's error chan, when it exists.
		panic(parseErr)
	}

	c.summaryDeltaMutex.Lock()
	defer c.summaryDeltaMutex.Unlock()

	for _, curDelta := range summary.Deltas {
		if _, ok := c.summaryDeltaSubscriptions[curDelta.MarketName]; ok {
			c.summaryDeltaSubscriptions[curDelta.MarketName] <- curDelta
		}
	}
}

func (c *Client) pipeEventSummaryDeltaLite(args json.RawMessage) {
	var summary socketPayloads.SummaryLite
	parseErr := json.Unmarshal(args, &summary)

	if parseErr != nil {
		//@todo pipe this to the client's error chan, when it exists.
		panic(parseErr)
	}

	c.summaryLiteDeltaMutex.Lock()
	defer c.summaryLiteDeltaMutex.Unlock()

	for _, curDelta := range summary.Deltas {
		if _, ok := c.summaryLiteDeltaSubscriptions[curDelta.MarketName]; ok {
			c.summaryLiteDeltaSubscriptions[curDelta.MarketName] <- curDelta
		}
	}
}
