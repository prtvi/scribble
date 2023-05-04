'use strict';
// clientId, clientName and poolId initialised from inline js

const canvas = document.querySelector('#canv');
const ctx = canvas.getContext('2d');

// ---------------- main ----------------

const paintUtils = {
	coords: { x: 0, y: 0 },
	color: '',
	isPainting: false,
};

const socket = initSocket();

const displayAllClientsInPoolTimer = setInterval(displayAllClientsInPool, 5000);

sendChatMsgBtn.addEventListener('click', sendChatMsgBtnEL);
window.addEventListener('load', initColor);
window.addEventListener('load', addCanvasEventListeners);

// ---------------- main ----------------

// ---------------- chat ----------------

const msgInp = document.querySelector('#msg');
const sendChatMsgBtn = document.querySelector('#sendMsg');

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

// ---------------- chat ----------------

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

// ---------------- canvas ----------------

// ---------------- get all clients list ----------------

async function getAllClients() {
	const res = await fetch(`api/get-clients-in-pool?poolId=${poolId}`);
	const data = await res.json();
	return data;
}

async function displayAllClientsInPool() {
	try {
		const allClients = await getAllClients();

		const membersDiv = document.querySelector('.members');
		membersDiv.innerHTML = '';

		allClients.forEach(n => {
			const clientNameHolder = document.createElement('div');
			const clientName = document.createElement('p');

			clientName.innerHTML = n.name;
			clientName.style.color = n.color;
			clientNameHolder.appendChild(clientName);

			membersDiv.appendChild(clientNameHolder);
		});
	} catch (error) {
		console.log('error, closing display all clients timer');
		clearInterval(displayAllClientsInPoolTimer);
	}
}

// ---------------- get all clients list ----------------

// ---------------- utils ----------------

function initSocket() {
	const wsUrl = `ws://${getDomain()}/ws?poolId=${poolId}&clientId=${clientId}&clientName=${clientName}`;

	const socket = new WebSocket(wsUrl);
	socket.onopen = () => console.log('Socket successfully connected');
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

			default:
				break;
		}
	}

	function socketOnClose() {
		console.log('Socket connection closed, stopping render all clients timer');
		clearInterval(displayAllClientsInPoolTimer);
	}

	return socket;
}

function wait(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

async function initColor() {
	const allClients = await getAllClients();
	const matchedClient = allClients.find(c => c.id === clientId);
	paintUtils.color = matchedClient.color;
}

// ---------------- utils ----------------
