'use strict';

// -------------------------------- UTILS --------------------------------

function log(...args) {
	if (allowLogs) console.log(...args);
}

function wait(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

function getDomain() {
	// extract domain from url
	const url = window.location.href;
	const fi = url.indexOf('/');
	const li = url.lastIndexOf('/');
	const domain = url.slice(fi + 2, li);

	return domain;
}

function getSecondsLeftFrom(futureTime) {
	const now = new Date().getTime();
	const diff = futureTime - now;
	return Math.round(diff / 1000);
}

function clearAllIntervals(...ids) {
	ids.forEach(i => clearInterval(i));
}

function getAvatarDom(avatarConfig) {
	function setBgPosition(element, x, y) {
		const offset = 48 * scaleAvatarBy;

		element.style.backgroundPositionX = `-${x * offset}px`;
		element.style.backgroundPositionY = `-${y * offset}px`;
	}

	const playerAvatar = document.createElement('div');
	playerAvatar.classList.add('avatar', 'player-avatar');

	const pColor = document.createElement('div');
	pColor.classList.add('color');
	setBgPosition(pColor, avatarConfig.color.x, avatarConfig.color.y);

	const pEyes = document.createElement('div');
	pEyes.classList.add('eyes');
	setBgPosition(pEyes, avatarConfig.eyes.x, avatarConfig.eyes.y);

	const pMouth = document.createElement('div');
	pMouth.classList.add('mouth');
	setBgPosition(pMouth, avatarConfig.mouth.x, avatarConfig.mouth.y);

	const pOwner = document.createElement('div');
	pOwner.classList.add('owner');
	if (avatarConfig.isOwner) pOwner.classList.add('active');

	const pCrowned = document.createElement('div');
	pCrowned.classList.add('crowned');
	if (avatarConfig.isCrowned) pCrowned.classList.add('active');

	playerAvatar.appendChild(pColor);
	playerAvatar.appendChild(pEyes);
	playerAvatar.appendChild(pMouth);
	playerAvatar.appendChild(pOwner);
	playerAvatar.appendChild(pCrowned);

	return playerAvatar;
}

function getPlayerDom(playerInfo, iteration) {
	// player name div
	const player = document.createElement('div');
	player.classList.add('player');

	// player num span
	const playerNum = document.createElement('span');
	playerNum.classList.add('num');
	playerNum.textContent = `#${iteration + 1}`;

	// player name span
	const playerName = document.createElement('span');
	playerName.classList.add('name');

	playerName.textContent = playerInfo.name;
	if (clientId === playerInfo.id) playerName.textContent += ' (you)';

	// player score span
	const playerScore = document.createElement('span');
	playerScore.classList.add('score');
	playerScore.textContent = `${playerInfo.score} points`;

	const playerNameAndScore = document.createElement('div');
	playerNameAndScore.appendChild(playerName);
	playerNameAndScore.appendChild(playerScore);

	// player avatar
	const playerAvatar = getAvatarDom(playerInfo.avatarConfig);

	// append everything to player div
	player.appendChild(playerNum);
	player.appendChild(playerNameAndScore);
	player.appendChild(playerAvatar);

	return player;
}

function sendChatMsgBtnEL(e) {
	// event listener to send chat message

	e.preventDefault();
	const msg = document.querySelector('.msg').value.trim();

	if (msg.length === 0 || msg === '') return;

	// create string response object
	const socketMsg = {
		type: 3,
		typeStr: messageTypeMap.get(3),
		content: msg,
		clientName,
		clientId,
		poolId,
	};

	// convert object to string to transmit
	sendViaSocket(socketMsg);
}

function getOverlayHtmlForTextOnly(overlayText) {
	return `
	<div class="overlay-content">
		<div>
			<p class="overlay-text">
				${overlayText}
			</p>
		</div>
	</div>`;
}

function displayOverlay(html) {
	// display overlay after some delay to render fade in animation
	// if event listeners are to be added to the given html, then use the same timeout to attach the event listeners

	overlay.style.opacity = 0;
	setTimeout(function () {
		overlay.innerHTML = html;
		overlay.style.display = 'flex';
		overlay.style.opacity = 1;
		adjustOverlay();
	}, overlayFadeInAnimationDuration);
}

function hideOverlay() {
	// render hiding animation using timeout
	overlay.style.opacity = 1;
	setTimeout(() => {
		overlay.style.opacity = 0;
		overlay.innerHTML = '';
	}, overlayFadeInAnimationDuration);

	// change overlay display property to none after the animation
	setTimeout(() => (overlay.style.display = 'none'), 1000);
}

function adjustOverlay() {
	// adjust overlay position on scroll
	const cc = document.querySelector('.canvas-container');
	overlay.style.top = `${cc.offsetTop}px`;
	overlay.style.height = `${cc.offsetHeight}px`;
	overlay.style.width = `${cc.offsetWidth}px`;
}

function saveToLocalStorage(key, value) {
	window.localStorage.setItem(key, JSON.stringify(value));
}

function getFromLocalStorage(key) {
	return window.localStorage.getItem(key);
}

function disableSketching() {
	const painterUtilsDiv = document.querySelector('.painter-utils');
	const clearCanvasBtn = document.querySelector('.pu.clear');
	const undoBtn = document.querySelector('.pu.undo');
	const strokes = document.querySelector('.strokes');
	const brushStrokeSelect = document.querySelector('.pu.brush-stroke-select');

	paintUtils.isAllowedToPaint = false;

	// display painter utils div and remove EL
	painterUtilsDiv.classList.add('hidden');

	clearCanvasBtn.removeEventListener('click', requestCanvasClear);
	undoBtn.removeEventListener('click', undo);
	strokes.removeEventListener('click', selectStrokeEL);
	brushStrokeSelect.removeEventListener('click', openBrushStrokeSelect);
}

function beginClientSketchingFlowInit(socketMessage) {
	// initialise the time at which this word expires
	const currentWordExpiresAt = new Date(
		socketMessage.currWordExpiresAt
	).getTime();

	const timeLeftSpan = document.querySelector('.timer span');
	timeLeftSpan.textContent = `${timeForEachWordInSeconds}s`;
	return runTimer(timeLeftSpan, currentWordExpiresAt);
}

function runTimer(timerElement, timeoutAt) {
	const countdownTimer = setInterval(function () {
		const secondsLeft = getSecondsLeftFrom(timeoutAt);
		if (secondsLeft > -1) timerElement.textContent = `${secondsLeft}s`;
		else clearInterval(countdownTimer);
	}, 1000);
	return countdownTimer;
}

function showZeroOnTimeLeftSpan() {
	// to display 0s left, in the event that everyone guesses the word before the timeout
	document.querySelector('.timer span').textContent = '0s';
}

function removeEventListenersOnGameStart() {
	const isOwner = JSON.parse(getFromLocalStorage('avatarConfig')).isOwner;
	if (isOwner)
		document
			.querySelector('.start-game-btn')
			.removeEventListener('click', startGameEl);

	document
		.querySelector('.joining-link-btn')
		.removeEventListener('click', copyJoiningLinkEL);

	document.querySelector('.joining-link-div').classList.add('hidden');
}

// event listeners

function initGlobalEventListeners() {
	// chat
	document
		.querySelector('.send-msg')
		.addEventListener('click', sendChatMsgBtnEL);

	// show number of characters typed in chat box
	const lenIndicator = document.querySelector('.input-wrapper span');
	document
		.querySelector('.msg')
		.addEventListener(
			'input',
			e => (lenIndicator.textContent = e.target.value.length)
		);

	// event listeners for drawing
	window.addEventListener('load', () => {
		canvas.addEventListener('mousedown', startPainting);
		canvas.addEventListener('mouseup', stopPainting);
		canvas.addEventListener('mousemove', paint);
	});

	// resize canvas on window resize
	window.addEventListener('resize', function () {
		const { w, h } = getCanvasSize();
		canvas.width = w;
		canvas.height = h;

		const cc = document.querySelector('.canvas-container');
		cc.style.width = `${w}px`;
		cc.style.height = `${h}px`;

		adjustOverlay();
	});

	// copy joining link
	document
		.querySelector('.joining-link-btn')
		.addEventListener('click', copyJoiningLinkEL);

	// add event listener to start game button to start game
	const isOwner = JSON.parse(getFromLocalStorage('avatarConfig')).isOwner;
	const startGameBtn = document.querySelector('.start-game-btn');

	if (isOwner) startGameBtn.addEventListener('click', startGameEl);
	else startGameBtn.classList.add('hidden');

	// adjust overlay position on scroll
	window.addEventListener('scroll', adjustOverlay);

	// modal
	const modal = document.getElementById('modal');
	document
		.querySelector('.close-modal')
		.addEventListener('click', () => (modal.style.display = 'none'));

	window.addEventListener('click', e => {
		if (e.target === modal && modal.style.display != 'none')
			modal.style.display = 'none';
	});
}

function copyJoiningLinkEL() {
	navigator.clipboard.writeText(joiningLink);
	appendChatMsgToDOM('Copied to clipboard!', '#0043ff');
}

function startGameEl() {
	const socketMsg = {
		type: 7,
		typeStr: messageTypeMap.get(7),
		clientId,
		clientName,
		poolId,
	};

	sendViaSocket(socketMsg);
}

function wordChooseEL(e) {
	// passing parameters to addEvenListener function by assigning the "this" element a new parameter, say words, then using it in the event handler like e.currentTarget.words

	const chosenWord = e.target.textContent.trim();
	if (!e.currentTarget.words.includes(chosenWord)) return;

	const socketMsg = {
		type: 34,
		typeStr: messageTypeMap.get(34),
		content: chosenWord,
		clientName,
		clientId,
		poolId,
	};

	sendViaSocket(socketMsg);
}

function selectStrokeEL(e) {
	if (!e.target.classList.contains('stroke')) return;

	const strokeElems = document.getElementsByClassName('stroke');
	Array.from(strokeElems).forEach(ele => ele.classList.remove('active'));
	e.target.classList.add('active');

	const size = +e.target.id.slice(1);
	paintUtils.lineWidth = size * 2;

	const viewImg = document.querySelector('.brush-stroke-select img');
	viewImg.style.width = `${size * minBrushStrokeSizeForImg}px`;
	viewImg.style.height = `${size * minBrushStrokeSizeForImg}px`;

	const strokeSelectMenu = document.querySelector('.stroke-select');
	setTimeout(() => strokeSelectMenu.classList.add('hidden'), 100);
}

function openBrushStrokeSelect() {
	const strokeSelectMenu = document.querySelector('.stroke-select');
	strokeSelectMenu.classList.remove('hidden');

	const strokes = document.querySelector('.strokes');
	strokes.addEventListener('click', selectStrokeEL);
}

// -------------------------------- CANVAS --------------------------------

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

function getMousePos(event) {
	const clientRect = canvas.getBoundingClientRect();
	return {
		x: Math.round(event.clientX - clientRect.left),
		y: Math.round(event.clientY - clientRect.top),
	};
}

function startPainting(event) {
	paintUtils.isPainting = true;

	paintUtils.prevMouse = {
		x: paintUtils.mouse.x,
		y: paintUtils.mouse.y,
	};
	paintUtils.mouse = getMousePos(event);

	paintUtils.points = [];
	paintUtils.points.push(paintUtils.mouse);
}

function stopPainting() {
	paintUtils.isPainting = false;
	paintUtils.pointsHistory.push(paintUtils.points);
}

async function paint(event) {
	if (!paintUtils.isPainting) return;
	if (!paintUtils.hasGameStarted) return;
	if (!paintUtils.isAllowedToPaint) return;

	ctx.lineWidth = paintUtils.lineWidth;
	ctx.lineCap = paintUtils.lineCap;
	ctx.lineJoin = paintUtils.lineJoin;
	ctx.strokeStyle = paintUtils.strokeStyle;

	paintUtils.prevMouse = {
		x: paintUtils.mouse.x,
		y: paintUtils.mouse.y,
	};
	paintUtils.mouse = getMousePos(event);
	paintUtils.points.push(paintUtils.mouse);

	ctx.beginPath();
	ctx.moveTo(paintUtils.prevMouse.x, paintUtils.prevMouse.y);
	ctx.lineTo(paintUtils.mouse.x, paintUtils.mouse.y);
	ctx.stroke();

	await wait(500);
	sendImgData();
}

function drawPaths() {
	clearCanvas();

	paintUtils.pointsHistory.forEach(path => {
		if (path.length === 0) return;

		ctx.beginPath();
		ctx.moveTo(path[0].x, path[0].y);

		for (let i = 1; i < path.length; i++) ctx.lineTo(path[i].x, path[i].y);

		ctx.stroke();
	});
}

function undo() {
	paintUtils.pointsHistory.splice(-1, 1);
	drawPaths();
	sendImgDataForUndoAction();
}

function requestCanvasClear() {
	// clear canvas and request clear on rest of the clients
	clearCanvas();

	paintUtils.points = [];
	paintUtils.pointsHistory = [];

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

function sendImgDataForUndoAction() {
	// called by undo function
	const socketMsg = {
		type: 41,
		typeStr: messageTypeMap.get(41),
		content: String(canvas.toDataURL('img/png')),
		clientName,
		clientId,
		poolId,
	};

	// sending canvas data
	sendViaSocket(socketMsg);
}

// -------------------------------- ON MESSAGE HANDLERS --------------------------------

// 1, 2, 3, 31, 312, 313
function appendChatMsgToDOM(msg, formatColor) {
	// adds the msg into the DOM

	if (msg.length === 0 || msg === '') return;

	const messagesDiv = document.querySelector('.messages');

	const newMsgDiv = document.createElement('div');
	newMsgDiv.classList.add('message');

	const text = document.createElement('span');
	text.style.color = formatColor || '#1d1d1f'; // f5f5f7

	newMsgDiv.style.backgroundColor = `${formatColor}20`;
	text.innerHTML = msg; // TODO: handle better

	newMsgDiv.appendChild(text);
	messagesDiv.appendChild(newMsgDiv);

	newMsgDiv.scrollIntoView();

	document.querySelector('.msg').value = '';
	document.querySelector('.input-wrapper span').textContent = 0;
}

// 32
function revealWordOnOverlayAndChat(socketMessage) {
	const message = `The word was '${socketMessage.content}'`;
	displayOverlay(getOverlayHtmlForTextOnly(message));
	appendChatMsgToDOM(message, '#ffa500');
	document.querySelector('.word span.content').textContent =
		socketMessage.content;
}

// 33
function showWordToChoose(socketMessage) {
	const words = JSON.parse(socketMessage.content);

	let html = `<div class="overlay-content">
	<div><p class="overlay-text">Your turn, choose a word to draw!</p></div>
	<div class="word-options">`;
	words.forEach(
		w => (html += `<div class="word-option"><span>${w}</span></div>`)
	);
	html += `</div> <div><div class="word-choose-timer"><span>${timeForChoosingWordInSeconds}s</span></div> </div>`;

	displayOverlay(html);

	setTimeout(() => {
		const optionsEle = overlay.querySelector('.word-options');
		optionsEle.words = words;
		optionsEle.addEventListener('click', wordChooseEL);

		const timeoutAt = new Date(socketMessage.timeoutAfter).getTime();
		const timerEle = overlay.querySelector('div.word-choose-timer span');
		timerEle.textContent = `${timeForChoosingWordInSeconds}s`;
		runTimer(timerEle, timeoutAt);
	}, overlayFadeInAnimationDuration + 50);
}

// 35
function showChoosingWordOnOverlay(socketMessage) {
	displayOverlay(
		getOverlayHtmlForTextOnly(
			`${socketMessage.currSketcherName} is choosing a word!`
		)
	);
}

// 4
function displayImgOnCanvas(socketMessage) {
	// display image data on canvas
	const img = new Image();
	// scale up/down canvas data based on current canvas size using outer bounds
	img.onload = () => ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
	img.setAttribute('src', socketMessage.content);
}

// 41
function displayUndoCanvas(socketMessage) {
	clearCanvas();
	displayImgOnCanvas(socketMessage);
}

// 5
function clearCanvas() {
	ctx.clearRect(0, 0, canvas.width, canvas.height);
}

// 6
function renderClients(allClients) {
	// called when the socket conn receives a message from server as type 6

	if (allClients.length === 0) return;

	const membersDiv = document.querySelector('.players');
	membersDiv.innerHTML = '';

	// parse array of objects into json
	allClients = JSON.parse(allClients);

	// render
	allClients.forEach((n, i) => membersDiv.appendChild(getPlayerDom(n, i)));
}

// 70
function startGame(socketMessage) {
	// called when socket receives message from server with type as 6
	if (!socketMessage.success) return;

	log('game started');

	// flag game started
	paintUtils.hasGameStarted = true;

	// remove event listeners
	removeEventListenersOnGameStart();

	// display game started overlay
	document.querySelector('.word span.content').textContent = 'Game started!';

	displayOverlay(getOverlayHtmlForTextOnly('Game started!'));
	document.querySelector('.word span.status').textContent = 'Game started';
}

// 71
function renderRoundDetails(socketMessage) {
	document.querySelector('.round span.curr-round').textContent =
		socketMessage.currRound;
	displayOverlay(getOverlayHtmlForTextOnly(`Round ${socketMessage.currRound}`));
}

// 8
function beginClientSketchingFlow(socketMessage) {
	hideOverlay();
	overlay
		.querySelector('.word-options')
		.removeEventListener('click', wordChooseEL);

	const wordExpiryCountdown = beginClientSketchingFlowInit(socketMessage);

	// for enabling drawing access if clientId matches
	const painterUtilsDiv = document.querySelector('.painter-utils');
	const clearCanvasBtn = document.querySelector('.pu.clear');
	const undoBtn = document.querySelector('.pu.undo');
	const brushStrokeSelect = document.querySelector('.pu.brush-stroke-select');

	paintUtils.isAllowedToPaint = true;

	// display the word
	document.querySelector('.word span.status').textContent = 'Draw this!';
	document.querySelector('.word span.content').textContent =
		socketMessage.currWord;

	// display painter utils div and add EL for clearing the canvas
	painterUtilsDiv.classList.remove('hidden');
	clearCanvasBtn.addEventListener('click', requestCanvasClear);
	undoBtn.addEventListener('click', undo);
	brushStrokeSelect.addEventListener('click', openBrushStrokeSelect);

	return wordExpiryCountdown;
}

// 87
function showSketcherBeginDrawing(socketMessage) {
	displayOverlay(
		getOverlayHtmlForTextOnly(
			`${socketMessage.currSketcherName} is now drawing!`
		)
	);

	setTimeout(hideOverlay, 2000);
}

// 88
function showSketcherIsDrawing(socketMessage) {
	const wordExpiryCountdown = beginClientSketchingFlowInit(socketMessage);

	let text = '';
	for (let i = 0; i < socketMessage.currWordLen; i++) text += '_ ';
	text = text.trim();

	document.querySelector('.word span.status').textContent = 'Guess this!';
	document.querySelector('.word span.content').textContent = text;

	return wordExpiryCountdown;
}

// 81
function disableSketchingTurnOver() {
	clearAllIntervals(wordExpiryTimer);
	disableSketching();
	showTimeUp();
}

// 82
function showTimeUp() {
	showZeroOnTimeLeftSpan();
	clearAllIntervals(wordExpiryTimer);
	displayOverlay(getOverlayHtmlForTextOnly('Time up!'));
}

// 83
function disableSketchingAllGuessed() {
	clearAllIntervals(wordExpiryTimer);
	disableSketching();
	showAllHaveGuessed();
}

// 84
function showAllHaveGuessed() {
	showZeroOnTimeLeftSpan();
	clearAllIntervals(wordExpiryTimer);
	displayOverlay(getOverlayHtmlForTextOnly('Everyone guessed the word!'));
}

// 9
function displayScores(socketMessage) {
	const dataArr = JSON.parse(socketMessage.content);

	let html = `<div class="overlay-content">
	<div><p class="overlay-text">Game over!</p></div>`;
	html += `<div> <table>
	<tr>
		<th>Name</th>
		<th>Score</th>
	</tr>`;
	dataArr.forEach(
		item => (html += `<tr><td>${item.name}</td><td>${item.score}</td></tr>`)
	);
	html += `</table> </div> </div>`;

	displayOverlay(html);
	appendChatMsgToDOM('Game over!', '#ff0000');
}

// 10
function makeMessageTypeMapGlobal(socketMessage) {
	const content = JSON.parse(socketMessage.content);

	timeForEachWordInSeconds = content.timeForEachWordInSeconds;
	timeForChoosingWordInSeconds = content.timeForChoosingWordInSeconds;
	allowLogs = content.printLogs;

	const m = content.messageTypeMap;
	const keys = Object.keys(m);
	messageTypeMap = new Map();

	keys.forEach(k => messageTypeMap.set(Number(k), m[k]));
}

// -------------------------------- SOCKET --------------------------------

function initSocket() {
	// initialises socket connection and adds corresponding function handlers to the socket

	const avatarConfig = getFromLocalStorage('avatarConfig');

	const wsUrl = `ws://${getDomain()}/ws?poolId=${poolId}&clientId=${clientId}&clientName=${clientName}&avatarConfig=${avatarConfig}`;

	const socket = new WebSocket(wsUrl);

	socket.onopen = () => log('Socket successfully connected!');
	socket.onerror = error => log('Socket error', error);
	socket.onmessage = socketOnMessage;
	socket.onclose = socketOnClose;

	return socket;
}

function socketOnMessage(message) {
	// runs when a message is received on the socket conn, runs the corresponding functions depending on message type

	// parse json string into json object
	const socketMessage = JSON.parse(message.data);

	if (socketMessage.type !== 4) log(socketMessage.type, socketMessage.typeStr);

	switch (socketMessage.type) {
		case 1:
			if (socketMessage.clientId === clientId)
				// if the current clientId and the clientId from response match then
				appendChatMsgToDOM(
					`You joined the room as <strong>${socketMessage.clientName}</strong>!`,
					''
				);
			else
				appendChatMsgToDOM(
					`<strong>${socketMessage.clientName}</strong> has joined the room!`,
					''
				);
			break;

		case 2:
			appendChatMsgToDOM(
				`<strong>${socketMessage.clientName}</strong> has left the room!`,
				''
			);
			break;

		case 3:
			appendChatMsgToDOM(
				`<strong>${socketMessage.clientName}</strong>: ${socketMessage.content}`,
				''
			);
			break;

		case 31:
			appendChatMsgToDOM(
				`<strong>${socketMessage.clientName}</strong> guessed the word!`,
				'#00ff00'
			);
			break;

		case 312:
			appendChatMsgToDOM(
				`Naughty <strong>@${socketMessage.clientName}</strong>`,
				'#ff0000'
			);
			break;

		case 313:
			appendChatMsgToDOM(
				`You can't reveal the word <strong>@${socketMessage.clientName}</strong>`,
				'#ff7f00'
			);
			break;

		case 32:
			revealWordOnOverlayAndChat(socketMessage);
			break;

		case 33:
			showWordToChoose(socketMessage);
			break;

		case 35:
			showChoosingWordOnOverlay(socketMessage);
			break;

		case 4:
			displayImgOnCanvas(socketMessage);
			break;

		case 41:
			displayUndoCanvas(socketMessage);
			break;

		case 5:
		case 51:
			clearCanvas();
			break;

		case 6:
			renderClients(socketMessage.content);
			break;

		case 69:
			appendChatMsgToDOM(
				'You need at least two players to start the game',
				'#457ef4'
			);
			break;

		case 70:
			startGame(socketMessage);
			break;

		case 71:
			renderRoundDetails(socketMessage);
			break;

		case 8:
			wordExpiryTimer = beginClientSketchingFlow(socketMessage);
			break;

		case 87:
			showSketcherBeginDrawing(socketMessage);
			break;

		case 88:
			wordExpiryTimer = showSketcherIsDrawing(socketMessage);
			break;

		case 81:
			disableSketchingTurnOver();
			break;

		case 82:
			showTimeUp();
			break;

		case 83:
			disableSketchingAllGuessed();
			break;

		case 84:
			showAllHaveGuessed();
			break;

		case 9:
			displayScores(socketMessage);
			break;

		case 10:
			makeMessageTypeMapGlobal(socketMessage);
			break;

		default:
			break;
	}
}

function socketOnClose() {
	// TODO: show to user when disconnected from server
	// on socket conn close, stop all timer or intervals
	log('Socket connection closed, stopping timers and timeouts!');
	clearAllIntervals(wordExpiryTimer);

	document.getElementById('modal').style.display = 'flex';
}

function sendViaSocket(socketMsg) {
	/*  socket.readyState: int
			0 - connecting
			1 - open
			2 - closing
			3 - closed
	*/

	if (socket.readyState === socket.OPEN) socket.send(JSON.stringify(socketMsg));
	else {
		log(
			'0: connecting | 1: open | 2: closing | 3: closed, current state:',
			socket.readyState
		);

		clearAllIntervals(wordExpiryTimer);
		document.getElementById('modal').style.display = 'flex';
	}
}

// -------------------------------- MAIN --------------------------------

// to be configured in css file too, #overlay{}, render animation/transition for changing innerHTML - https://stackoverflow.com/questions/29640486
const overlayFadeInAnimationDuration = 300;
const scaleAvatarBy = 0.5;
const minBrushStrokeSizeForImg = 6;

// canvas, canvas ctx and overlay init
const { canvas, ctx, overlay } = initCanvasAndOverlay();

// utils for painting on canvas
const paintUtils = {
	isPainting: false,
	hasGameStarted: false,
	isAllowedToPaint: false,
	points: [],
	pointsHistory: [],
	mouse: { x: 0, y: 0 },
	prevMouse: { x: 0, y: 0 },
	lineWidth: 2,
	lineCap: 'round',
	lineJoin: 'round',
	strokeStyle: '#000',
};

let messageTypeMap,
	timeForEachWordInSeconds,
	timeForChoosingWordInSeconds,
	wordExpiryTimer,
	allowLogs;

// init socket connection and check game begin status
const socket = initSocket();
initGlobalEventListeners();
