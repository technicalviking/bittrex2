package socketPayloads

type operationType int

//Operation types for listings inside the Exchange payload
const (
	Add = operationType(iota)
	Remove
	Update
)

//ExchangeOrder values describing an order in the order book.
type ExchangeOrder struct {
	Type     operationType
	Rate     decimal
	Quantity decimal
}

//ExchangeFill values describing fill operations for an order book.
type ExchangeFill struct {
	FillID    int
	OrderType string
	Rate      decimal
	Quantity  decimal
	TimeStamp date
}

//ExchangeDelta payload for the "SubscribeToExchangeDeltas" response
type ExchangeDelta struct {
	MarketName string
	Nonce      int
	Buys       []ExchangeOrder
	Sells      []ExchangeOrder
	Fills      []ExchangeFill
}
