package socketPayloads

//Order fields describing an individual Order within bittrex
type Order struct {
	UUID              guid    `json:"U"`
	ID                long    `json:"I"`
	OrderUUID         guid    `json:"OU"`
	Exchange          string  `json:"E"`
	OrderType         string  `json:"OT"`
	Quantity          decimal `json:"Q"`
	QuantityRemaining decimal `json:"q"`
	Limit             decimal `json:"X"`
	CommissionPaid    decimal `json:"n"`
	Price             decimal `json:"P"`
	PricePerUnit      decimal `json:"PU"`
	Opened            date    `json:"Y"`
	Closed            date    `json:"C"`
	IsOpen            bool    `json:"i"`
	CancelInitiated   bool    `json:"CI"`
	ImmediateOrCancel bool    `json:"K"`
	IsConditional     bool    `json:"k"`
	Condition         string  `json:"J"`
	ConditionTarget   decimal `json:"j"`
	Updated           date    `json:"u"`
}

//OrderResponse Payload response for Order Delta (uO)
type OrderResponse struct {
	AccountUUID guid  `json:"w"`
	Nonce       int   `json:"N"`
	Type        int   `json:"TY"`
	Order       Order `json:"o"`
}
