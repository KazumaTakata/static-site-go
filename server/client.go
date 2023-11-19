package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request, fileChanged chan bool) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {

		<-fileChanged

		fmt.Println("filechanged!!")

		w, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write([]byte("changed!!"))

		// Add queued chat messages to the current websocket message.
		if err := w.Close(); err != nil {
			return
		}
	}

}
