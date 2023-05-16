'use strict';

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

function displayImgOnCanvas(imgData) {
	// display image data on canvas
	var img = new Image();
	img.onload = () => ctx.drawImage(img, 0, 0);
	img.setAttribute('src', imgData);
}

function requestCanvasClear() {
	// broadcast clear canvas
	const responseMsg = {
		type: 5,
		content: 'clear canvas',
		clientId,
		clientName,
		poolId,
	};

	sendViaSocket(responseMsg);
}

function clearCanvas() {
	ctx.clearRect(0, 0, canvas.width, canvas.height);
}

function sendImgData() {
	// called by paint function
	const responseMsg = {
		type: 4,
		content: String(canvas.toDataURL('img/png')),
		clientName,
		clientId,
		poolId,
	};

	// sending canvas data
	sendViaSocket(responseMsg);
}

window.addEventListener('load', () => {
	document.addEventListener('mousedown', startPainting);
	document.addEventListener('mouseup', stopPainting);
	document.addEventListener('mousemove', paint);
});
