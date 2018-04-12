//bulk of this code borrowed from github.com/hweom/signalr, not forked because of debugging I had to do to get the bittrex stuff to work.

package signalr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type negotiationResponse struct {
	URL                     string
	ConnectionToken         string
	ConnectionID            string
	KeepAliveTimeout        float32
	DisconnectTimeout       float32
	ConnectionTimeout       float32
	TryWebSockets           bool
	ProtocolVersion         string
	TransportConnectTimeout float32
	LogPollDelay            float32
}

//Client object representing connection to the signalr socket api
type Client struct {
	OnMessageError func(err error)
	OnClientMethod func(hub, method string, arguments []json.RawMessage)
	// When client disconnects, the causing error is sent to this channel. Valid only after Connect().
	DisconnectedChannel chan bool
	// Additional header parameters to add to the negotiation HTTP request.
	RequestHeader http.Header

	params negotiationResponse
	socket *websocket.Conn
	nextID int

	// Futures for server call responses and a guarding mutex.
	responseFutures map[string]chan *serverMessage
	mutex           sync.Mutex
	dispatchRunning bool
}

type serverMessage struct {
	Cursor     string            `json:"C"`
	Data       []json.RawMessage `json:"M"`
	Result     json.RawMessage   `json:"R"`
	Identifier string            `json:"I"`
	Error      string            `json:"E"`
}

func (sc *Client) connectWebsocket(address string, params negotiationResponse, hubs []string) (*websocket.Conn, error) {
	var connectionData = make([]struct {
		Name string `json:"Name"`
	}, len(hubs))
	for i, h := range hubs {
		connectionData[i].Name = h
	}
	connectionDataBytes, err := json.Marshal(connectionData)
	if err != nil {
		return nil, err
	}

	var connectionParameters = url.Values{}
	connectionParameters.Set("transport", "webSockets")
	connectionParameters.Set("clientProtocol", "1.5")
	connectionParameters.Set("connectionToken", params.ConnectionToken)
	connectionParameters.Set("connectionData", string(connectionDataBytes))

	var connectionURL = url.URL{Scheme: "wss", Host: address, Path: "signalr/connect"}
	connectionURL.RawQuery = connectionParameters.Encode()

	if conn, _, err := websocket.DefaultDialer.Dial(connectionURL.String(), sc.RequestHeader); err != nil {
		return nil, err
	} else {
		return conn, nil
	}
}

func (sc *Client) negotiate(scheme, address string) (negotiationResponse, error) {
	var response negotiationResponse

	var negotiationURL = url.URL{Scheme: scheme, Host: address, Path: "/signalr/negotiate"}

	client := &http.Client{}

	request, err := http.NewRequest("GET", negotiationURL.String(), nil)

	if err != nil {
		return response, err
	}

	for k, values := range sc.RequestHeader {
		for _, val := range values {
			request.Header.Add(k, val)
		}
	}

	reply, err := client.Do(request)
	if err != nil {
		return response, err
	}

	defer reply.Body.Close()

	if body, err := ioutil.ReadAll(reply.Body); err != nil {
		return response, err
	} else if err := json.Unmarshal(body, &response); err != nil {
		return response, fmt.Errorf("Failed to parse message '%s': %s", string(body), err.Error())
	} else {
		return response, nil
	}
}

func (sc *Client) routeResponse(response *serverMessage) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if c, ok := sc.responseFutures[response.Identifier]; ok {
		c <- response
		close(c)
		delete(sc.responseFutures, response.Identifier)
	}
}

func (sc *Client) createResponseFuture(identifier string) (chan *serverMessage, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if !sc.dispatchRunning {
		return nil, fmt.Errorf("Dispatch is not running")
	}

	var c = make(chan *serverMessage)
	sc.responseFutures[identifier] = c

	return c, nil
}

func (sc *Client) deleteResponseFuture(identifier string) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	delete(sc.responseFutures, identifier)
}

