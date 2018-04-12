package bittrex

const (
	baseURL          string = "https://bittrex.com"
	v1APIURL         string = baseURL + "/v1.1"
	v2APIURL         string = baseURL + "/v2.0"
	websocketBaseURI string = "https://beta.bittrex.com"
	websocketHub     string = "c2" //SignalR main hub
	defaultTimeout   int64  = 30
	//signalR events
	eventOrderDelta       string = "uO"
	eventBalanceDelta     string = "uB"
	eventMarketDelta      string = "uE"
	eventSummaryDelta     string = "uS"
	eventSummaryDeltaLite string = "uL"
)
