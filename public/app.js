const msgInp = document.querySelector('#msg');
const sendMsgBtn = document.querySelector('#sendMsg');
const messagesDiv = document.querySelector('.messages');

const addMsgToDOM = function (msg) {
	if (msg.length === 0) return;

	const msgDiv = document.createElement('div');
	const text = document.createTextNode(msg);
	msgDiv.appendChild(text);
	messagesDiv.appendChild(msgDiv);

	msgInp.value = '';
};

const poolId = new URLSearchParams(window.location.search).get('join');
const clientId = String(Date.now());

const socket = new WebSocket(
	`ws://localhost:1323/ws?poolId=${poolId}&clientId=${clientId}`
);

const connect = () => {
	socket.onopen = event => console.log('Successfully Connected', event);

	socket.onmessage = msg => {
		const data = JSON.parse(msg.data);
		// if (data.type !== 0)
		addMsgToDOM(msg.data);
	};

	socket.onclose = event => console.log('Socket Closed Connection: ', event);

	socket.onerror = error => console.log('Socket Error: ', error);
};

connect();

sendMsgBtn.addEventListener('click', e => {
	e.preventDefault();
	socket.send(msgInp.value);
});
