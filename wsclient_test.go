package wsclient

import (
	"testing"
)

func TestClient(t *testing.T) {

	ws := NewWSClient("wss://pubsub-local.comms.razerzone.com:7070/ws")
	ws.Connect()

}
