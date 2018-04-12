package bittrex

import (
	"encoding/json"
	"time"
)

type decimal = float64

//BittrexTimestamp - time.Time supertype created for parsing date values from bittrex v1.1 and v2.0 json payloads.
type BittrexTimestamp time.Time

//UnmarshalJSON implement json.Unmarshaler interface
func (bt *BittrexTimestamp) UnmarshalJSON(raw []byte) error {
	if raw == nil || len(raw) == 0 {
		bt = nil
		return nil
	}

	var strTimestamp string //"2014-07-09T07:19:30.15"

	if err := json.Unmarshal(raw, &strTimestamp); err != nil {
		return err
	}

	*bt = time.Parse("2006-01-02T03:04:05", strTimestamp)
}

//String implement stringer interface
func (bt *BittrexTimestamp) String() string {
	cast := time.Time(*bt)
	return cast.Format("2006-01-02T03:04:05")
}

//Time cast BittrexTimeStamp back to time.Time
func (bt *BittrexTimestamp) Time() time.Time {
	return time.Time(*bt)
}

//MarketDescription Result element as described under /public/getmarkets
type MarketDescription struct {
	MarketCurrency     string           `json:"MarketCurrency"`
	BaseCurrency       string           `json:"BaseCurrency"`
	MarketCurrencyLong string           `json:"MarketCurrencyLong"`
	BaseCurrencyLong   string           `json:"BaseCurrencyLong"`
	MinTradeSize       decimal          `json:"MinTradeSize"`
	MarketName         string           `json:"MarketName"`
	IsActive           bool             `json:"IsActive"`
	Created            BittrexTimestamp `json:"Created"`
}

//Currency Result element as described under /public/getcurrencies
type Currency struct {
	Currency        string  `json:"Currency"`
	CurrencyLong    string  `json:"CurrencyLong"`
	MinConfirmation int     `json:"MinConfirmation"`
	TxFee           decimal `json:"TxFee"`
	IsActive        bool    `json:"IsActive"`
	CoinType        string  `json:"CoinType"`
	BaseAddress     string  `json:"BaseAddress"`
}

//Ticker Result element as described under /public/getticker
type Ticker struct {
	Bid  decimal `json:"Bid"`
	Ask  decimal `json:"Ask"`
	Last decimal `json:"Last"`
}

//MarketSummary result element as described under /public/getmarketsummaries
type MarketSummary struct {
	MarketName        string           `json:"MarketName"`        // : "BTC-888",
	High              decimal          `json:"High"`              // : 0.00000919,
	Low               decimal          `json:"Low"`               // : 0.00000820,
	Volume            decimal          `json:"Volume"`            // : 74339.61396015,
	Last              decimal          `json:"Last"`              // : 0.00000820,
	BaseVolume        decimal          `json:"BaseVolume"`        // : 0.64966963,
	TimeStamp         BittrexTimestamp `json:"TimeStamp"`         // : "2014-07-09T07:19:30.15",
	Bid               decimal          `json:"Bid"`               // : 0.00000820,
	Ask               decimal          `json:"Ask"`               // : 0.00000831,
	OpenBuyOrders     int              `json:"OpenBuyOrders"`     // : 15,
	OpenSellOrders    int              `json:"OpenSellOrders"`    // : 15,
	PrevDay           decimal          `json:"PrevDay"`           // : 0.00000821,
	Created           BittrexTimestamp `json:"Created"`           // : "2014-03-20T06:00:00",
	DisplayMarketName string           `json:"DisplayMarketName"` // : null
}

//OrderElement element found under 'buy' or 'sell' in an OrderBook
type OrderElement struct {
	Quantity decimal `json:"Quantity"`
	Rate     decimal `json:"Rate"`
}

//OrderBook Result body of /public/getorderbook
type OrderBook struct {
	Buy  []OrderElement `json:"buy"`
	Sell []OrderElement `json:"sell"`
}

//Trade result element as described under /public/getmarkethistory
type Trade struct {
	ID        string           `json:"Id"`        // : 319435,
	TimeStamp BittrexTimestamp `json:"TimeStamp"` // : "2014-07-09T03:21:20.08",
	Quantity  decimal          `json:"Quantity"`  // : 0.30802438,
	Price     decimal          `json:"Price"`     // : 0.01263400,
	Total     decimal          `json:"Total"`     // : 0.00389158,
	FillType  string           `json:"FillType"`  // : "FILL",
	OrderType string           `json:"OrderType"` // : "BUY" or "SELL"
}

