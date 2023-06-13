'use strict';

function initCanvasAndOverlay() {
	const canvas = document.querySelector('.canv');
	const ctx = canvas.getContext('2d');

	const { w, h } = getCanvasSize();
	canvas.width = w;
	canvas.height = h;

	const cc = document.querySelector('.canvas-container');
	cc.style.width = `${w}px`;
	cc.style.height = `${h}px`;

	const overlay = document.querySelector('#overlay');
	overlay.style.top = `${cc.offsetTop}px`;
	overlay.style.height = `${cc.offsetHeight}px`;

	return { canvas, ctx, overlay };
}

function getCanvasSize() {
	const w = window.innerWidth;
	const cw = w;
	const ch = cw / 1.5;

	return { w: cw, h: ch };
}

// drawing on canvas

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

// render canvas data

// clear canvas

function requestCanvasClear() {
	// clear canvas and request clear on rest of the clients
	ctx.clearRect(0, 0, canvas.width, canvas.height);

	// broadcast clear canvas
	const socketMsg = {
		type: 5,
		typeStr: messageTypeMap.get(5),
		clientId,
		clientName,
		poolId,
	};

	sendViaSocket(socketMsg);
}

// send image data

function sendImgData() {
	// called by paint function
	const socketMsg = {
		type: 4,
		typeStr: messageTypeMap.get(4),
		content: String(canvas.toDataURL('img/png')),
		clientName,
		clientId,
		poolId,
	};

	// sending canvas data
	sendViaSocket(socketMsg);
}
