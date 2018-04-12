package socketPayloads

import (
	"encoding/json"
	"time"
)

//creating type aliases to be used internally here.
type guid = string

/*
bittrex describes decimal values as "string encoded decimals", but the json body does not actually wrap the values in quotes.
Thankfully, the value as originally described is within the bounds of a float64 (18 significant digits, with a precision of 8),
so shouldn't be too much data loss for the consumer of this library to transform the value into a better type.

NOTE:  I'm not using shopspring/decimal here because that library's performance is fucking garbage.  Any significant operation within that library
involves a call to the function 'rescale', which creates (potentially multiple) temporary 'Decimal' values which only live for the purposes of the operation.
*/
type decimal = float64

//unmarshal dates.
type date time.Time

func (d *date) UnmarshalJSON(raw []byte) error {
	var timestamp int64
	if timeParseErr := json.Unmarshal(raw, &timestamp); timeParseErr != nil {
		return timeParseErr
	}
	//the timestamp given is in milliseconds.  I wonder if bittrex is assuming node users?
	*d = date(time.Unix(timestamp/1000, 0))
	return nil
}

func (d *date) Get() time.Time {
	return time.Time(*d)
}