//TransactionID Result body of /market/buylimit and /market/sellimit
type TransactionID struct {
	UUID string `json:"uuid"`
}

//OrderDescription result element as described under /market/getopenorders
type OrderDescription struct {
	UUID              string           `json:"Uuid"`              // : null,
	OrderUUID         string           `json:"OrderUuid"`         // : "09aa5bb6-8232-41aa-9b78-a5a1093e0211",
	Exchange          string           `json:"Exchange"`          // : "BTC-LTC",
	OrderType         string           `json:"OrderType"`         // : "LIMIT_SELL",
	Quantity          decimal          `json:"Quantity"`          // : 5.00000000,
	QuantityRemaining decimal          `json:"QuantityRemaining"` // : 5.00000000,
	Limit             decimal          `json:"Limit"`             // : 2.00000000,
	CommissionPaid    decimal          `json:"CommissionPaid"`    // : 0.00000000,
	Price             decimal          `json:"Price"`             // : 0.00000000,
	PricePerUnit      decimal          `json:"PricePerUnit"`      // : null,
	Opened            BittrexTimestamp `json:"Opened"`            // : "2014-07-09T03:55:48.77",
	Closed            BittrexTimestamp `json:"Closed"`            // : null,
	CancelInitiated   bool             `json:"CancelInitiated"`   // : false,
	ImmediateOrCancel bool             `json:"ImmediateOrCancel"` // : false,
	IsConditional     bool             `json:"IsConditional"`     // : false,
	Condition         string           `json:"Condition"`         // : null,
	ConditionTarget   string           `json:"ConditionTarget"`   // : null
}

//AccountOrderDescription result body of /account/getorder
type AccountOrderDescription struct {
	AccountID                  string           `json:"AccountId"`                  // : null,
	OrderUUID                  string           `json:"OrderUuid"`                  // : "0cb4c4e4-bdc7-4e13-8c13-430e587d2cc1",
	Exchange                   string           `json:"Exchange"`                   // : "BTC-SHLD",
	Type                       string           `json:"Type"`                       // : "LIMIT_BUY",
	Quantity                   decimal          `json:"Quantity"`                   // : 1000.00000000,
	QuantityRemaining          decimal          `json:"QuantityRemaining"`          // : 1000.00000000,
	Limit                      decimal          `json:"Limit"`                      // : 0.00000001,
	Reserved                   decimal          `json:"Reserved"`                   // : 0.00001000,
	ReserveRemaining           decimal          `json:"ReserveRemaining"`           // : 0.00001000,
	CommissionReserved         decimal          `json:"CommissionReserved"`         // : 0.00000002,
	CommissionReserveRemaining decimal          `json:"CommissionReserveRemaining"` // : 0.00000002,
	CommissionPaid             decimal          `json:"CommissionPaid"`             // : 0.00000000,
	Price                      decimal          `json:"Price"`                      // : 0.00000000,
	PricePerUnit               decimal          `json:"PricePerUnit"`               // : null,
	Opened                     BittrexTimestamp `json:"Opened"`                     // : "2014-07-13T07:45:46.27",
	Closed                     BittrexTimestamp `json:"Closed"`                     // : null,
	IsOpen                     bool             `json:"IsOpen"`                     // : true,
	Sentinel                   string           `json:"Sentinel"`                   // : "6c454604-22e2-4fb4-892e-179eede20972",
	CancelInitiated            bool             `json:"CancelInitiated"`            // : false,
	ImmediateOrCancel          bool             `json:"ImmediateOrCancel"`          // : false,
	IsConditional              bool             `json:"IsConditional"`              // : false,
	Condition                  string           `json:"Condition"`                  // : "NONE",
	ConditionTarget            string           `json:"ConditionTarget"`            // : null
}

