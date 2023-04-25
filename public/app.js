const msgInp = document.querySelector('#msg');
const sendMsgBtn = document.querySelector('#sendMsg');
const messagesDiv = document.querySelector('.messages');

const main = document.querySelector('.main');

const addMsgToDOM = function (msg) {
	if (msg.length === 0) return;

	const msgDiv = document.createElement('div');
	const text = document.createTextNode(msg);
	msgDiv.appendChild(text);
	messagesDiv && messagesDiv.appendChild(msgDiv);

	msgInp.value = '';
};

const connect = async () => {
	const wsUrl = `ws://localhost:1323/ws?poolId=${poolId}&clientId=${clientId}&clientName=${clientName}`;

	console.log('connecting socket', wsUrl);

	const socket = new WebSocket(wsUrl);

	socket.onopen = () => console.log('Successfully Connected');
	socket.onmessage = msg => addMsgToDOM(msg.data);
	socket.onclose = () => console.log('Socket Closed Connection');
	socket.onerror = error => console.log('Socket Error', error);

	sendMsgBtn.addEventListener('click', e => {
		e.preventDefault();
		socket.send(msgInp.value);
	});
};

const poolId = document.getElementsByName('poolId')[0].value;
const clientName = document.getElementsByName('clientName')[0].value;
const clientId = String(Date.now());

connect();
