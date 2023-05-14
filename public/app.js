'use strict';

// canvas
const canvas = document.querySelector('.canv');
const ctx = canvas.getContext('2d');

// ---------------- main ----------------
// utils for painting on canvas
const paintUtils = {
	coords: { x: 0, y: 0 },
	color: `#${clientColor}`,
	isPainting: false,
	hasGameStarted: false,
	isAllowedToPaint: false,
};

// init socket connection and check game begin status
const socket = initSocket();
checkGameBeginStat();

var wordExpiryTimerId, currentWordExpiresAt;

// render all clients in pool on UI every n seconds
const renderClientsTimerId = setInterval(getAllClientsEL, 10 * 1000);
window.addEventListener('load', addCanvasEventListeners);
document.querySelector('.send-msg').addEventListener('click', sendChatMsgBtnEL);

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
		// 5 === clear canvas
		// 6 === all client info
		// 7 === start game
		// 8 === request next word
		// 9 === finish game and display scores

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
				clearCanvas();
				break;

			case 6:
				renderClients(socketMessage.content);
				break;

			case 7:
				startGame(socketMessage);
				break;

			case 8:
				beginClientSketchingFlow(socketMessage);
				break;

			case 9:
				displayScores(socketMessage);
				break;

			default:
				break;
		}
	}

	function socketOnClose() {
		// on socket conn close, stop all timer or intervals
		console.log('Socket connection closed, stopping timers and timeouts!');
		clearAllIntervals(renderClientsTimerId);
	}

	return socket;
}

function sendViaSocket(responseMsg) {
	if (socket.readyState === socket.OPEN)
		socket.send(JSON.stringify(responseMsg));
	else
		console.log(
			'socket already closed | yet opening | in closing state',
			socket.readyState
		);
}

