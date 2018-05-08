package signalr

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

type callHubRequest struct {
	Hub        string        `json:"H"`
	Method     string        `json:"M"`
	Arguments  []interface{} `json:"A"`
	Identifier string        `json:"I"`
}

func (sc *Client) newCallHubRequest(hub, method string, params []interface{}) callHubRequest {
	sc.callHubIDMutex.Lock()
	requestID := sc.nextID
	sc.nextID++
	sc.callHubIDMutex.Unlock()

	return callHubRequest{
		Hub:        hub,
		Method:     method,
		Arguments:  params,
		Identifier: fmt.Sprintf("%d", requestID),
	}
}

func (sc *Client) sendHubMessage(data []byte) error {
	sc.callHubMutex.Lock()
	defer sc.callHubMutex.Unlock()

	return sc.socket.WriteMessage(websocket.TextMessage, data)
}

//CallHub Call server hub method. Dispatch() function must be running, otherewise this method will return an error.
func (sc *Client) CallHub(hub, method string, params ...interface{}) (json.RawMessage, error) {
	if !sc.isDispatchRunning() {
		return nil, errors.New("dispatch not running")
	}

	request := sc.newCallHubRequest(hub, method, params)

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	if err := sc.sendHubMessage(data); err != nil {
		return nil, err
	}

	responseKey := request.Identifier
	responseChannel := sc.createResponseFuture(responseKey)

	var (
		response *serverMessage
		ok       bool
	)

	if response, ok = <-responseChannel; !ok {
		return nil, fmt.Errorf("Call to server returned no result")
	}

	if len(response.Error) > 0 {
		return nil, fmt.Errorf("%s", response.Error)
	}

	return response.Result, nil
}

func (sc *Client) createResponseFuture(identifier string) chan *serverMessage {
	sc.responseMutex.Lock()
	defer sc.responseMutex.Unlock()

	sc.responseFutures[identifier] = make(chan *serverMessage)
	return sc.responseFutures[identifier]
}

//this method should only be called when a message has an identifier corresponding with a request.
func (sc *Client) routeResponse(response *serverMessage) {
	sc.responseMutex.RLock()
	c, ok := sc.responseFutures[response.Identifier]
	defer sc.responseMutex.RUnlock()

	if ok {
		c <- response
		sc.deleteResponseFuture(response.Identifier)
	}
}

func (sc *Client) deleteResponseFuture(identifier string) {
	sc.responseMutex.Lock()
	defer sc.responseMutex.Unlock()

	close(sc.responseFutures[identifier])
	delete(sc.responseFutures, identifier)
}
