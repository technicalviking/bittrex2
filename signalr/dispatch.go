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

func (sc *Client) beginDispatch() {
	for {
		sc.dispatch()
		if err := sc.reconnectWebsocket(); err != nil {
			sc.setState(Disconnected)
			sc.outputError(err)
			return
		}
	}
}

// Start dispatch loop. This function will return when error occurs. When this
// happens, all the connections are closed and user can run Connect()
// and Dispatch() again on the same client.
func (sc *Client) dispatch() {
	if sc.isDispatchRunning() {
		return
	}

	sc.setState(Connected)
	sc.setDispatchState(true)
	defer sc.setDispatchState(false)

	t := time.NewTicker(time.Second)
	dataChan := sc.listenToWebSocket()

	for {
		select {
		case data, ok := <-dataChan:
			if !ok {
				t.Stop()
				return
			}
			sc.handleSocketData(data)
		case <-t.C:
			sc.keepAliveMutex.RLock()
			keepAliveTime := sc.keepAliveTime
			sc.keepAliveMutex.RUnlock()

			if time.Since(keepAliveTime) > time.Duration(sc.negotiationParams.KeepAliveTimeout)*time.Second {
				t.Stop()
				sc.socket.Close()
				sc.outputError(newError("keepalive timeout reached.  RECONNECTING."))
				return
			}
		}
	}
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
			var (
				data []byte
				err  error
			)

			if _, data, err = sc.socket.ReadMessage(); err != nil {
				sc.outputError(err)
				return
			}

			sc.updateKeepAlive()

			var message serverMessage
			if err = json.Unmarshal(data, &message); err != nil {
				sc.outputError(newError("Unable to parse message: %s\n", err.Error()))
				continue
			}
			socketDataChan <- message
		}
	}()

	return socketDataChan
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
			go sc.OnClientMethod(hubCall.HubName, hubCall.Method, hubCall.Arguments)
		}
	}
}
