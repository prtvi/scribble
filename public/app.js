'use strict';

// canvas
const canvas = document.querySelector('.canv');
const ctx = canvas.getContext('2d');

// chat messages
const msgInp = document.querySelector('.msg');
const sendChatMsgBtn = document.querySelector('.send-msg');

// ---------------- main ----------------

// utils for paintting on canvas
const paintUtils = {
	coords: { x: 0, y: 0 },
	color: `#${clientColor}`,
	isPainting: false,
	hasGameStarted: false,
	isAllowedToPaint: false, // not used yet
};

// timers to start when game has not yet started
let startGameTimer,
	startGameAfterInterval,
	secondsLeft = gameStartsInSeconds;

// init socket connection and check game begin status
const socket = initSocket();
checkGameBeginStat();

// render all clients in pool on UI every 5 seconds
const renderClientsTimer = setInterval(getAllClientsEL, 5 * 1000);
// event listeners to send chat messages event listeners for canvas painting
sendChatMsgBtn.addEventListener('click', sendChatMsgBtnEL);
window.addEventListener('load', addCanvasEventListeners);

// ------------------------------------- utils -------------------------------------

function initSocket() {
	// initialises socket connection and adds corresponding function handlers to the socket

	const wsUrl = `ws://${getDomain()}/ws?poolId=${poolId}&clientId=${clientId}&clientName=${clientName}&clientColor=${clientColor}`;
	const socket = new WebSocket(wsUrl);

	socket.onopen = () => {
		// on socket open success, get all clients and render them on UI
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
		// runs when a message is received on the socket conn, runs the corresponding functions depending on message type

		// parse json string into json object
		const msg = JSON.parse(message.data);

		// msg.type
		// 1 === CONNECTED
		// 2 === DISCONNECTED
		// 3 === string data
		// 4 === canvas data
		// 5 === all client info
		// 6 === start game ack

		switch (msg.type) {
			case 1:
				if (msg.clientId === clientId)
					// if the current clientId and the clientId from response match then
					appendChatMsgToDOM(`You joined the pool as ${msg.clientName}!`);
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

			case 6:
				startGame(msg.content);
				break;

			default:
				break;
		}
	}

	function socketOnClose() {
		// on socket conn close, stop all timer or intervals
		console.log('Socket connection closed, stopping timers and timeouts!');
		clearInterval(renderClientsTimer);
		clearInterval(startGameTimer);
		clearTimeout(startGameAfterInterval);
	}

	return socket;
}

function checkGameBeginStat() {
	// checks if game has already started based on "hasGameStarted" variable

	// if game has not started then, begin the countdown and render time left
	// and request start game on time-up
	// also listen to start game btn press to start the game
	if (!hasGameStarted) {
		console.log('game not started');

		// start game countdown to show user how much time is left
		startGameTimer = setInterval(renderCountdownEL, 1000);

		// start game after this timeout
		startGameAfterInterval = setTimeout(
			requestStartGameEL,
			(gameStartsInSeconds + 2) * 1000
		);

		// add event listener to start game button to start game
		document
			.querySelector('.start-game-btn')
			.addEventListener('click', requestStartGameEL);
	} else {
		// if game has already begun, then alter the hasGameStarted field in paintUtils
		console.log('started game');
		paintUtils.hasGameStarted = true;
	}
}

function wait(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

// ------------------------------------- start game countdown -------------------------------------

function renderCountdownEL() {
	// renders the countdown to display to the user remaining time for game to start
	document.querySelector('.loading').textContent = secondsLeft;
	secondsLeft -= 1;
}

function requestStartGameEL() {
	// runs when the game starts, makes socket conn call to server to start the game

	// clear the countdown timer to show time left to start game
	clearInterval(startGameTimer);

	// generate response and send
	const responseMsg = {
		type: 6,
		content: 'start the game bro!',
		poolId,
	};

	socket.send(JSON.stringify(responseMsg));
}

function startGame(msg) {
	// called when socket receives message from server with type as 6
	if (msg !== 'true') return;

	console.log('started game ...');

	// clear interval after response from server
	clearTimeout(startGameAfterInterval);

	// hide the div and toggle hasGameStarted
	document.querySelector('.start-game').style.display = 'none';
	paintUtils.hasGameStarted = true;
}

// ------------------------------------- get all clients and render -------------------------------------

function getAllClientsEL() {
	// makes a socket connection call to request client info list
	const responseMsg = {
		type: 5,
		content: '',
		poolId,
	};

	socket.send(JSON.stringify(responseMsg));
}

function renderClients(allClients) {
	// called when the socket conn receives a message from server as type 5

	const membersDiv = document.querySelector('.members');
	membersDiv.innerHTML = '';

	// parse array of objects into json
	allClients = JSON.parse(allClients);

	// render
	allClients.forEach(n => {
		const clientNameHolder = document.createElement('div');
		const clientName = document.createElement('p');

		clientName.innerHTML = n.name;
		clientName.style.color = `#${n.color}`;
		clientNameHolder.appendChild(clientName);

		membersDiv.appendChild(clientNameHolder);
	});
}

// ------------------------------------- chat -------------------------------------

function appendChatMsgToDOM(msg) {
	// adds the msg into the DOM

	if (msg.length === 0 || msg === '') return;

	const messagesDiv = document.querySelector('.messages');
	const msgDiv = document.createElement('div');
	const text = document.createTextNode(msg);
	msgDiv.appendChild(text);
	messagesDiv.appendChild(msgDiv);

	msgInp.value = '';
}

function sendChatMsgBtnEL(e) {
	// event listener to send chat message

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

// ------------------------------------- canvas -------------------------------------

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
