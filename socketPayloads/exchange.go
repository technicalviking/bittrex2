package socketPayloads

type operationType int

//Operation types for listings inside the Exchange payload
const (
	Add = operationType(iota)
	Remove
	Update
	Cancel
)

//ExchangeOrder values describing an order in the order book.
type ExchangeOrder struct {
	Type     operationType `json:"TY"`
	Rate     decimal       `json:"R"`
	Quantity decimal       `json:"Q"`
}

//ExchangeFill values describing fill operations for an order book.
type ExchangeFill struct {
	FillID    int     `json:"FI"`
	OrderType string  `json:"OT"`
	Rate      decimal `json:"R"`
	Quantity  decimal `json:"Q"`
	TimeStamp date    `json:"T"`
}

//ExchangeDelta payload for the "SubscribeToExchangeDeltas" response
type ExchangeDelta struct {
	MarketName string          `json:"M"`
	Nonce      int             `json:"N"`
	Buys       []ExchangeOrder `json:"Z"`
	Sells      []ExchangeOrder `json:"S"`
	Fills      []ExchangeFill  `json:"f"`
}

//////////////////////////////////////////////

//ExchangeStateOrder Buy or Sell struct within ExchangeState
type ExchangeStateOrder struct {
	Quantity decimal `json:"Q"`
	Rate     decimal `json:"R"`
}

//ExchangeStateFill Describes a filled order (?) within the ExchangeState
type ExchangeStateFill struct {
	ID        int     `json:"I"`
	TimeStamp date    `json:"T"`
	Quantity  decimal `json:"Q"`
	Price     decimal `json:"P"`
	Total     decimal `json:"t"`
	FillType  string  `json:"F"`
	OrderType string  `json:"OT"`
}

//ExchangeState response payload for use with "QueryExchangeState"
type ExchangeState struct {
	MarketName string               `json:"M"`
	Nonce      int                  `json:"N"`
	Buys       []ExchangeStateOrder `json:"Z"`
	Sells      []ExchangeStateOrder `json:"S"`
	Fills      []ExchangeStateFill  `json:"f"`
}
