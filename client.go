package bittrex

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/technicalviking/bittrex2/signalr"
	"github.com/technicalviking/bittrex2/socketPayloads"
)

//Client object representing connection to bittrex api.
type Client struct {
	apiKey       string
	apiSecret    string
	timeout      time.Duration
	socketClient *signalr.Client

	orderSubscription   chan socketPayloads.OrderResponse
	balanceSubscription chan socketPayloads.BalanceDelta

	summaryDeltaMutex         sync.RWMutex
	summaryDeltaSubscriptions map[string]chan socketPayloads.Summary
	isSubbedToSummaryDelta    bool

	exchangeDeltaMutex         sync.RWMutex
	exchangeDeltaSubscriptions map[string]chan socketPayloads.ExchangeDelta

	summaryLiteDeltaMutex         sync.RWMutex
	summaryLiteDeltaSubscriptions map[string]chan socketPayloads.SummaryLiteDelta
	isSubbedToSummaryLiteDelta    bool

	errChanMutex sync.RWMutex
	errChan      chan error
}

//New construct a new Client object representing an interface to the various bittrex APIs.
func New(key string, secret string) (*Client, error) {
	newClient := &Client{
		apiKey:                        key,
		apiSecret:                     secret,
		timeout:                       time.Duration(defaultTimeout) * time.Second,
		summaryDeltaSubscriptions:     make(map[string]chan socketPayloads.Summary),
		summaryLiteDeltaSubscriptions: make(map[string]chan socketPayloads.SummaryLiteDelta),
		exchangeDeltaSubscriptions:    make(map[string]chan socketPayloads.ExchangeDelta),
		errChan:                       make(chan error, 5),
	}

	return newClient, nil
}

//ConnectWebSocket provide functionality to connect to the signalr endpoint.
func (c *Client) ConnectWebSocket() error {

	if c.socketClient != nil {
		c.socketClient.Close()
	}

	if newClientErr := c.connectNewSignalClient(); newClientErr != nil {
		return newClientErr
	}

	if c.apiKey != "" && c.apiSecret != "" {
		if authenticateErr := c.authNewSignalClient(); authenticateErr != nil {
			return authenticateErr
		}
	}

	c.addListeners()

	return nil
}

func (c *Client) connectNewSignalClient() error {
	client, clientErr := signalr.New()

	if clientErr != nil {
		return clientErr
	}

	if connectErr := client.Connect(websocketBaseURI, []string{websocketHub}); connectErr != nil {
		return fmt.Errorf("Unable to create bittrex signal client at url %s:  %+v", websocketBaseURI, connectErr)
	}

	c.socketClient = client

	return nil
}

//authNewSignalClient authenticate client to retrieve balance and order notifications
func (c *Client) authNewSignalClient() error {
	//authenticate the client.
	authContext, authErr := c.socketClient.CallHub(websocketHub, "GetAuthContext", c.apiKey)

	if authErr != nil {
		return fmt.Errorf("Unable to authenticate bittrex client: %+v", authErr)
	}

	var parsedAuthContext string
	json.Unmarshal(authContext, &parsedAuthContext)

	signedChallenge := c.sign(parsedAuthContext)

	challengeResp, challengeErr := c.socketClient.CallHub(websocketHub, "Authenticate", c.apiKey, signedChallenge)

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

	return nil
}

func (c *Client) addListeners() {
	//@TODO generate an error channel.
	c.socketClient.OnMessageError = c.socketOnErrorMethod

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
			c.pipeMarketExchangeDelta(decodedArg)
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

func (c *Client) socketOnErrorMethod(err error) {
	go func(e error) {
		c.errChanMutex.Lock()
		c.errChan <- err
		c.errChanMutex.Unlock()
	}(err)
}

//SubscribeToWebsocketErrors retrieve reference to error chan for websocket
func (c *Client) SubscribeToWebsocketErrors() chan error {
	return c.errChan
}

// GetWebSocketState passthrough of the Signalr state method.
func (c *Client) GetWebSocketState() signalr.ClientState {
	return c.socketClient.State()
}
