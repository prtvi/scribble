const msgInp = document.querySelector('#msg');
const sendMsgBtn = document.querySelector('#sendMsg');
const messagesDiv = document.querySelector('.messages');

const socket = new WebSocket('ws://localhost:8080/ws');

const addMsgToDOM = function (msg) {
	if (msg.length === 0) return;

	const msgDiv = document.createElement('div');
	const text = document.createTextNode(msg);
	msgDiv.appendChild(text);
	messagesDiv.appendChild(msgDiv);

	msgInp.value = '';
};

const connect = () => {
	socket.onopen = () => console.log('Successfully Connected');

	socket.onmessage = msg => addMsgToDOM(msg.data);

	socket.onclose = event => console.log('Socket Closed Connection: ', event);

	socket.onerror = error => console.log('Socket Error: ', error);
};

const sendMsg = msg => socket.send(msg);

connect();

sendMsgBtn.addEventListener('click', e => {
	e.preventDefault();
	sendMsg(msgInp.value);
});
