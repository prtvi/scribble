const msgInp = document.querySelector('#msg');
const sendMsgBtn = document.querySelector('#sendMsg');
const messagesDiv = document.querySelector('.messages');
const membersDiv = document.querySelector('.members');

const canvas = document.querySelector('#canv');
const ctx = canvas.getContext('2d');

// -------- main

const coord = { x: 0, y: 0 };
let paint = false;
let color = '';

const [poolId, clientName, clientId] = initCredentials();
const wsUrl = `ws://${getDomain()}/ws?poolId=${poolId}&clientId=${clientId}&clientName=${clientName}`;

// establish socket connection
const socket = new WebSocket(wsUrl);
socket.onopen = () => console.log('Socket successfully connected');
socket.onmessage = socketOnMessage;
socket.onclose = socketOnClose;
socket.onerror = error => console.log('Socket error', error);

sendMsgBtn.addEventListener('click', sendMsgBtnEL);
window.addEventListener('load', windowEL);

const renderAllClientsTimer = setInterval(renderAllClients, 5000);

// -------- main

function initCredentials() {
	return [
		document.getElementsByName('poolId')[0].value,
		document.getElementsByName('clientName')[0].value,
		String(Date.now()),
	];
}

function getDomain() {
	// extract domain from url
	const url = window.location.href;
	const fi = url.indexOf('/');
	const li = url.lastIndexOf('/');
	const domain = url.slice(fi + 2, li);

	return domain;
}

function wait(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

async function initColor() {
	const allClients = await getAllClients();
	matchedClient = allClients.find(c => c.id === clientId);
	color = matchedClient.color;
}

async function getAllClients() {
	const res = await fetch(`api/get-clients-in-pool?poolId=${poolId}`);
	const data = await res.json();
	return data;
}

async function renderAllClients() {
	try {
		const allClients = await getAllClients();

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
		console.log('error, closing render all clients timer');
		clearInterval(renderAllClientsTimer);
	}
}

function addMsgToDOM(msg) {
	// adds the content into the DOM
	if (msg.length === 0 || msg === '') return;

	const msgDiv = document.createElement('div');
	const text = document.createTextNode(msg);
	msgDiv.appendChild(text);
	messagesDiv.appendChild(msgDiv);

	msgInp.value = '';
}

function displayImgOnCanvas(imgData) {
	// display image data on canvas
	var img = new Image();
	img.onload = () => ctx.drawImage(img, 0, 0);
	img.setAttribute('src', imgData);
}

function sendMessage(msg) {
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

async function sendImgData() {
	const respBody = {
		type: 4,
		content: String(canvas.toDataURL('img/png')),
		clientName,
		clientId,
	};

	// sending canvas data
	await wait(500);
	socket.send(JSON.stringify(respBody));
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
				addMsgToDOM(`You joined the pool as ${clientName}!`);
			else addMsgToDOM(`${msg.clientName} has joined the pool!`);
			break;

		case 2:
			addMsgToDOM(`${msg.clientName} has left the pool!`);
			break;

		case 3:
			addMsgToDOM(`${msg.clientName}: ${msg.content}`);
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
	clearInterval(renderAllClientsTimer);
}

function getPosition(event) {
	coord.x = event.clientX - canvas.offsetLeft;
	coord.y = event.clientY - canvas.offsetTop;
}

function startPainting(event) {
	paint = true;
	getPosition(event);
}

function stopPainting() {
	paint = false;
}

function sketch(event) {
	if (!paint) return;

	ctx.beginPath();

	ctx.lineWidth = 5;
	ctx.lineCap = 'round';
	ctx.strokeStyle = color;

	ctx.moveTo(coord.x, coord.y);

	getPosition(event);

	ctx.lineTo(coord.x, coord.y);
	ctx.stroke();

	sendImgData();
}

function sendMsgBtnEL(e) {
	e.preventDefault();
	sendMessage(msgInp.value);
}

async function windowEL() {
	await initColor();

	document.addEventListener('mousedown', startPainting);
	document.addEventListener('mouseup', stopPainting);
	document.addEventListener('mousemove', sketch);
}
