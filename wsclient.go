/*
Package wsclient implements a WebSocket client.

Example:

	ws := wsclient.NewWSClient("ws://localhost:7070/ws")
	err := ws.Connect()
	if err != nil {
		panic(err)
	}
	ws.OnOpen(func() {
		fmt.Printf("connection opened")
	})
	ws.OnMessage(func(data []byte) {
		fmt.Printf("got message")
	})
	ws.OnClose(func() {
		fmt.Println("connection closed")
	})
	ws.SendJSON(wsclient.M{
		"type": "chat",
		"payload": "hello world",
		"sender": {
			"name": "Bob",
		},
	})

*/
package wsclient

// WSClient is a WebSocket client
type WSClient struct {
}

// M is a convenient alias for map[string]interface{}
type M map[string]interface{}

// NewWSClient returns a new instance of WSClient given the WebSocket URL
func NewWSClient(url string) *WSClient {
	return &WSClient{}
}

// OnOpen is a callback function when the connection is opened
func (c *WSClient) OnOpen(fn func()) {
}

// OnMessage is the callback function when a data is received from the server
func (c *WSClient) OnMessage(fn func(data []byte)) {
}

// OnClose is a the callback function when the connection is closed
func (c *WSClient) OnClose(fn func()) {
}

// Connect connects to the WebSocket server
func (c *WSClient) Connect() error {
	return nil
}

// SendJSON sends a JSON encoded message to the server
func (c *WSClient) SendJSON(j M) error {
	return nil
}

// Close closes the connection from the server
func (c *WSClient) Close() error {
	return nil
}
