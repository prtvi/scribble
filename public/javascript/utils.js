'use strict';

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

function gameStartTimer() {
	// start game countdown to show user how much time is left for game to start
	return runTimer(document.querySelector('.loading'), gameStartTime);
}

function getClientNameDiv(clientInfo, iteration) {
	// client name div
	const clientNameDiv = document.createElement('div');
	clientNameDiv.classList.add('member');

	// client num span
	const clientNumSpan = document.createElement('span');
	clientNumSpan.classList.add('member-num');
	clientNumSpan.textContent = `#${iteration + 1}`;

	// client name span
	const clientNameSpan = document.createElement('span');
	clientNameSpan.classList.add('member-name');
	clientNameSpan.style.color = `#${clientInfo.color}`;

	if (clientName === clientInfo.name)
		clientNameSpan.textContent = `${clientInfo.name} (you)`;
	else clientNameSpan.textContent = clientInfo.name;

	// client score span
	const clientScoreSpan = document.createElement('span');
	clientScoreSpan.classList.add('member-score');
	clientScoreSpan.textContent = `${clientInfo.score} points`;

	// append everything to client name div
	clientNameDiv.appendChild(clientNumSpan);
	clientNameDiv.appendChild(clientNameSpan);
	clientNameDiv.appendChild(clientScoreSpan);

	return clientNameDiv;
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

// render animation/transition for changing innerHTML - https://stackoverflow.com/questions/29640486

const overlayFadeInAnimationDuration = 300; // to be configured in css file too, #overlay{}

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
}

function disableSketching() {
	const painterUtilsDiv = document.querySelector('.painter-utils');
	const clearCanvasBtn = document.querySelector('.clear-canvas');

	paintUtils.isAllowedToPaint = false;

	// display painter utils div and remove EL
	painterUtilsDiv.classList.add('hidden');
	clearCanvasBtn.removeEventListener('click', requestCanvasClear);
}

function beginClientSketchingFlowInit(socketMessage) {
	// initialise the time at which this word expires
	const currentWordExpiresAt = new Date(
		socketMessage.currWordExpiresAt
	).getTime();

	const timeLeftSpan = document.querySelector('.time-left span');
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
	document.querySelector('.time-left span').textContent = '0s';
}

function removeEventListenersOnGameStart() {
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
		document.addEventListener('mousedown', startPainting);
		document.addEventListener('mouseup', stopPainting);
		document.addEventListener('mousemove', paint);
	});

	// copy joining link
	document
		.querySelector('.joining-link-btn')
		.addEventListener('click', copyJoiningLinkEL);

	// add event listener to start game button to start game
	document
		.querySelector('.start-game-btn')
		.addEventListener('click', startGameEl);

	// adjust overlay position on scroll
	window.addEventListener('scroll', adjustOverlay);
}

function copyJoiningLinkEL() {
	navigator.clipboard.writeText(joiningLink);
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