function checkGameBeginStat() {
	// checks if game has already started based on "has Game Started" variable

	if (hasGameStarted) return;

	// if game has not started then, begin the countdown and render time left
	// and request start game on time-up
	// also listen to start game btn press to start the game

	console.log('game not started');

	// add event listener to start game button to start game
	const startGameBtn = document.querySelector('.start-game-btn');
	startGameBtn.addEventListener('click', requestStartGameEL);

	// start game countdown to show user how much time is left for game to start
	const startGameTimerId = setInterval(
		() =>
			(document.querySelector('.loading').textContent =
				getSecondsLeftFrom(gameStartTime)),
		1000
	);

	// start game after this timeout
	const startGameAfterTimeoutId = setTimeout(
		requestStartGameEL,
		getSecondsLeftFrom(gameStartTime) * 1000
	);

	function requestStartGameEL() {
		// runs when the game starts, makes socket conn call to server to start the game
		// clear the countdown timers
		clearAllIntervals(startGameTimerId, startGameAfterTimeoutId);

		// generate response and send
		const responseMsg = {
			type: 7,
			content: 'start the game bro!',
			poolId,
		};

		sendViaSocket(responseMsg);
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

function displayScores(socketMessage) {
	console.table(socketMessage);

	const dataArr = JSON.parse(socketMessage.content);

	let html = `<table>
	<tr>
		<th>Name</th>
		<th>Score</th>
	</tr>`;
	dataArr.forEach(
		item => (html += `<tr><td>${item.name}</td><td>${item.score}</td></tr>`)
	);
	html += `</table>`;

	document.querySelector('.score-board').innerHTML = html;

	clearAllIntervals(renderClientsTimerId, wordExpiryTimerId);
}

// ------------------------ start game ------------------------

function beginClientSketchingFlow(socketMessage) {
	console.table(socketMessage);

	// initialise the time at which this word expires
	currentWordExpiresAt = new Date(socketMessage.currWordExpiresAt).getTime();

	// start timer for the word expiry
	wordExpiryTimerId = setInterval(async () => {
		const timeLeftDiv = document.querySelector('.time-left-for-word');
		timeLeftDiv.classList.remove('hidden');

		const secondsLeft = getSecondsLeftFrom(currentWordExpiresAt);
		timeLeftDiv.querySelector('span').textContent = secondsLeft;

		if (secondsLeft <= 0) {
			clearInterval(wordExpiryTimerId);
			console.log('timer for word cleared');

			// trigger next word for next player: TODO
			// requestCanvasClear();

			// const responseMsg = {
			// 	type: 8,
			// 	content: 'next word',
			// };

			// await wait(5 * 1000);
			// sendViaSocket(responseMsg);
		}
	}, 1000);

	// for enabling drawing access if clientId matches
	if (clientId === socketMessage.currSketcherId) {
		paintUtils.isAllowedToPaint = true;

		// display the word by unhiding the painter-utils div
		document.querySelector('.painter-utils').classList.remove('hidden');
		document.querySelector('.your-word').textContent = socketMessage.currWord;

		// add EL for clearing the canvas
		document
			.querySelector('.clear-canvas')
			.addEventListener('click', requestCanvasClear);
	} else {
		paintUtils.isAllowedToPaint = false;
		document.querySelector('.painter-utils').classList.add('hidden');
		document.querySelector('.your-word').textContent = '';
	}
}

function startGame(socketMessage) {
	// called when socket receives message from server with type as 6
	if (socketMessage.content !== 'true') return;

	console.log('game started by server');
	paintUtils.hasGameStarted = true;

	// hide the div and toggle paintUtils.has Game Started
	const startGameDiv = document.querySelector('.start-game');
	startGameDiv && startGameDiv.classList.add('hidden');

	beginClientSketchingFlow(socketMessage);
}

// ------------------------ get all clients and render ------------------------

function getAllClientsEL() {
	// makes a socket connection call to request client info list
	const responseMsg = {
		type: 6,
		content: '',
		poolId,
	};

	sendViaSocket(responseMsg);
}

function renderClients(allClients) {
	// called when the socket conn receives a message from server as type 5

	const membersDiv = document.querySelector('.members');
	membersDiv.innerHTML = '';

	// parse array of objects into json
	allClients = JSON.parse(allClients);

	// render
	allClients.forEach((n, i) => {
		const clientNameHolder = document.createElement('div');
		const clientName = document.createElement('p');

		clientName.innerHTML = `#${i + 1} ${n.name}: ${n.score} points`;
		clientName.style.color = `#${n.color}`;
		clientNameHolder.appendChild(clientName);

		membersDiv.appendChild(clientNameHolder);
	});
}

// ------------------------------------- chat -------------------------------------

// event listeners to send chat messages

function appendChatMsgToDOM(msg) {
	// adds the msg into the DOM

	if (msg.length === 0 || msg === '') return;

	const messagesDiv = document.querySelector('.messages');
	const msgDiv = document.createElement('div');
	const text = document.createTextNode(msg);
	msgDiv.appendChild(text);
	messagesDiv.appendChild(msgDiv);

	document.querySelector('.msg').value = '';
}

function sendChatMsgBtnEL(e) {
	// event listener to send chat message

	e.preventDefault();
	const msg = document.querySelector('.msg').value;

	if (msg.length === 0 || msg === '') return;

	// create string response object
	const responseMsg = {
		type: 3,
		content: msg,
		clientName,
		clientId,
	};

	// convert object to string to transmit
	sendViaSocket(responseMsg);
}

// ------------------------------------- canvas -------------------------------------

// event listeners for canvas painting

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
	};

	// sending canvas data
	sendViaSocket(responseMsg);
}

function addCanvasEventListeners() {
	document.addEventListener('mousedown', startPainting);
	document.addEventListener('mouseup', stopPainting);
	document.addEventListener('mousemove', paint);
}

function clearAllIntervals(...ids) {
	ids.forEach(i => clearInterval(i));
}
