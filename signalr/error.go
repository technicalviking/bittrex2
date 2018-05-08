package signalr

import "fmt"

//Error implement error interface in a way that can be detected with type assertion.
type Error string

//Error implement Error interface
func (s Error) Error() string {
	return string(s)
}

func newError(format string, params ...interface{}) Error {
	return Error(fmt.Sprintf(format, params...))
}

func (sc *Client) outputError(e error) {
	if sc.OnMessageError != nil {
		sc.OnMessageError(e)
	}
}
