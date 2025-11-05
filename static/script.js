// Configuración del canvas
const canvas = document.getElementById("board");
const ctx = canvas.getContext("2d");

// Ajustar tamaño del canvas
function resizeCanvas() {
    const container = canvas.parentElement;
    canvas.width = container.clientWidth - 40;
    canvas.height = container.clientHeight - 80;
    
    // Dibujar fondo blanco
    ctx.fillStyle = 'white';
    ctx.fillRect(0, 0, canvas.width, canvas.height);
}

// Inicializar canvas
resizeCanvas();
window.addEventListener('resize', resizeCanvas);

// Variables de dibujo
let drawing = false;
let lastX = 0;
let lastY = 0;
let currentColor = '#e74c3c';
let currentTool = 'brush';
let lineWidth = 3;

// Connect to WebSocket
const socket = new WebSocket(`ws://${location.host}/ws`);

// Handle incoming strokes
socket.onopen = () => {
    console.log("Conectado al servidor WebSocket");
    addSystemMessage("Conectado al servidor");
};

socket.onclose = () => {
    console.log("Desconectado del servidor WebSocket");
    addSystemMessage("Desconectado del servidor");
};

socket.onerror = (error) => {
    console.error("Error en WebSocket:", error);
    addSystemMessage("Error de conexión");
};

socket.onmessage = (event) => {
    try {
        const data = JSON.parse(event.data);

        // Si es un mensaje de chat
        if (data.type === 'chat') {
            console.log(data.sender)
            console.log(data.message)
            addMessage(data.sender, data.message, false);
            return;
        }

        // Si es un mensaje de sistema
        if (data.type === 'system') {
            addSystemMessage(data.message);
            updateOnlineUsers(data.onlineUsers || 1);
            return;
        }

        // Si es un trazo de dibujo
        const w = canvas.width;
        const h = canvas.height;
        drawLine(
            data.x0 * w, 
            data.y0 * h, 
            data.x1 * w, 
            data.y1 * h, 
            data.color || '#333',
            data.lineWidth || 3,
            false
        );
    } catch (error) {
        console.error("Error procesando mensaje:", error);
    }
};

// Eventos de dibujo
canvas.addEventListener("mousedown", (e) => {
    drawing = true;
    [lastX, lastY] = getMousePos(e);
});

canvas.addEventListener("mouseup", () => (drawing = false));
canvas.addEventListener("mouseout", () => (drawing = false));

canvas.addEventListener("mousemove", (e) => {
    if (!drawing) return;
    const [x, y] = getMousePos(e);
    
    if (currentTool === 'brush') {
        drawLine(lastX, lastY, x, y, currentColor, lineWidth, true);
    } else if (currentTool === 'eraser') {
        drawLine(lastX, lastY, x, y, 'white', lineWidth * 3, true);
    }
    
    [lastX, lastY] = [x, y];
});

// Para dispositivos táctiles
canvas.addEventListener('touchstart', (e) => {
    e.preventDefault();
    drawing = true;
    [lastX, lastY] = getMousePos(e.touches[0]);
});

canvas.addEventListener('touchmove', (e) => {
    e.preventDefault();
    if (!drawing) return;
    const [x, y] = getMousePos(e.touches[0]);
    
    if (currentTool === 'brush') {
        drawLine(lastX, lastY, x, y, currentColor, lineWidth, true);
    } else if (currentTool === 'eraser') {
        drawLine(lastX, lastY, x, y, 'white', lineWidth * 3, true);
    }
    
    [lastX, lastY] = [x, y];
});

canvas.addEventListener('touchend', () => (drawing = false));

// Función auxiliar para obtener posición del mouse/touch
function getMousePos(e) {
    const rect = canvas.getBoundingClientRect();
    const scaleX = canvas.width / rect.width;
    const scaleY = canvas.height / rect.height;
    
    return [
        (e.clientX - rect.left) * scaleX,
        (e.clientY - rect.top) * scaleY
    ];
}

