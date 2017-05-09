/*
Package wsclient implements a WebSocket client.

Example:

	ws := wsclient.NewWSClient("ws://localhost:7070/ws")

	ws.OnOpen(func() {
		fmt.Printf("connection opened")
		ws.SendJSON(wsclient.M{
			"type": "chat",
			"payload": "hello world",
			"sender": {
				"name": "Bob",
			},
		})
	})
	ws.OnMessage(func(data []byte) {
		fmt.Printf("got message")
	})
	ws.OnClose(func() {
		fmt.Println("connection closed")
	})
	ws.OnError(func(err error) {
		panic(err)
	})
	ws.Connect()


*/
package wsclient

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSClient is a WebSocket client
type WSClient struct {
	u        string
	ws       *websocket.Conn
	send     chan []byte
	closed   bool
	closedMu sync.RWMutex

	onOpen    func()
	onMessage func(data []byte)
	onClose   func()
	onError   func(e error)
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

// M is a convenient alias for map[string]interface{}
type M map[string]interface{}

// NewWSClient returns a new instance of WSClient given the WebSocket URL
func NewWSClient(url string) *WSClient {
	return &WSClient{
		u:    url,
		send: make(chan []byte),
	}
}

// OnOpen is a callback function when the connection is opened
func (c *WSClient) OnOpen(fn func()) {
	c.onOpen = fn
}

// OnMessage is the callback function when a data is received from the server
func (c *WSClient) OnMessage(fn func(data []byte)) {
	c.onMessage = fn
}

// OnClose is the callback function when the connection is closed
func (c *WSClient) OnClose(fn func()) {
	c.onClose = fn
}

// OnError is a callback function for handling errors
func (c *WSClient) OnError(fn func(err error)) {
	c.onError = fn
}

// Connect connects to the WebSocket server
func (c *WSClient) Connect() {
	go func() {
		var err error
		//log.Printf("wsclient connecting to: %s", c.u)
		c.ws, _, err = websocket.DefaultDialer.Dial(c.u, nil)
		if err != nil {
			fmt.Printf("Connect error: %s", err.Error())
			if c.onError != nil {
				c.onError(err)
			}
			return
		}
		//log.Printf("wsclient connected to: %s", c.u)
		go c.writePump()
		go c.readPump()

		if c.onOpen != nil {
			c.onOpen()
		}
	}()
}

// SendJSON sends a JSON encoded message to the server
func (c *WSClient) SendJSON(j M) error {

	b, err := json.Marshal(j)
	if err != nil {
		log.Printf("SendJSON: Marshal error: %s", err.Error())
		return err
	}
	//log.Printf("Sending: '%s'", string(b))

	c.send <- b

	return nil
}

// Close closes the connection from the server
func (c *WSClient) Close() {
	go func() {
		if c.isClosed() {
			log.Printf("Close: already closed")
			return
		}
		c.closedMu.Lock()
		c.closed = true
		c.closedMu.Unlock()
		if c.ws != nil {
			c.ws.Close()
		}
		if c.onClose != nil {
			c.onClose()
		}
		close(c.send)
		log.Printf("Close done")
	}()
	return
}

func (c *WSClient) writePump() {
	defer func() {
		c.Close()
		log.Printf("writePump: done")
	}()
	for {
		select {
		case mesg, ok := <-c.send:
			if !ok {
				return
			}
			if err := c.write(websocket.TextMessage, mesg); err != nil {
				log.Printf("write: error: %s", err.Error())
				return
			}
		}
	}
}

func (c *WSClient) readPump() {
	defer func() {
		c.Close()
		log.Printf("readPump: done")
	}()
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Read error: %s", err.Error())
			}
			break
		}
		if c.onMessage != nil {
			c.onMessage(message)
		}
	}
}

func (c *WSClient) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	if mt != websocket.CloseMessage {
		if mt == websocket.PingMessage {
			log.Printf("mt: ping")
		} else {
			//log.Printf("mt: %d write: '%s'", mt, string(payload))
		}
	}
	return c.ws.WriteMessage(mt, payload)
}

func (c *WSClient) isClosed() bool {
	c.closedMu.RLock()
	defer c.closedMu.RUnlock()
	return c.closed
}
