const msgInp = document.querySelector('#msg');
const sendMsgBtn = document.querySelector('#sendMsg');
const messagesDiv = document.querySelector('.messages');

const addMsgToDOM = function (msg) {
	if (msg.length === 0) return;

	const msgDiv = document.createElement('div');
	const text = document.createTextNode(msg);
	msgDiv.appendChild(text);
	messagesDiv && messagesDiv.appendChild(msgDiv);

	msgInp.value = '';
};

const connect = async () => {
	const wsUrl = `ws://${domain}/ws?poolId=${poolId}&clientId=${clientId}&clientName=${clientName}`;

	console.log('connecting socket', wsUrl);
	const socket = new WebSocket(wsUrl);

	socket.onopen = () => console.log('Successfully Connected');

	socket.onmessage = msg => {
		const msgJson = JSON.parse(msg.data);
		console.log('received:', msgJson);

		addMsgToDOM(`${msgJson.clientName}: ${msgJson.content}`);
	};

	socket.onclose = () => console.log('Socket Closed Connection');
	socket.onerror = error => console.log('Socket Error', error);

	sendMsgBtn.addEventListener('click', e => {
		e.preventDefault();

		const responseMsg = {
			type: 0,
			content: msgInp.value,
			clientName,
			clientId,
		};
		console.log('sending:', responseMsg);

		socket.send(JSON.stringify(responseMsg));
	});
};

const url = window.location.href;
const li = url.lastIndexOf('/');
const fi = url.indexOf('/');
const domain = url.slice(fi + 2, li);

const poolId = document.getElementsByName('poolId')[0].value;
const clientName = document.getElementsByName('clientName')[0].value;
const clientId = String(Date.now());

connect();
