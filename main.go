package main

import (
	"encoding/json"
	"fmt"

	"github.com/technicalviking/bittrex2/signalr"
)

func main() {
	client := signalr.New()

	client.OnClientMethod = func(hub, method string, arguments []json.RawMessage) {
		fmt.Println("Message Received: ")
		fmt.Println("HUB: ", hub)
		fmt.Println("METHOD: ", method)
		fmt.Println("ARGUMENTS: ", arguments)
	}

	client.OnMessageError = func(err error) {
		fmt.Println("ERROR OCCURRED: ", err)
	}

	client.Connect("https", "beta.bittrex.com/signalr", "c2")

	client.CallHub("c2", "GetAuthContext", "")
}
