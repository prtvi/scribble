'use strict';

function initCanvasAndOverlay() {
	const canvas = document.querySelector('.canv');
	const ctx = canvas.getContext('2d');

	const { w, h } = getCanvasSize();

	canvas.width = w;
	canvas.height = h;

	const overlay = document.querySelector('#overlay');
	const cc = document.querySelector('.canvas-container');

	overlay.style.top = `${cc.offsetTop}px`;
	overlay.style.height = `${cc.offsetHeight}px`;

	return { canvas, ctx, overlay };
}

function getCanvasSize() {
	const w = window.innerWidth;
	const cw = w - 10;
	const ch = cw / 1.5;

	return { w: cw, h: ch };
}

// drawing

function updatePositionCanvas(event) {
	paintUtils.coords.x = event.clientX - canvas.offsetLeft;
	paintUtils.coords.y = event.clientY - canvas.offsetTop;
}

function startPainting(event) {
	paintUtils.isPainting = true;
	updatePositionCanvas(event);
}

function stopPainting() {
	paintUtils.isPainting = false;
}

async function paint(event) {
	if (!paintUtils.isPainting) return;
	if (!paintUtils.hasGameStarted) return;
	if (!paintUtils.isAllowedToPaint) return;

	ctx.beginPath();

	ctx.lineWidth = 5;
	ctx.lineCap = 'round';
	ctx.strokeStyle = paintUtils.color;

	ctx.moveTo(paintUtils.coords.x, paintUtils.coords.y);

	updatePositionCanvas(event);

	ctx.lineTo(paintUtils.coords.x, paintUtils.coords.y);
	ctx.stroke();

	await wait(500);
	sendImgData();
}

function displayImgOnCanvas(socketMessage) {
	// display image data on canvas

	if (clientId === socketMessage.currSketcherId) return;

	var img = new Image();
	// scale up/down canvas data based on current canvas size using outer bounds
	img.onload = () => ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
	img.setAttribute('src', socketMessage.content);
}

function requestCanvasClear() {
	// broadcast clear canvas
	const socketMsg = {
		type: 5,
		typeStr: 'clear_canvas',
		clientId,
		clientName,
		poolId,
	};

	sendViaSocket(socketMsg);
}

function clearCanvas() {
	hideOverlay();
	ctx.clearRect(0, 0, canvas.width, canvas.height);
}

function sendImgData() {
	// called by paint function
	const socketMsg = {
		type: 4,
		typeStr: 'canvas_data',
		content: String(canvas.toDataURL('img/png')),
		clientName,
		clientId,
		poolId,
	};

	// sending canvas data
	sendViaSocket(socketMsg);
}

window.addEventListener('load', () => {
	document.addEventListener('mousedown', startPainting);
	document.addEventListener('mouseup', stopPainting);
	document.addEventListener('mousemove', paint);
});
