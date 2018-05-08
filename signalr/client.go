//bulk of this code borrowed from github.com/hweom/signalr, not forked because of debugging I had to do to get the bittrex stuff to work.

package signalr

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/technicalviking/bittrex2/cloudflare"
)

//ClientState int representing current state of the SignalR Client
type ClientState int

//SignalR Client State Values
const (
	Disconnected ClientState = iota
	Connecting
	Reconnecting
	Connected
)

//Client object representing connection to the signalr socket api
type Client struct {
	state ClientState
	//When errors happen for any reason, this callback is called.  This includes when the websocket closes remotely.
	OnMessageError func(err error)
	//This method is called whenever a message comes down through the websocket.
	OnClientMethod func(hub, method string, arguments []json.RawMessage)
	// Additional header parameters to add to the negotiation HTTP request.
	RequestHeader http.Header

	negotiationParams negotiationResponse
	keepAliveTime     time.Time
	keepAliveMutex    sync.RWMutex

	socket       *websocket.Conn
	callHubMutex sync.Mutex

	nextID         int
	callHubIDMutex sync.Mutex

	// Futures for server call responses and a guarding mutex.
	responseFutures map[string]chan *serverMessage
	responseMutex   sync.RWMutex

	//keep track of last message id in case reconnect is needed.
	lastMessageID string

	dispatchRunning bool
	dispatchMutex   sync.RWMutex

	//setting a persisting http client to allow for the usage of cloudflare scraper.
	client *http.Client

	url *url.URL

	hubs []string

	maxRetries int
}

//Close close the websocket connection
func (sc *Client) Close() {
	sc.socket.Close()
}

//SetMaxRetries - number of times the client will try to automatically reconnect.
func (sc *Client) SetMaxRetries(retries int) {
	sc.maxRetries = retries
}

//State get connected state.
func (sc *Client) State() ClientState {
	return sc.state
}

//New constructor for SignalR connection client.
func New() (*Client, error) {

	scraper, err := cloudflare.NewTransport(http.DefaultTransport)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Transport: scraper}

	return &Client{
		RequestHeader:   http.Header{},
		nextID:          1,
		responseFutures: make(map[string]chan *serverMessage),
		client:          client,
		maxRetries:      5,
	}, nil
}
