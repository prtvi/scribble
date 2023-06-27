'use strict';

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

function displayOverlay(html) {
	overlay.innerHTML = html;
	overlay.style.display = 'flex';
	adjustOverlay();
}

function hideOverlay() {
	overlay.innerHTML = '';
	overlay.style.display = 'none';
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
	hideOverlay();

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

//
//
//
//
//

// event listeners

function copyJoiningLinkEL() {
	navigator.clipboard.writeText(joiningLink);
}

function hideAndRemoveElForJoiningLink() {
	document
		.querySelector('.joining-link-btn')
		.removeEventListener('click', copyJoiningLinkEL);

	document.querySelector('.joining-link-div').classList.add('hidden');
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