//AccountOrderHistoryDescription result element of /account/getorderhistory
type AccountOrderHistoryDescription struct {
	OrderUUID         string           `json:"OrderUuid"`         // : "fd97d393-e9b9-4dd1-9dbf-f288fc72a185",
	Exchange          string           `json:"Exchange"`          // : "BTC-LTC",
	TimeStamp         BittrexTimestamp `json:"TimeStamp"`         // : "2014-07-09T04:01:00.667",
	OrderType         string           `json:"OrderType"`         // : "LIMIT_BUY",
	Limit             decimal          `json:"Limit"`             // : 0.00000001,
	Quantity          decimal          `json:"Quantity"`          // : 100000.00000000,
	QuantityRemaining decimal          `json:"QuantityRemaining"` // : 100000.00000000,
	Commission        decimal          `json:"Commission"`        // : 0.00000000,
	Price             decimal          `json:"Price"`             // : 0.00000000,
	PricePerUnit      decimal          `json:"PricePerUnit"`      // : null,
	IsConditional     bool             `json:"IsConditional"`     // : false,
	Condition         string           `json:"Condition"`         // : null,
	ConditionTarget   string           `json:"ConditionTarget"`   // : null,
	ImmediateOrCancel bool             `json:"ImmediateOrCancel"` // : false
}

//TransactionHistoryDescription result element of /account/getwithdrawalhistory and /account/getdeposithistory
type TransactionHistoryDescription struct {
	PaymentUUID    string           `json:"PaymentUuid"`    // : "b52c7a5c-90c6-4c6e-835c-e16df12708b1",
	Currency       string           `json:"Currency"`       // : "BTC",
	Amount         decimal          `json:"Amount"`         // : 17.00000000,
	Address        string           `json:"Address"`        // : "1DeaaFBdbB5nrHj87x3NHS4onvw1GPNyAu",
	Opened         BittrexTimestamp `json:"Opened"`         // : "2014-07-09T04:24:47.217",
	Authorized     bool             `json:"Authorized"`     // : true,
	PendingPayment bool             `json:"PendingPayment"` // : false,
	TxCost         decimal          `json:"TxCost"`         // : 0.00020000,
	TxID           string           `json:"TxId"`           // : null,
	Canceled       bool             `json:"Canceled"`       // : true,
	InvalidAddress bool             `json:"InvalidAddress"` // : false
}

//AccountBalance result element as described under /account/getbalances. also the result body of /account/getbalance
type AccountBalance struct {
	Currency      string  `json:"Currency"`      // : "DOGE",
	Balance       decimal `json:"Balance"`       // : 0.00000000,
	Available     decimal `json:"Available"`     // : 0.00000000,
	Pending       decimal `json:"Pending"`       // : 0.00000000,
	CryptoAddress string  `json:"CryptoAddress"` // : "DLxcEt3AatMyr2NTatzjsfHNoB9NT62HiF",
	Requested     bool    `json:"Requested"`     // : false,
	UUID          string  `json:"Uuid"`          // : null
}

//WalletAddress result body of /account/getdepositaddress
type WalletAddress struct {
	Currency string `json:"Currency"` // : "VTC"
	Address  string `json:"Address"`  // : "Vy5SKeKGXUHKS2WVpJ76HYuKAu3URastUo"
}

//Candle result element as described under v2.0/pub/market/getticks
type Candle struct {
	TimeStamp BittrexTimestamp `json:"T"`
	Open      decimal          `json:"O"`
	Close     decimal          `json:"C"`
	High      decimal          `json:"H"`
	Low       decimal          `json:"L"`
	//Volume amount traded in the altcoin (Ex: the LTC in BTC-LTC)
	Volume decimal `json:"V"`
	//Volume amount traded in the base coin (Ex: the BTC in BTC-LTC)
	BaseVolume decimal `json:"BV"`
}

//OrderUpdate Update to an order listed under buys and sells in ExchangeState
type OrderUpdate struct {
	OrderElement //embed
	Type         int
}

//Fill structure found inside an ExchangeState object
type Fill struct {
	OrderElement //embed
	OrderType    string
	Timestamp    BittrexTimestamp
}

// ExchangeState contains fills and order book updates for a market.
type ExchangeState struct {
	MarketName string
	Nounce     int
	Buys       []OrderUpdate
	Sells      []OrderUpdate
	Fills      []Fill
	Initial    bool
}
