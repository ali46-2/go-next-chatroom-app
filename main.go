package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var subscribers = make(map[string][]*websocket.Conn)
var topics = []string{"anime", "books", "games", "movies", "music"}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("failed to create connection: ", err)
		return
	}

	defer conn.Close()

	urlParts := strings.Split(r.URL.Path, "/")
	topic := urlParts[len(urlParts)-1]

	subscribers[topic] = append(subscribers[topic], conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read message from connection %v: %v", conn.RemoteAddr(), err)
			return
		}

		for _, sub := range subscribers[topic] {
			go func(ws *websocket.Conn) {
				if err := sub.WriteMessage(messageType, message); err != nil {
					log.Printf("connection %v failed to write message: %v", conn.RemoteAddr(), err)
				}
			}(sub)
		}
	}
}

func main() {
	for _, topic := range topics {
		http.HandleFunc("/ws/"+topic, handleConnection)
	}

	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
