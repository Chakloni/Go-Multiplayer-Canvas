package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { 
		return true 
	},
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

var clients = make(map[*Client]bool)
var broadcast = make(chan []byte)
var mu sync.Mutex

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("🎨 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	client := &Client{conn: ws, send: make(chan []byte)}
	mu.Lock()
	clients[client] = true
	mu.Unlock()

	log.Println("🟢 New client connected 🙋")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("🔴 Client disconnected 👋:", err)
			mu.Lock()
			delete(clients, client)
			mu.Unlock()
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			select {
			case client.send <- msg:
				go func(c *Client, m []byte) {
					c.conn.WriteMessage(websocket.TextMessage, m)
				}(client, msg)
			default:
				close(client.send)
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}