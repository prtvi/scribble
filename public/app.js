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

const connect = async () => {
	const res = await fetch(`/check-pool?poolId=${poolId}&clientId=${clientId}`);
	const data = await res.json();
	console.log(data);

	if (data.code !== 200) {
		messagesDiv.innerHTML = `<p>${data.msg}</p>`;
		sendMsgBtn.disabled = true;
		return;
	}

	const socket = new WebSocket(
		`ws://localhost:1323/ws?poolId=${poolId}&clientId=${clientId}`
	);

	socket.onopen = () => console.log('Successfully Connected');
	socket.onmessage = msg => {
		// const data = JSON.parse(msg.data);
		// if (data.type !== 0)
		addMsgToDOM(msg.data);
	};

	socket.onclose = () => console.log('Socket Closed Connection');
	socket.onerror = error => console.log('Socket Error', error);

	sendMsgBtn.addEventListener('click', e => {
		e.preventDefault();
		socket.send(msgInp.value);
	});
};

const poolId = new URLSearchParams(window.location.search).get('join');
const clientId = String(Date.now());

connect();
