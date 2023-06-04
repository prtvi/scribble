'use strict';

// utils for painting on canvas
const paintUtils = {
	coords: { x: 0, y: 0 },
	color: `#${clientColor}`,
	isPainting: false,
	hasGameStarted: false,
	isAllowedToPaint: false,
};

function getCanvasSize() {
	const w = window.innerWidth;
	const cw = w - 10;
	const ch = cw / 1.5;

	return { w: cw, h: ch };
}

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

	clearAllIntervals(wordExpiryTimerIdG);
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
	allClients.forEach((n, i) => {
		const clientNameDiv = document.createElement('div');
		clientNameDiv.classList.add('member');

		const clientNum = document.createElement('span');
		clientNum.classList.add('member-num');

		const clientName = document.createElement('span');
		clientName.classList.add('member-name');
		clientName.style.color = `#${n.color}`;

		const clientScore = document.createElement('span');
		clientScore.classList.add('member-score');

		clientNum.textContent = `#${i + 1}`;
		clientName.textContent = n.name;
		clientScore.textContent = `${n.score} points`;

		clientNameDiv.appendChild(clientNum);
		clientNameDiv.appendChild(clientName);
		clientNameDiv.appendChild(clientScore);

		membersDiv.appendChild(clientNameDiv);
	});
}

//  chat

function appendChatMsgToDOM(msg) {
	// adds the msg into the DOM

	if (msg.length === 0 || msg === '') return;

	const messagesDiv = document.querySelector('.messages');

	const msgDiv = document.createElement('div');
	msgDiv.classList.add('message');

	const text = document.createElement('span');
	text.textContent = msg;

	msgDiv.appendChild(text);
	messagesDiv.appendChild(msgDiv);

	msgDiv.scrollIntoView();

	document.querySelector('.msg').value = '';
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

document.querySelector('.send-msg').addEventListener('click', sendChatMsgBtnEL);

// copy joining link
document
	.querySelector('.joining-link')
	.addEventListener('click', () => navigator.clipboard.writeText(joiningLink));

function renderRoundDetails(socketMessage) {
	document.querySelector(
		'.round-details'
	).textContent = `Round: ${socketMessage.currRound}`;
}

function toggleOverlay() {
	const overlay = document.querySelector('#overlay');
	overlay.classList.toggle('hidden');
}
