package bittrex

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/technicalviking/bittrex2/bittrex/socketPayloads"
	"github.com/technicalviking/bittrex2/signalr"
)

//Client object representing connection to bittrex api.
type Client struct {
	apiKey       string
	apiSecret    string
	timeout      time.Duration
	socketClient *signalr.Client

	summaryDeltaMutex         sync.Mutex
	summaryDeltaSubscriptions map[string]chan socketPayloads.SummaryDelta

	exchangeDeltaMutex         sync.Mutex
	exchangeDeltaSubscriptions map[string]chan socketPayloads.ExchangeDelta

	summaryLiteDeltaMutex         sync.Mutex
	summaryLiteDeltaSubscriptions map[string]chan socketPayloads.SummaryLiteDelta
}

//New construct a new Client object representing an interface to the various bittrex APIs.
func New(key string, secret string) (*Client, error) {
	newClient := &Client{
		apiKey:    key,
		apiSecret: secret,
		timeout:   time.Duration(defaultTimeout) * time.Second,
	}

	if newClientErr := newClient.connectNewSignalClient(); newClientErr != nil {
		return nil, newClientErr
	}

	newClient.addListeners()

	return newClient, nil
}

func (c *Client) connectNewSignalClient() error {
	client := signalr.New()

	socketURL, _ := url.Parse(websocketBaseURI)

	if connectErr := client.Connect(socketURL.Scheme, socketURL.Host, []string{websocketHub}); connectErr != nil {
		return fmt.Errorf("Unable to create bittrex client: %+v", connectErr)
	}

	//authenticate the client.
	authContext, authErr := client.CallHub(websocketHub, "GetAuthContext", c.apiKey)

	if authErr != nil {
		return fmt.Errorf("Unable to authenticate bittrex client: %+v", authErr)
	}

	var parsedAuthContext string
	json.Unmarshal(authContext, &parsedAuthContext)

	signedChallenge := c.sign(parsedAuthContext)

	challengeResp, challengeErr := client.CallHub(websocketHub, "Authenticate", c.apiKey, signedChallenge)

	if challengeErr != nil {
		return fmt.Errorf("Signed challenge not accepted: %+v", challengeErr)
	}

	var challengeOK bool

	parseErr := json.Unmarshal(challengeResp, &challengeOK)

	if parseErr != nil {
		return fmt.Errorf("Unable to parse response from authenticate call: %+v", parseErr)
	}

	if challengeOK == false {
		return fmt.Errorf("Signed challenge not accepted, no error")
	}

	c.socketClient = client

	return nil
}

func (c *Client) addListeners() {
	//@TODO generate an error channel.
	c.socketClient.OnMessageError = func(err error) {
		fmt.Println("ERROR OCCURRED: ", err)
	}

	c.socketClient.OnClientMethod = c.socketOnClientMethod
}

func (c *Client) socketOnClientMethod(hub, method string, arguments []json.RawMessage) {
	for _, arg := range arguments {
		var parseErr error
		var decodedArg []byte
		decodedArg, parseErr = socketPayloads.Parse(arg)

		if parseErr != nil {
			fmt.Println("parse error!", parseErr.Error())
		}

		switch method {
		case eventOrderDelta:
			c.pipeEventOrderDelta(decodedArg)
		case eventBalanceDelta:
			c.pipeBalanceDelta(decodedArg)
		case eventMarketDelta:
			c.pipeMarketDelta(decodedArg)
		case eventSummaryDelta:
			c.pipeEventSummaryDelta(decodedArg)
		case eventSummaryDeltaLite:
			c.pipeEventSummaryDeltaLite(decodedArg)
		default:
			//@TODO pass this into an error channel.
			fmt.Printf("unknown event type: %s \n", method)
			fmt.Printf("unknown event body %s", string(decodedArg))
		}

	}

}
