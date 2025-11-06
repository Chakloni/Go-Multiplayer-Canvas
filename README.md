# 🎨 Multiplayer Canvas

**Multiplayer Canvas** es una aplicación web colaborativa que permite a múltiples usuarios dibujar y chatear en tiempo real sobre un mismo lienzo compartido.  
El backend está desarrollado en **Go (Golang)** utilizando **WebSockets**, mientras que el frontend está hecho en **HTML, CSS y JavaScript puro**.

---

## 🚀 Características principales

- 🎨 **Lienzo compartido en tiempo real**: todos los usuarios ven los trazos al instante.  
- 💬 **Chat integrado**: permite conversar mientras se dibuja.  
- 🧹 **Herramientas básicas**: pincel, borrador y línea.  
- 🌈 **Selector de color** y grosor del trazo.  
- 💾 **Guardar dibujo** como imagen PNG.  
- 👥 **Indicador de usuarios conectados**.  
- 📱 **Diseño adaptable (responsive)** para móviles y pantallas pequeñas.

---

## 🧩 Estructura del proyecto

```
📦 multiplayer-canvas
├── main.go                # Servidor principal en Go (maneja WebSockets y archivos estáticos)
├── static/
│   ├── index.html         # Interfaz principal
│   ├── script.js          # Lógica de dibujo y chat
│   └── style.css          # Estilos del frontend
└── README.md              # Este archivo
```

---

## ⚙️ Requisitos previos

Antes de ejecutar el proyecto, asegúrate de tener instalado:

- [Go 1.20+](https://go.dev/dl/)
- Navegador web moderno (Chrome, Firefox, Edge)
- Opcional: [Render](https://render.com/), [Railway](https://railway.app/), o cualquier servicio compatible con Go para desplegarlo.

---

## 🧠 Funcionamiento técnico

### Backend (Go)
El servidor:
- Sirve archivos estáticos desde la carpeta `./static`.
- Administra las conexiones WebSocket (`/ws`).
- Mantiene una lista de clientes activos en memoria.
- Difunde cada mensaje o trazo a todos los clientes conectados.

Flujo principal:
1. El cliente se conecta vía WebSocket.
2. Cada mensaje recibido (ya sea de chat o trazo) se retransmite a todos.
3. Cuando un cliente se desconecta, se elimina del registro de usuarios.

Archivo clave: [`main.go`](./main.go)
```go
http.Handle("/", http.FileServer(http.Dir("./static")))
http.HandleFunc("/ws", handleConnections)
go handleMessages()
```

---

### Frontend (HTML + JS + CSS)
- `index.html` define la estructura de la interfaz con el lienzo y el chat.
- `script.js` gestiona los eventos del mouse/táctiles, la conexión WebSocket y la sincronización del dibujo/chat.
- `style.css` proporciona una estética moderna con un tema oscuro y acentos de color.

---

## 💻 Ejecución local

1. Clona este repositorio:
   ```bash
   git clone https://github.com/tuusuario/multiplayer-canvas.git
   cd multiplayer-canvas
   ```

2. Ejecuta el servidor Go:
   ```bash
   go run .
   ```

3. Abre tu navegador en:
   ```
   http://localhost:8080
   ```

4. Conecta múltiples pestañas o dispositivos y empieza a dibujar 🎨

---

## 🌐 Despliegue en Render / Railway

### Render
1. Crea un nuevo servicio web en [Render](https://render.com/).
2. Sube tu repositorio o conéctalo a GitHub.
3. Usa este comando de build:
   ```bash
   go build -o app .
   ```
4. Comando de ejecución:
   ```bash
   ./app
   ```
5. Render asignará un `PORT` automáticamente — el servidor ya lo detecta con:
   ```go
   port := os.Getenv("PORT")
   if port == "" {
       port = "8080"
   }
   ```

---

## 🧰 Tecnologías utilizadas

| Tecnología | Uso |
|-------------|-----|
| **Go (Golang)** | Backend y servidor WebSocket |
| **Gorilla WebSocket** | Manejo de conexiones en tiempo real |
| **HTML5 Canvas** | Dibujo compartido en el navegador |
| **CSS3** | Diseño visual y responsividad |
| **JavaScript (ES6)** | Lógica del chat, dibujo y sincronización |

---

## 🧑‍💻 Autor

**Edgar Joel Villela Castañeda**  
📧 [edjovilellaca@ittepic.edu.mx]  
💻 Proyecto académico y demostrativo de colaboración en tiempo real con Go y WebSockets.

---

## 🪄 Posibles mejoras futuras

- 🗂️ Persistencia del lienzo (guardar estado en servidor o BD).  
- 🔐 Autenticación de usuarios.  
- 🖼️ Galería de dibujos compartidos.  
- 💬 Chat con historial.  
- ✏️ Herramientas adicionales (formas, relleno, texto, grosor variable).

---

## 📄 Licencia

Este proyecto se distribuye bajo la licencia **MIT**, por lo que puedes usarlo, modificarlo y compartirlo libremente con atribución.