func (sc *Client) tryStartDispatch() error {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if sc.dispatchRunning {
		return fmt.Errorf("Another Dispatch() is running")
	}
	sc.DisconnectedChannel = make(chan bool)
	sc.dispatchRunning = true

	return nil
}

func (sc *Client) endDispatch() {
	// Close all the waiting response futures.
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.dispatchRunning = false
	for _, c := range sc.responseFutures {
		close(c)
	}
	sc.responseFutures = make(map[string]chan *serverMessage)
	close(sc.DisconnectedChannel)
}

// Start dispatch loop. This function will return when error occurs. When this
// happens, all the connections are closed and user can run Connect()
// and Dispatch() again on the same client.
func (sc *Client) dispatch(connectedChannel chan bool) {
	if err := sc.tryStartDispatch(); err != nil {
		panic("Dispatch is already running")
	}

	defer sc.endDispatch()

	close(connectedChannel)

	for {
		var message serverMessage

		var hubCall struct {
			HubName   string            `json:"H"`
			Method    string            `json:"M"`
			Arguments []json.RawMessage `json:"A"`
		}

		_, data, err := sc.socket.ReadMessage()
		if err != nil {
			sc.socket.Close()
			break
		} else if err := json.Unmarshal(data, &message); err != nil {
			if sc.OnMessageError != nil {
				sc.OnMessageError(err)
			}
		} else {
			if len(message.Identifier) > 0 {
				// This is a response to a hub call.
				sc.routeResponse(&message)
			} else if len(message.Data) == 1 {
				if err := json.Unmarshal(message.Data[0], &hubCall); err == nil && len(hubCall.HubName) > 0 && len(hubCall.Method) > 0 {
					// This is a client Hub method call from server.
					if sc.OnClientMethod != nil {
						sc.OnClientMethod(hubCall.HubName, hubCall.Method, hubCall.Arguments)
					}
				}
			}
		}
	}
}

//CallHub Call server hub method. Dispatch() function must be running, otherwise this method will never return.
func (sc *Client) CallHub(hub, method string, params ...interface{}) (json.RawMessage, error) {
	var request = struct {
		Hub        string        `json:"H"`
		Method     string        `json:"M"`
		Arguments  []interface{} `json:"A"`
		Identifier int           `json:"I"`
	}{
		Hub:        hub,
		Method:     method,
		Arguments:  params,
		Identifier: sc.nextID,
	}

	sc.nextID++

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	var responseKey = fmt.Sprintf("%d", request.Identifier)
	responseChannel, err := sc.createResponseFuture(responseKey)
	if err != nil {
		return nil, err
	}

	if err := sc.socket.WriteMessage(websocket.TextMessage, data); err != nil {
		return nil, err
	}

	defer sc.deleteResponseFuture(responseKey)

	if response, ok := <-responseChannel; !ok {
		return nil, fmt.Errorf("Call to server returned no result")
	} else if len(response.Error) > 0 {
		return nil, fmt.Errorf("%s", response.Error)
	} else {
		return response.Result, nil
	}
}

//Connect create websocket connection to SignalR endpoint.
func (sc *Client) Connect(scheme, host string, hubs []string) error {
	var params negotiationResponse
	var err error

	// Negotiate parameters.
	if params, err = sc.negotiate(scheme, host); err != nil {
		return err
	}

	sc.params = params

	var ws *websocket.Conn

	// Connect Websocket.
	if ws, err = sc.connectWebsocket(host, sc.params, hubs); err != nil {
		return err
	}

	sc.socket = ws

	var connectedChannel = make(chan bool)
	go sc.dispatch(connectedChannel)
	<-connectedChannel

	return nil
}

//Close close the websocket connection
func (sc *Client) Close() {
	sc.socket.Close()
}

//New constructor for SignalR connection client.
func New() *Client {
	return &Client{
		RequestHeader:   http.Header{},
		nextID:          1,
		responseFutures: make(map[string]chan *serverMessage),
	}
}
