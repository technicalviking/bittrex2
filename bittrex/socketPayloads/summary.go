package socketPayloads

//SummaryDelta Element of the "Deltas" array within the Summary struct
type SummaryDelta struct {
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

//Summary response payload for "SubscribeToSummaryDeltas"
type Summary struct {
	Nonce  int            `json:"N"`
	Deltas []SummaryDelta `json:"D"`
}
