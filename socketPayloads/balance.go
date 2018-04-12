package socketPayloads

//Delta balance for a given currency
type Delta struct {
	UUID          guid    `json:"U"`
	AccountID     int     `json:"W"`
	Currency      string  `json:"c"`
	Balance       decimal `json:"b"`
	Available     decimal `json:"a"`
	Pending       decimal `json:"z"`
	CryptoAddress string  `json:"p"`
	Requested     bool    `json:"r"`
	Updated       date    `json:"u"`
	AutoSell      bool    `json:"h"`
}

//Balance response body for balance delta events (uB)
type Balance struct {
	Nonce int `json:"N"`
	Delta `json:"d"`
}
