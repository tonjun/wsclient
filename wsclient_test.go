package wsclient

import (
	"log"
	"testing"
	//"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	done := make(chan bool)

	ws := NewWSClient("ws://localhost:8080")
	ws.Connect()
	ws.OnOpen(func() {
		log.Println("connection opened")

		ws.SendJSON(M{
			"op": "get-time",
		})

		ws.OnMessage(func(data []byte) {
			log.Printf("OnMessage: '%s'", string(data))
			assert.Equal(t, []byte("{\"op\":\"get-time-response\"}"), data)

			//time.Sleep(10 * time.Millisecond)
			done <- true
		})

	})

	log.Printf("waiting to finish")
	<-done

	ws.OnClose(func() {
		log.Printf("onClose")
		done <- true
	})
	ws.Close()

	<-done

	//time.Sleep(100 * time.Millisecond)
}
