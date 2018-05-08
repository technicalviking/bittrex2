package signalr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	defaultScheme     string = "https"
	socketScheme      string = "wss"
	signalREndpoint   string = "signalr"
	negotiatePath     string = signalREndpoint + "/negotiate"
	connectEndpoint   string = signalREndpoint + "/connect"
	reconnectEndpoint string = signalREndpoint + "/reconnect"
)

type negotiationResponse struct {
	ConnectionToken string
	//commenting these out because they're not currently used, but just in case I need them in the future.
	URL                     string
	ConnectionID            string
	KeepAliveTimeout        float32
	DisconnectTimeout       float32
	ConnectionTimeout       float32
	TryWebSockets           bool
	ProtocolVersion         string
	TransportConnectTimeout float32
	LogPollDelay            float32
}

//Connect create websocket connection to SignalR endpoint.
func (sc *Client) Connect(connectURL string, hubs []string) error {

	sc.setConnectionURL(connectURL)
	sc.hubs = hubs

	var err error

	sc.state = Connecting

	// Negotiate parameters.
	if err = sc.negotiate(); err != nil {
		sc.state = Disconnected
		return err
	}

	// Connect Websocket.
	if err = sc.connectWebsocket(); err != nil {
		sc.state = Disconnected
		return err
	}

	go sc.dispatch()
	return nil
}

func (sc *Client) setConnectionURL(connectURL string) error {
	var err error
	if sc.url, err = url.Parse(connectURL); err != nil {
		return err
	}

	//if the user didn't define a scheme, set it to https.
	if sc.url.Scheme == "" {
		sc.url.Scheme = "https"
	}

	return nil
}

func (sc *Client) getConnectionURL() *url.URL {
	//make a copy of the stored base url.
	result, _ := url.Parse(sc.url.String())
	return result
}

func (sc *Client) negotiate() error {
	var (
		request  *http.Request
		response *http.Response
		result   negotiationResponse
		err      error
		body     []byte
	)

	negotiationURL := sc.getConnectionURL()
	negotiationURL.Path = negotiatePath

	query := negotiationURL.Query()
	query.Set("clientProtocol", "1.5")
	query.Set("_", fmt.Sprintf("%d", time.Now().Unix()*1000)) //prevent 304 responses.

	negotiationURL.RawQuery = query.Encode()

	if request, err = http.NewRequest("GET", negotiationURL.String(), nil); err != nil {
		return err
	}

	for k, values := range sc.RequestHeader {
		for _, val := range values {
			request.Header.Add(k, val)
		}
	}

	if response, err = sc.client.Do(request); err != nil {
		return err
	}

	defer response.Body.Close()

	if body, err = ioutil.ReadAll(response.Body); err != nil {
		return err
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("Failed to parse message '%s': %s", string(body), err.Error())
	}

	sc.negotiationParams = result
	return nil
}

func (sc *Client) connectWebsocket() error {
	var err error

	connectionURL := sc.getConnectionURL()
	connectionURL.Scheme = socketScheme
	connectionURL.Path = connectEndpoint
	connectionURL.RawQuery = url.Values{
		"transport":       []string{"webSockets"},
		"clientProtocol":  []string{sc.negotiationParams.ProtocolVersion},
		"connectionToken": []string{sc.negotiationParams.ConnectionToken},
		"connectionData":  []string{string(castNamesToString(sc.hubs))},
		"_":               []string{fmt.Sprintf("%d", time.Now().Unix()*1000)},
	}.Encode()

	socketDialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		Jar:              sc.client.Jar,
	}

	i := 0
	for ; i < sc.maxRetries; i++ {
		backoff := math.Pow(2.0, float64(i))
		time.Sleep(time.Second * time.Duration(backoff))
		if sc.socket, _, err = socketDialer.Dial(connectionURL.String(), sc.RequestHeader); err != nil {
			sc.outputError(err)
			continue
		}

		break
	}

	if i == sc.maxRetries {
		return newError("MAX RETRIES REACHED.  ABORTING CONNECTION.")
	}

	return nil
}

func (sc *Client) reconnectWebsocket() error {
	var err error

	sc.state = Reconnecting

	connectionURL := sc.getConnectionURL()
	connectionURL.Scheme = socketScheme
	connectionURL.Path = connectEndpoint
	connectionURL.RawQuery = url.Values{
		"transport":       []string{"webSockets"},
		"clientProtocol":  []string{sc.negotiationParams.ProtocolVersion},
		"connectionToken": []string{sc.negotiationParams.ConnectionToken},
		"connectionData":  []string{string(castNamesToString(sc.hubs))},
		"messageId":       []string{sc.lastMessageID},
	}.Encode()

	socketDialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		Jar:              sc.client.Jar,
	}

	i := 0
	for ; i < sc.maxRetries; i++ {
		backoff := math.Pow(2.0, float64(i))
		time.Sleep(time.Second * time.Duration(backoff))
		if sc.socket, _, err = socketDialer.Dial(connectionURL.String(), sc.RequestHeader); err != nil {
			sc.outputError(err)
			continue
		}

		break
	}

	if i == sc.maxRetries {
		return newError("MAX RETRIES REACHED.  ABORTING CONNECTION.")
	}

	return nil
}

func castNamesToString(hubs []string) []byte {
	var connectionData = make([]struct {
		Name string `json:"Name"`
	}, len(hubs))
	for i, h := range hubs {
		connectionData[i].Name = h
	}
	connectionDataBytes, _ := json.Marshal(connectionData)

	return connectionDataBytes
}
