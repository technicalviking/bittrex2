package signalr

import (
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/gorilla/websocket"
)

//Client object representing SignalR connection.
type Client struct {
	socketURL  string
	socketConn *websocket.Conn
}

//Connect establish socket connection to signalr endpoint.
func (c *Client) Connect() error {
	conn, resp, connectErr := websocket.DefaultDialer.Dial(c.socketURL, nil)

	if connectErr != nil {
		return connectErr
	}

	if resp.StatusCode != 200 {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Unable to connect websocket to signalr hub: %s", string(responseBody))
	}

	c.socketConn = conn
	return nil
}

//New constructor for Signalr client
func New(socketURL string) (*Client, error) {

	_, parseErr := url.Parse(socketURL)

	if parseErr != nil {
		return nil, fmt.Errorf("Unable to parse URL: %+v", parseErr)
	}

	return &Client{
		socketURL: socketURL,
	}, nil
}