// Draw line & send to server
function drawLine(x0, y0, x1, y1, color, width, emit) {
    ctx.beginPath();
    ctx.moveTo(x0, y0);
    ctx.lineTo(x1, y1);
    ctx.strokeStyle = color;
    ctx.lineWidth = width;
    ctx.lineCap = "round";
    ctx.lineJoin = "round";
    ctx.stroke();

    if (!emit) return;

    const w = canvas.width;
    const h = canvas.height;

    const data = {
        x0: x0 / w,
        y0: y0 / h,
        x1: x1 / w,
        y1: y1 / h,
        color: color,
        lineWidth: width
    };

    if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(data));
    }
}

// Selectores de color
const colorOptions = document.querySelectorAll('.color-option');
colorOptions.forEach(option => {
    option.addEventListener('click', () => {
        colorOptions.forEach(opt => opt.classList.remove('active'));
        option.classList.add('active');
        currentColor = option.getAttribute('data-color');
    });
});

// Herramientas
const toolButtons = document.querySelectorAll('.tool-btn');
toolButtons.forEach(button => {
    button.addEventListener('click', () => {
        toolButtons.forEach(btn => btn.classList.remove('active'));
        button.classList.add('active');
        currentTool = button.getAttribute('data-tool');
    });
});

// Botón limpiar
document.getElementById('clearBtn').addEventListener('click', () => {
    if (confirm('¿Estás seguro de que quieres limpiar el lienzo?')) {
        ctx.fillStyle = 'white';
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        
        // Opcional: enviar comando de limpiar a otros clientes
        if (socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({ type: 'clear' }));
        }
    }
});

// Botón guardar
document.getElementById('saveBtn').addEventListener('click', () => {
    const link = document.createElement('a');
    link.download = 'mi-dibujo.png';
    link.href = canvas.toDataURL('image/png');
    link.click();
});

// Funcionalidad del chat
const chatBox = document.getElementById('chatBox');
const input = document.getElementById('input');
const sendButton = document.getElementById('send');

function addMessage(sender, content, isOwn = false) {
    const messageDiv = document.createElement('div');
    messageDiv.className = isOwn ? 'message own' : 'message other';
    
    const messageHeader = document.createElement('div');
    messageHeader.className = 'message-header';
    
    const senderSpan = document.createElement('span');
    senderSpan.className = 'message-sender';
    senderSpan.textContent = isOwn ? 'Tú' : sender;
    
    const timeSpan = document.createElement('span');
    timeSpan.className = 'message-time';
    const now = new Date();
    timeSpan.textContent = `${now.getHours()}:${now.getMinutes().toString().padStart(2, '0')}`;
    
    messageHeader.appendChild(senderSpan);
    messageHeader.appendChild(timeSpan);
    
    const contentDiv = document.createElement('div');
    contentDiv.className = 'message-content';
    contentDiv.textContent = content;
    
    messageDiv.appendChild(messageHeader);
    messageDiv.appendChild(contentDiv);
    
    chatBox.appendChild(messageDiv);
    chatBox.scrollTop = chatBox.scrollHeight;
}

function addSystemMessage(content) {
    const messageDiv = document.createElement('div');
    messageDiv.className = 'message system';
    messageDiv.textContent = content;
    chatBox.appendChild(messageDiv);
    chatBox.scrollTop = chatBox.scrollHeight;
}

function updateOnlineUsers(count) {
    const onlineElement = document.querySelector('.online-users span');
    if (onlineElement) {
        onlineElement.textContent = `${count} usuarios en línea`;
    }
}

function sendMessage() {
    const message = input.value.trim();
    if (message && socket.readyState === WebSocket.OPEN) {
        // Enviar mensaje al servidor
        socket.send(JSON.stringify({
            type: 'chat',
            message: message
        }));
        
        // Mostrar mensaje localmente inmediatamente
        addMessage('Usuario', message, true);
        input.value = '';
    }
}

sendButton.addEventListener('click', sendMessage);
input.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        sendMessage();
    }
});