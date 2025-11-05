const canvas = document.getElementById("board");
const ctx = canvas.getContext("2d");

canvas.width = 800;
canvas.height = 500;

ctx.lineWidth = 3;
ctx.lineCap = "round";
ctx.strokeStyle = "#333";

// Connect to WebSocket
const socket = new WebSocket(`ws://${location.host}/ws`);

// Handle incoming strokes
socket.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log("Recibiendo...")
  drawLine(data.x0, data.y0, data.x1, data.y1, false);
};

// Drawing state
let drawing = false;
let lastX = 0;
let lastY = 0;

canvas.addEventListener("mousedown", (e) => {
  drawing = true;
  [lastX, lastY] = [e.offsetX, e.offsetY];
});

canvas.addEventListener("mouseup", () => (drawing = false));
canvas.addEventListener("mouseout", () => (drawing = false));

canvas.addEventListener("mousemove", (e) => {
  if (!drawing) return;
  const [x, y] = [e.offsetX, e.offsetY];
  drawLine(lastX, lastY, x, y, true);
  [lastX, lastY] = [x, y];
});

// Draw line & send to server
function drawLine(x0, y0, x1, y1, emit) {
  ctx.beginPath();
  ctx.moveTo(x0, y0);
  ctx.lineTo(x1, y1);
  ctx.stroke();

  if (!emit) return;

  const w = canvas.width;
  const h = canvas.height;

  const data = {
    x0: x0 / w,
    y0: y0 / h,
    x1: x1 / w,
    y1: y1 / h,
  };

  console.log("enviando...")
  socket.send(JSON.stringify(data));
}