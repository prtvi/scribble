'use strict';

// canvas
const canvas = document.querySelector('.canv');
const ctx = canvas.getContext('2d');

// chat messages
const msgInp = document.querySelector('.msg');
const sendChatMsgBtn = document.querySelector('.send-msg');

// ---------------- main ----------------
// utils for painting on canvas
const paintUtils = {
	coords: { x: 0, y: 0 },
	color: `#${clientColor}`,
	isPainting: false,
	hasGameStarted: false,
	isAllowedToPaint: false,
};

// timers
let startGameTimerId,
	startGameAfterIntervalId,
	wordExpiryTimerId,
	renderClientsTimerId;

// init socket connection and check game begin status
const socket = initSocket();
checkGameBeginStat();

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
		const socketMessage = JSON.parse(message.data);

		// socketMessage.type
		// 1 === CONNECTED
		// 2 === DISCONNECTED
		// 3 === string data
		// 4 === canvas data
		// 5 === all client info
		// 6 === start game

		switch (socketMessage.type) {
			case 1:
				if (socketMessage.clientId === clientId)
					// if the current clientId and the clientId from response match then
					appendChatMsgToDOM(
						`You joined the pool as ${socketMessage.clientName}!`
					);
				else
					appendChatMsgToDOM(
						`${socketMessage.clientName} has joined the pool!`
					);
				break;

			case 2:
				appendChatMsgToDOM(`${socketMessage.clientName} has left the pool!`);
				break;

			case 3:
				appendChatMsgToDOM(
					`${socketMessage.clientName}: ${socketMessage.content}`
				);
				break;

			case 4:
				displayImgOnCanvas(socketMessage.content);
				break;

			case 5:
				renderClients(socketMessage.content);
				break;

			case 6:
				startGame(socketMessage);
				break;

			default:
				break;
		}
	}

	function socketOnClose() {
		// on socket conn close, stop all timer or intervals
		console.log('Socket connection closed, stopping timers and timeouts!');
		clearInterval(renderClientsTimerId);
		clearInterval(startGameTimerId);
		clearTimeout(startGameAfterIntervalId);
		clearInterval(wordExpiryTimerId);
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

		function requestStartGameEL() {
			// runs when the game starts, makes socket conn call to server to start the game
			// clear the countdown timers
			clearInterval(startGameTimerId);
			clearTimeout(startGameAfterIntervalId);

			// generate response and send
			const responseMsg = {
				type: 6,
				content: 'start the game bro!',
				poolId,
			};

			socket.send(JSON.stringify(responseMsg));
		}

		// add event listener to start game button to start game
		const startGameBtn = document.querySelector('.start-game-btn');
		startGameBtn.addEventListener('click', requestStartGameEL);

		// start game countdown to show user how much time is left
		startGameTimerId = setInterval(
			() =>
				(document.querySelector('.loading').textContent =
					getSecondsLeftFrom(gameStartTime)),
			1000
		);

		// start game after this timeout
		startGameAfterIntervalId = setTimeout(
			requestStartGameEL,
			getSecondsLeftFrom(gameStartTime) * 1000
		);
	} else {
		// if game has already begun, then alter the hasGameStarted field in paintUtils
		console.log(
			'game has started already, but I have no idea about the status of the game!'
		);
		paintUtils.hasGameStarted = true;
	}
}

function wait(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

function getSecondsLeftFrom(futureTime) {
	const now = new Date().getTime();
	const diff = futureTime - now;
	return Math.round(diff / 1000);
}

// ------------------------ start game ------------------------

function beginClientSketchingFlow(socketMessage) {
	// initialise the time at which this word expires
	const currentWordExpiresAt = new Date(
		socketMessage.currWordExpiresAt
	).getTime();

	const timeLeftDiv = document.querySelector('.time-left-for-word');

	// start timer for the word expiry
	wordExpiryTimerId = setInterval(() => {
		timeLeftDiv.classList.remove('hidden');

		const secondsLeft = getSecondsLeftFrom(currentWordExpiresAt);
		if (secondsLeft <= 0) {
			clearInterval(wordExpiryTimerId);
			console.log('timer for word cleared');

			// trigger next word for next player: TODO
		}

		timeLeftDiv.querySelector('span').textContent = secondsLeft;
	}, 1000);

	// for enabling drawing access if clientId matches
	if (clientId === socketMessage.currSketcherId) {
		paintUtils.isAllowedToPaint = true;

		// display the word by unhiding the painter-utils div
		document.querySelector('.painter-utils').classList.remove('hidden');
		document.querySelector('.your-word').textContent = socketMessage.currWord;

		// add EL for clearing the canvas
		document.querySelector('.clear-canvas').addEventListener('click', () => {
			ctx.clearRect(0, 0, canvas.width, canvas.height);

			console.log('send canvas data after clearing not working');
			sendImgData();
		});
	}
}

function startGame(socketMessage) {
	// called when socket receives message from server with type as 6
	if (socketMessage.content !== 'true') return;

	console.log('game started by server');
	paintUtils.hasGameStarted = true;

	// hide the div and toggle hasGameStarted
	const startGameDiv = document.querySelector('.start-game');
	if (startGameDiv) startGameDiv.classList.add('hidden');

	beginClientSketchingFlow(socketMessage);
}

// ------------------------ get all clients and render ------------------------

// render all clients in pool on UI every n seconds
renderClientsTimerId = setInterval(getAllClientsEL, 10 * 1000);

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

		clientName.innerHTML = `${n.name} with ${n.score} points`;
		clientName.style.color = `#${n.color}`;
		clientNameHolder.appendChild(clientName);

		membersDiv.appendChild(clientNameHolder);
	});
}

// ------------------------------------- chat -------------------------------------

// event listeners to send chat messages
sendChatMsgBtn.addEventListener('click', sendChatMsgBtnEL);

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

// event listeners for canvas painting
window.addEventListener('load', addCanvasEventListeners);

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

function sendImgData() {
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
