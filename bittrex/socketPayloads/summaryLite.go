package socketPayloads

//SummaryLiteDelta Element of the "Deltas" array within the SummaryLite struct
type SummaryLiteDelta struct {
	MarketName string  `json:""`
	Last       decimal `json:""`
	BaseVolume decimal `json:""`
}

//SummaryLite response payload for "SubscribeToSummaryLiteDeltas"
type SummaryLite struct {
	Deltas []SummaryLiteDelta `json:"D"`
}
