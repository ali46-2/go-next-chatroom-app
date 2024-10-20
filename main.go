package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var subscribers = make(map[*websocket.Conn]bool)

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("failed to create connection: ", err)
		return
	}

	defer conn.Close()

	subscribers[conn] = true

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read message from connection %v: %v", conn.RemoteAddr(), err)
			return
		}

		for sub := range subscribers {
			go func(ws *websocket.Conn) {
				if err := sub.WriteMessage(messageType, message); err != nil {
					log.Printf("connection %v failed to write message: %v", conn.RemoteAddr(), err)
				}
			}(sub)
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)
	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
