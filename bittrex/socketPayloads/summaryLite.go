package socketPayloads

//SummaryLiteDelta Element of the "Deltas" array within the SummaryLite struct
type SummaryLiteDelta struct {
	MarketName string  `json:"M"`
	Last       decimal `json:"l"`
	BaseVolume decimal `json:"m"`
}

//SummaryLite response payload for "SubscribeToSummaryLiteDeltas"
type SummaryLite struct {
	Deltas []SummaryLiteDelta `json:"D"`
}
