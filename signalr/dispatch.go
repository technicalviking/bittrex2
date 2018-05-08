package signalr

import (
	"encoding/json"
	"time"
)

type serverMessage struct {
	Cursor     string            `json:"C"`
	Data       []json.RawMessage `json:"M"`
	Result     json.RawMessage   `json:"R"`
	Identifier string            `json:"I"`
	Error      string            `json:"E"`
}

type hubCallResponse struct {
	HubName   string            `json:"H"`
	Method    string            `json:"M"`
	Arguments []json.RawMessage `json:"A"`
}

func (sc *Client) isDispatchRunning() bool {
	sc.dispatchMutex.RLock()
	defer sc.dispatchMutex.RUnlock()
	return sc.dispatchRunning

}

func (sc *Client) setDispatchState(newState bool) {
	sc.dispatchMutex.Lock()
	defer sc.dispatchMutex.Unlock()

	sc.dispatchRunning = newState
}

func (sc *Client) listenToWebSocket() chan serverMessage {
	socketDataChan := make(chan serverMessage)

	go func() {
		defer close(socketDataChan)
		for {
			_, data, err := sc.socket.ReadMessage()
			if err != nil {
				sc.outputError(err)
				return
			}

			var message serverMessage
			if err := json.Unmarshal(data, &message); err != nil {
				sc.outputError(newError("Unable to unmarshal message: %s\n", err.Error()))
			}

			sc.keepAliveMutex.Lock()
			sc.keepAliveTime = time.Now()
			sc.keepAliveMutex.Unlock()

			socketDataChan <- message
		}
	}()

	return socketDataChan
}

// Start dispatch loop. This function will return when error occurs. When this
// happens, all the connections are closed and user can run Connect()
// and Dispatch() again on the same client.
func (sc *Client) dispatch() {

	if sc.isDispatchRunning() {
		return
	}

	sc.state = Connected
	sc.setDispatchState(true)
	t := time.NewTicker(time.Second)

	defer func() {
		sc.setDispatchState(false)
		t.Stop()
		if e := sc.reconnectWebsocket(); e != nil {
			sc.state = Disconnected
			sc.outputError(e)
			return
		}

		go sc.dispatch()
	}()

	dataChan := sc.listenToWebSocket()

	for {
		select {
		case data, ok := <-dataChan:
			if !ok {
				return
			}
			sc.handleSocketData(data)
		case <-t.C:
			sc.keepAliveMutex.RLock()
			keepAliveTime := sc.keepAliveTime
			sc.keepAliveMutex.RUnlock()

			if time.Since(keepAliveTime) > time.Duration(sc.negotiationParams.KeepAliveTimeout)*time.Second {
				sc.socket.Close()
				sc.outputError(newError("keepalive timeout reached.  RECONNECTING."))
				return
			}
		}
	}
}

func (sc *Client) handleSocketData(message serverMessage) {

	// This is a response to a hub call.
	if len(message.Identifier) > 0 {
		sc.routeResponse(&message)
		return
	}

	for _, curData := range message.Data {
		var hubCall hubCallResponse

		if err := json.Unmarshal(curData, &hubCall); err != nil {
			sc.outputError(newError("Unable to unmarshal message data: %s", err.Error()))
			continue
		}

		// check if this is a client Hub method call from server.
		if hubCall.HubName != "" && hubCall.Method != "" && sc.OnClientMethod != nil {
			sc.OnClientMethod(hubCall.HubName, hubCall.Method, hubCall.Arguments)
		}
	}
}
