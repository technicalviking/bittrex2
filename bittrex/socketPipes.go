package bittrex

import (
	"encoding/json"

	"github.com/technicalviking/bittrex2/bittrex/socketPayloads"
)

func (c *Client) pipeEventOrderDelta(args json.RawMessage) {

}
func (c *Client) pipeBalanceDelta(args json.RawMessage) {

}
func (c *Client) pipeMarketDelta(args json.RawMessage) {

}
func (c *Client) pipeEventSummaryDelta(args json.RawMessage) {
	var summary socketPayloads.Summary
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

}
