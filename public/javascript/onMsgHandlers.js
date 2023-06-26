'use strict';

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
	displayOverlay(
		`<div class="overlay-div"><p class="overlay-p">The word was '${socketMessage.content}'</p></div>`
	);

	appendChatMsgToDOM(`The word was '${socketMessage.content}'`, '#ffa500');
}

// 33
function showWordToChoose(socketMessage) {
	const words = JSON.parse(socketMessage.content);

	let html = `<div class="overlay-div">
	<p class="overlay-p">Your turn! Choose a word to draw</p>`;
	words.forEach(w => (html += `<span class="word-option">${w}</span>`));
	html += `<span class="word-choose-timer"></span> </div>`;

	displayOverlay(html);

	overlay.querySelector('div').addEventListener('click', function (e) {
		const chosenWord = e.target.textContent.trim();
		if (!words.includes(chosenWord)) return;

		const socketMsg = {
			type: 34,
			typeStr: messageTypeMap.get(34),
			content: chosenWord,
			clientName,
			clientId,
			poolId,
		};

		sendViaSocket(socketMsg);
	});

	const timeoutAt = new Date(socketMessage.timeoutAfter).getTime();
	const timerEle = overlay.querySelector('span.word-choose-timer');
	timerEle.textContent = `${timeForChoosingWordInSeconds}s`;
	runTimer(timerEle, timeoutAt);
}

// 35
function showChoosingWordOnOverlay(socketMessage) {
	displayOverlay(
		`<div class="overlay-div"><p class="overlay-p">${socketMessage.currSketcherName} is choosing a word!</p></div>`
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

// 5
function clearCanvas() {
	hideOverlay();
	ctx.clearRect(0, 0, canvas.width, canvas.height);
}

// 6
function renderClients(allClients) {
	// called when the socket conn receives a message from server as type 6

	if (allClients.length === 0) return;

	const membersDiv = document.querySelector('.members');
	membersDiv.innerHTML = '';

	// parse array of objects into json
	allClients = JSON.parse(allClients);

	// render
	allClients.forEach((n, i) => membersDiv.appendChild(getClientNameDiv(n, i)));
}

// 70
function startGame(socketMessage) {
	// called when socket receives message from server with type as 6
	if (!socketMessage.success) return;

	console.log('game started by server');

	// flag game started
	paintUtils.hasGameStarted = true;
	clearAllIntervals(startGameTimerId);

	// remove event listeners
	document
		.querySelector('.start-game-btn')
		.removeEventListener('click', startGameEl);

	hideAndRemoveElForJoiningLink();

	// display game started overlay
	displayOverlay(
		`<div class="overlay-div"><p class="overlay-p">Game started</p></div>`
	);

	document.querySelector('.time-left span').textContent = 'Game started';
}

// 71
function renderRoundDetails(socketMessage) {
	const roundDiv = document.querySelector('.round');
	roundDiv.classList.remove('hidden');

	roundDiv.querySelector(
		'span'
	).textContent = `Round: ${socketMessage.currRound}`;

	displayOverlay(
		`<div class="overlay-div"><p class="overlay-p">Round: ${socketMessage.currRound}</p></div>`
	);
}

// 8
function beginClientSketchingFlow(socketMessage) {
	const wordExpiryCountdown = beginClientSketchingFlowInit(socketMessage);

	// for enabling drawing access if clientId matches
	const wordDiv = document.querySelector('.word');
	wordDiv.classList.remove('hidden');
	const wordSpan = wordDiv.querySelector('span');

	const painterUtilsDiv = document.querySelector('.painter-utils');
	const clearCanvasBtn = document.querySelector('.clear-canvas');
	paintUtils.isAllowedToPaint = true;

	// display the word
	wordSpan.textContent = socketMessage.currWord;

	// display painter utils div and add EL for clearing the canvas
	painterUtilsDiv.classList.remove('hidden');
	clearCanvasBtn.addEventListener('click', requestCanvasClear);

	return wordExpiryCountdown;
}

// 88
function showClientDrawing(socketMessage) {
	const wordExpiryCountdown = beginClientSketchingFlowInit(socketMessage);

	const wordDiv = document.querySelector('.word');
	wordDiv.classList.remove('hidden');
	const wordSpan = wordDiv.querySelector('span');
	wordSpan.textContent = `${socketMessage.currWordLen} characters`;

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
	clearAllIntervals(wordExpiryTimer);
	displayOverlay(
		`<div class="overlay-div"><p class="overlay-p">Time up!</p></div>`
	);
}

// 83
function disableSketchingAllGuessed() {
	clearAllIntervals(wordExpiryTimer);
	disableSketching();
	showAllHaveGuessed();
}

// 84
function showAllHaveGuessed() {
	clearAllIntervals(wordExpiryTimer);
	displayOverlay(
		`<div class="overlay-div"><p class="overlay-p">Everyone guessed the word!</p></div>`
	);
}

// 9
function displayScores(socketMessage) {
	const dataArr = JSON.parse(socketMessage.content);

	let html = `<div class="overlay-div">
	<p class="overlay-p">Game over!</p>
	<table>
	<tr>
		<th>Name</th>
		<th>Score</th>
	</tr>`;
	dataArr.forEach(
		item => (html += `<tr><td>${item.name}</td><td>${item.score}</td></tr>`)
	);
	html += `</table> </div>`;

	displayOverlay(html);
}

// 10
function makeMessageTypeMapGlobal(socketMessage) {
	const content = JSON.parse(socketMessage.content);

	timeForEachWordInSeconds = content.timeForEachWordInSeconds;
	timeForChoosingWordInSeconds = content.timeForChoosingWordInSeconds;

	const m = content.messageTypeMap;
	const keys = Object.keys(m);
	messageTypeMap = new Map();

	keys.forEach(k => messageTypeMap.set(Number(k), m[k]));
}
