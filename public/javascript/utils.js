'use strict';

// utils for painting on canvas
const paintUtils = {
	coords: { x: 0, y: 0 },
	color: `#${clientColor}`,
	isPainting: false,
	hasGameStarted: false,
	isAllowedToPaint: false,
};

function wait(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

function getSecondsLeftFrom(futureTime) {
	const now = new Date().getTime();
	const diff = futureTime - now;
	return Math.round(diff / 1000);
}

function clearAllIntervals(...ids) {
	ids.forEach(i => clearInterval(i));
}

function displayScores(socketMessage) {
	const dataArr = JSON.parse(socketMessage.content);

	let html = `<div> <table>
	<tr>
		<th>Name</th>
		<th>Score</th>
	</tr>`;
	dataArr.forEach(
		item => (html += `<tr><td>${item.name}</td><td>${item.score}</td></tr>`)
	);
	html += `</table> </div>`;

	overlay.innerHTML = html;
	displayOverlay();

	clearAllIntervals(wordExpiryTimerIdG);
}

function disableSketching(socketMessage) {
	if (clientId !== socketMessage.currSketcherId) return;

	const painterUtilsDiv = document.querySelector('.painter-utils');
	const clearCanvasBtn = document.querySelector('.clear-canvas');

	paintUtils.isAllowedToPaint = false;

	// display painter utils div and remove EL
	painterUtilsDiv.classList.add('hidden');
	clearCanvasBtn.removeEventListener('click', requestCanvasClear);
}

// render clients

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

//  chat

function appendChatMsgToDOM(msg, formatColor) {
	// adds the msg into the DOM

	if (msg.length === 0 || msg === '') return;

	const messagesDiv = document.querySelector('.messages');

	const msgDiv = document.createElement('div');
	msgDiv.classList.add('message');

	const text = document.createElement('span');
	text.textContent = msg;
	text.style.color = formatColor || '#000';

	msgDiv.appendChild(text);
	messagesDiv.appendChild(msgDiv);

	msgDiv.scrollIntoView();

	document.querySelector('.msg').value = '';
	document.querySelector('.input-wrapper span').textContent = 0;
}

function sendChatMsgBtnEL(e) {
	// event listener to send chat message

	e.preventDefault();
	const msg = document.querySelector('.msg').value;

	if (msg.length === 0 || msg === '') return;

	// create string response object
	const socketMsg = {
		type: 3,
		typeStr: 'text_msg',
		content: msg,
		clientName,
		clientId,
		poolId,
	};

	// convert object to string to transmit
	sendViaSocket(socketMsg);
}

function renderRoundDetails(socketMessage) {
	document.querySelector(
		'.round span'
	).textContent = `Round: ${socketMessage.currRound}`;
}

function displayOverlay() {
	overlay.style.display = 'flex';
}

function hideOverlay() {
	overlay.style.display = 'none';
}

document.querySelector('.send-msg').addEventListener('click', sendChatMsgBtnEL);

// show number of characters typed in chat box
document.querySelector('.msg').addEventListener('input', function (e) {
	document.querySelector('.input-wrapper span').textContent =
		e.target.value.length;
});

// copy joining link
document
	.querySelector('.joining-link')
	.addEventListener('click', () => navigator.clipboard.writeText(joiningLink));
