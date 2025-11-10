package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client representa un cliente conectado
type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	username string
	hub      *Hub
    lastMessage time.Time
}

// Hub mantiene el conjunto de clientes activos y difunde mensajes
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// Message representa la estructura de mensajes
type Message struct {
	Type    string      `json:"type"`
	Content string      `json:"content,omitempty"`
	Message ChatMessage `json:"message,omitempty"`
	// Campos para dibujo
	X0        float64 `json:"x0,omitempty"`
	Y0        float64 `json:"y0,omitempty"`
	X1        float64 `json:"x1,omitempty"`
	Y1        float64 `json:"y1,omitempty"`
	Color     string  `json:"color,omitempty"`
	LineWidth float64 `json:"lineWidth,omitempty"`
}

type ChatMessage struct {
	Username string `json:"username"`
	Msg      string `json:"msg"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NewHub crea un nuevo Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run inicia el hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.broadcastUserCount()
			log.Printf("Cliente conectado. Total: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			h.broadcastUserCount()
			log.Printf("Cliente desconectado. Total: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Si el canal está lleno, desconectar el cliente
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) broadcastUserCount() {
	h.mu.RLock()
	count := len(h.clients)
	h.mu.RUnlock()

	msg := Message{
		Type:    "userCount",
		Content: fmt.Sprintf("%d", count),
	}
	
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling user count: %v", err)
		return
	}

	h.broadcast <- data
}

// readPump lee mensajes del websocket
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}

		// Validar que el mensaje no esté vacío
		if len(message) == 0 {
			continue
		}

		// Reenviar el mensaje a todos los clientes
		c.hub.broadcast <- message
	}
}

// writePump escribe mensajes al websocket
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Agregar mensajes en cola al mensaje actual
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleWebSocket maneja las conexiones WebSocket
func handleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
		hub:  hub,
	}

	client.hub.register <- client

	// Iniciar goroutines para leer y escribir
	go client.writePump()
	go client.readPump()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	hub := NewHub()
	go hub.Run()

	// Servir archivos estáticos
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Endpoint WebSocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})

	// Endpoint de salud (útil para monitoreo)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	log.Printf("🎨 Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}