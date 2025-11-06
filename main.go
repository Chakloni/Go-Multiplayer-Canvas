package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}
var userCount = 0 // 👈 contador global

type Message struct {
    Type    string `json:"type"`
    Content string `json:"content"`
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Error upgrading connection: %v", err)
        return
    }
    defer ws.Close()

    clients[ws] = true
    userCount++ // 👈 incrementa al conectar
    broadcastUserCount()

    for {
        var msg Message
        err := ws.ReadJSON(&msg)
        if err != nil {
            delete(clients, ws)
            userCount-- // 👈 decrementa al desconectar
            broadcastUserCount()
            break
        }
        broadcast <- msg
    }
}

func handleMessages() {
    for {
        msg := <-broadcast
        for client := range clients {
            err := client.WriteJSON(msg)
            if err != nil {
                log.Printf("Error sending message: %v", err)
                client.Close()
                delete(clients, client)
                userCount--
                broadcastUserCount()
            }
        }
    }
}

func broadcastUserCount() {
    msg := Message{
        Type:    "userCount",
        Content: fmt.Sprintf("%d users online", userCount),
    }
    for client := range clients {
        client.WriteJSON(msg)
    }
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)
    http.HandleFunc("/ws", handleConnections)

    go handleMessages()

    log.Printf("🎨 Server running on http://localhost:%s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
