# 🎨 Multiplayer Canvas

**Multiplayer Canvas** is a collaborative web application that allows multiple users to draw and chat in real time on a shared canvas.  
The backend is built with **Go (Golang)** using **WebSockets**, and the frontend is made with **HTML, CSS, and pure JavaScript**.

---

## 🚀 Main Features

- 🎨 **Real-time shared canvas**: every stroke is instantly visible to all users.  
- 💬 **Integrated chat**: communicate while drawing.  
- 🧹 **Basic tools**: brush, eraser, and line.  
- 🌈 **Color picker** and adjustable stroke thickness.  
- 💾 **Save your drawing** as a PNG image.  
- 👥 **Connected users indicator**.  
- 📱 **Responsive design** for mobile and small screens.

---

## 🧩 Project Structure

```
📦 multiplayer-canvas
├── main.go                # Main Go server (handles WebSockets and static files)
├── static/
│   ├── index.html         # Main UI
│   ├── script.js          # Drawing and chat logic
│   └── style.css          # Frontend styles
└── README.md              # This file
```

---

## ⚙️ Requirements

Before running the project, make sure you have installed:

- [Go 1.20+](https://go.dev/dl/)
- Modern web browser (Chrome, Firefox, Edge)
- Optional: [Render](https://render.com/), [Railway](https://railway.app/), or any Go-compatible hosting service

---

## 🧠 Technical Overview

### Backend (Go)
The server:
- Serves static files from the `./static` directory.
- Manages WebSocket connections (`/ws` endpoint).
- Keeps a list of active clients in memory.
- Broadcasts each message or stroke to all connected clients.

Main flow:
1. Client connects via WebSocket.  
2. Each received message (chat or stroke) is broadcast to everyone.  
3. When a client disconnects, it’s removed from the user registry.

Key file: [`main.go`](./main.go)
```go
http.Handle("/", http.FileServer(http.Dir("./static")))
http.HandleFunc("/ws", handleConnections)
go handleMessages()
```

---

### Frontend (HTML + JS + CSS)
- `index.html` defines the layout with canvas and chat.  
- `script.js` manages mouse/touch events, WebSocket communication, and synchronization.  
- `style.css` provides a clean dark theme with color accents.

---

## 💻 Run Locally

1. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/multiplayer-canvas.git
   cd multiplayer-canvas
   ```

2. Run the Go server:
   ```bash
   go run main.go
   ```

3. Open your browser at:
   ```
   http://localhost:8080
   ```

4. Open multiple tabs or devices and start drawing 🎨

---

## 🌐 Deploy on Render / Railway

### Render
1. Create a new web service on [Render](https://render.com/).  
2. Upload your repo or connect it to GitHub.  
3. Build command:
   ```bash
   go build -o app .
   ```
4. Start command:
   ```bash
   ./app
   ```
5. Render automatically assigns a `PORT` — the server handles it with:
   ```go
   port := os.Getenv("PORT")
   if port == "" {
       port = "8080"
   }
   ```

---

## 🧰 Technologies Used

| Technology | Purpose |
|-------------|----------|
| **Go (Golang)** | Backend and WebSocket server |
| **Gorilla WebSocket** | Real-time connection handling |
| **HTML5 Canvas** | Collaborative drawing on the web |
| **CSS3** | Visual design and responsiveness |
| **JavaScript (ES6)** | Drawing logic, chat, and synchronization |

---

## 🧑‍💻 Author

**Edgar Joel Villela Castañeda**  
📧 [edjovilellaca@ittepic.edu.mx]  
💻 Academic and demo project showcasing real-time collaboration using Go and WebSockets.

---

## 💡 Future Improvements

- 🗂️ Canvas persistence (save state on server or DB).  
- 🔐 User authentication.  
- 🖼️ Shared drawing gallery.  
- 💬 Chat history.  
- ✏️ Extra tools (shapes, fill, text, variable thickness).

---

## 📄 License

This project is released under the **MIT License**, meaning you can freely use, modify, and share it with proper attribution.
