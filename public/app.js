'use strict';
// clientId, clientName and poolId initialised from inline js

const canvas = document.querySelector('#canv');
const ctx = canvas.getContext('2d');

const msgInp = document.querySelector('#msg');
const sendChatMsgBtn = document.querySelector('#send-msg');

const loading = document.querySelector('.loading');
const startGameDiv = document.querySelector('.start-game');
const startGameBtn = document.querySelector('.start-game-btn');
const now = new Date();

// ---------------- main ----------------

const paintUtils = {
	coords: { x: 0, y: 0 },
	color: `#${clientColor}`,
	isPainting: false,
	hasGameStarted: false,
	isAllowedToPaint: false, // not used yet
};

const socket = initSocket();

const renderClientsTimer = setInterval(getAllClientsEL, 5 * 1000);
const startGameCountdownTimer = setInterval(startGameCountdownEL, 1000);
const startGameAfterTimeout = setTimeout(
	startGameAfterTimeoutEL,
	Math.ceil(gameStartTime - now)
);

sendChatMsgBtn.addEventListener('click', sendChatMsgBtnEL);
window.addEventListener('load', addCanvasEventListeners);
startGameBtn.addEventListener('click', startGameAfterTimeoutEL);

// ---------------- chat ----------------

function appendChatMsgToDOM(msg) {
	// adds the content into the DOM
	if (msg.length === 0 || msg === '') return;

	const messagesDiv = document.querySelector('.messages');

	const msgDiv = document.createElement('div');
	const text = document.createTextNode(msg);
	msgDiv.appendChild(text);
	messagesDiv.appendChild(msgDiv);

	msgInp.value = '';
}

function sendChatMsgBtnEL(e) {
	e.preventDefault();
	const msg = msgInp.value;

	if (msg.length === 0 || msg === '') return;

	// create string response object
	const responseMsg = {
		type: 3,
		content: msg,
		clientName,
		clientId,
	};

	// convert object to string to transmit
	socket.send(JSON.stringify(responseMsg));
}

// ---------------- canvas ----------------

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

async function sendImgData() {
	// called by paint function
	const respBody = {
		type: 4,
		content: String(canvas.toDataURL('img/png')),
		clientName,
		clientId,
	};

	// sending canvas data
	socket.send(JSON.stringify(respBody));
}

function addCanvasEventListeners() {
	document.addEventListener('mousedown', startPainting);
	document.addEventListener('mouseup', stopPainting);
	document.addEventListener('mousemove', paint);
}

// ---------------- get all clients list and render ----------------

function renderClients(allClients) {
	const membersDiv = document.querySelector('.members');
	membersDiv.innerHTML = '';

	allClients = JSON.parse(allClients);

	allClients.forEach(n => {
		const clientNameHolder = document.createElement('div');
		const clientName = document.createElement('p');

		clientName.innerHTML = n.name;
		clientName.style.color = `#${n.color}`;
		clientNameHolder.appendChild(clientName);

		membersDiv.appendChild(clientNameHolder);
	});
}

function getAllClientsEL() {
	const responseMsg = {
		type: 5,
		content: '',
		poolId,
	};

	socket.send(JSON.stringify(responseMsg));
}

// ---------------- start game countdown ----------------

function startGameCountdownEL() {
	// does the countdown to display to the user remaining time for game to start

	const now = new Date();
	loading.textContent = Math.ceil((gameStartTime - now) / 1000);
}

async function startGameAfterTimeoutEL() {
	// runs when the game starts

	const res = await fetch(`/api/start-game?poolId=${poolId}`);
	const data = await res.json();

	if (!data.success) return;

	console.log('starting game ...');
	clearInterval(startGameCountdownTimer);
	clearTimeout(startGameAfterTimeout);

	startGameDiv.style.display = 'none';
	paintUtils.hasGameStarted = true;
}

// ---------------- utils ----------------

function initSocket() {
	const wsUrl = `ws://${getDomain()}/ws?poolId=${poolId}&clientId=${clientId}&clientName=${clientName}&clientColor=${clientColor}`;

	const socket = new WebSocket(wsUrl);

	socket.onopen = () => {
		console.log('Socket successfully connected!');
		getAllClientsEL();
	};

	socket.onmessage = socketOnMessage;
	socket.onclose = socketOnClose;
	socket.onerror = error => console.log('Socket error', error);

	function getDomain() {
		// extract domain from url
		const url = window.location.href;
		const fi = url.indexOf('/');
		const li = url.lastIndexOf('/');
		const domain = url.slice(fi + 2, li);

		return domain;
	}

	function socketOnMessage(message) {
		// parse json string into json object
		const msg = JSON.parse(message.data);

		// if message type is 1 === CONNECTED
		// if message type is 2 === DISCONNECTED
		// if message type is 3 === string data
		// if message type is 4 === canvas data
		// if message type is 5 === all client info

		switch (msg.type) {
			case 1:
				// if the current clientName and the clientName from response match then
				if (msg.clientName === clientName)
					appendChatMsgToDOM(`You joined the pool as ${clientName}!`);
				else appendChatMsgToDOM(`${msg.clientName} has joined the pool!`);
				break;

			case 2:
				appendChatMsgToDOM(`${msg.clientName} has left the pool!`);
				break;

			case 3:
				appendChatMsgToDOM(`${msg.clientName}: ${msg.content}`);
				break;

			case 4:
				displayImgOnCanvas(msg.content);
				break;

			case 5:
				renderClients(msg.content);
				break;

			default:
				break;
		}
	}

	function socketOnClose() {
		console.log('Socket connection closed, stopping timers and timeouts!');
		clearInterval(renderClientsTimer);
		clearInterval(startGameCountdownTimer);
		clearTimeout(startGameAfterTimeout);
	}

	return socket;
}

function wait(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}
