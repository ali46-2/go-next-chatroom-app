package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type connections map[*websocket.Conn]struct{}

var (
	upgrader    = websocket.Upgrader{}
	subscribers = make(map[string]connections)
	topics      = []string{"anime", "books", "games", "movies", "music"}
	mutex       sync.Mutex
)

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("failed to create connection: ", err)
		return
	}

	urlParts := strings.Split(r.URL.Path, "/")
	topic := urlParts[len(urlParts)-1]

	mutex.Lock()
	if _, exists := subscribers[topic]; !exists {
		subscribers[topic] = make(connections)
	}
	subscribers[topic][conn] = struct{}{}
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		conn.Close()
		delete(subscribers[topic], conn)
		mutex.Unlock()
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read message from connection %v: %v", conn.RemoteAddr(), err)
			return
		}

		mutex.Lock()
		var wg sync.WaitGroup
		wg.Add(len(subscribers[topic]))

		for sub := range subscribers[topic] {
			go func(ws *websocket.Conn) {
				defer wg.Done()
				if err := sub.WriteMessage(messageType, message); err != nil {
					log.Printf("connection %v failed to write message: %v", conn.RemoteAddr(), err)
				}
			}(sub)
		}

		wg.Wait()
		mutex.Unlock()
	}
}

func main() {
	for _, topic := range topics {
		http.HandleFunc("/ws/"+topic, handleConnection)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
