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

// get all clients and render

function getAllClientsEL() {
	// makes a socket connection call to request client info list
	const responseMsg = {
		type: 6,
		content: '',
		clientId,
		clientName,
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

		clientName.innerHTML = `#${i + 1} ${n.name}${
			n.score === 0 ? '' : `: ${n.score} points`
		}`;
		clientName.style.color = `#${n.color}`;
		clientNameHolder.appendChild(clientName);

		membersDiv.appendChild(clientNameHolder);
	});
}

//  chat

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
		poolId,
	};

	// convert object to string to transmit
	sendViaSocket(responseMsg);
}

const renderClientsTimerId = setInterval(getAllClientsEL, 15 * 1000);
document.querySelector('.send-msg').addEventListener('click', sendChatMsgBtnEL);
