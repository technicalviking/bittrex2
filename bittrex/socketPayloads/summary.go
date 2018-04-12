package socketPayloads

//Summary Element of the "Deltas" array within the SummaryDeltaResponse struct or "Summaries" array within SummaryQueryResponses
type Summary struct {
	MarketName     string  `json:"M"`
	High           decimal `json:"H"`
	Low            decimal `json:"L"`
	Volume         decimal `json:"V"`
	Last           decimal `json:"l"`
	BaseVolume     decimal `json:"m"`
	TimeStamp      date    `json:"T"`
	Bid            decimal `json:"B"`
	Ask            decimal `json:"A"`
	OpenBuyOrders  int     `json:"G"`
	OpenSellOrders int     `json:"g"`
	PrevDay        decimal `json:"PD"`
	Created        date    `json:"x"`
}

//SummaryDeltaResponse response payload for "SubscribeToSummaryDeltas"
type SummaryDeltaResponse struct {
	Nonce  int       `json:"N"`
	Deltas []Summary `json:"D"`
}

//SummaryQueryResponse response payload for "QuerySummaryState"
type SummaryQueryResponse struct {
	Nonce     int       `json:"N"`
	Summaries []Summary `json:"s"`
}
