const msgInp = document.querySelector('#msg');
const sendMsgBtn = document.querySelector('#sendMsg');
const messagesDiv = document.querySelector('.messages');

const addMsgToDOM = function (msg) {
	// adds the content into the DOM
	if (msg.length === 0) return;

	const msgDiv = document.createElement('div');
	const text = document.createTextNode(msg);
	msgDiv.appendChild(text);
	messagesDiv && messagesDiv.appendChild(msgDiv);

	msgInp.value = '';
};

const connect = async () => {
	// creates a socket connection to the server and renders events

	// establish socket connection using the domain, poolId, clientId and the clientName
	const wsUrl = `ws://${domain}/ws?poolId=${poolId}&clientId=${clientId}&clientName=${clientName}`;
	console.log('connecting socket', wsUrl);
	const socket = new WebSocket(wsUrl);

	// on socket open event
	socket.onopen = () => console.log('Successfully Connected');

	// on socket, received messages from backend
	socket.onmessage = msg => {
		// parse json string into json object
		const msgJson = JSON.parse(msg.data);
		console.log('received:', msgJson);

		// if message type is 1 === CONNECTED
		if (msgJson.type === 1) {
			// if the current clientName and the clientName from response match then
			if (msgJson.clientName === clientName)
				addMsgToDOM(`You joined the pool as ${clientName}!`);
			// else
			else addMsgToDOM(`${msgJson.clientName} has joined the pool!`);
			// if message type is 2 === DISCONNECTED
		} else if (msgJson.type == 2) {
			addMsgToDOM(`${msgJson.clientName} has left the pool!`);
			// if message type is 3 === JSON/string data
		} else if (msgJson.type === 3) {
			addMsgToDOM(`${msgJson.clientName}: ${msgJson.content}`);
		}
	};

	// on socket close event
	socket.onclose = () => console.log('Socket Closed Connection');
	// on socket error event
	socket.onerror = error => console.log('Socket Error', error);

	// add event listener to the button only if socket connection is established
	sendMsgBtn.addEventListener('click', e => {
		e.preventDefault();

		// create response object
		const responseMsg = {
			type: 3,
			content: msgInp.value,
			clientName,
			clientId,
		};
		console.log('sending:', responseMsg);

		// convert object to string to transmit
		socket.send(JSON.stringify(responseMsg));
	});
};

// extract domain from url
const url = window.location.href;
const fi = url.indexOf('/');
const li = url.lastIndexOf('/');
const domain = url.slice(fi + 2, li);

// get the poolId and clientName from the form
// generate clientId using timestamp
const poolId = document.getElementsByName('poolId')[0].value;
const clientName = document.getElementsByName('clientName')[0].value;
const clientId = String(Date.now());

connect();

// draw on canvas code

const canvas = document.querySelector('#canv');
const ctx = canvas.getContext('2d');

const coord = { x: 0, y: 0 };
let paint = false;

const getPosition = function (event) {
	coord.x = event.clientX - canvas.offsetLeft;
	coord.y = event.clientY - canvas.offsetTop;
};

const startPainting = function (event) {
	paint = true;
	getPosition(event);
};

const stopPainting = function () {
	paint = false;
};

const sketch = function (event) {
	if (!paint) return;

	ctx.beginPath();

	ctx.lineWidth = 5;
	ctx.lineCap = 'round';
	ctx.strokeStyle = 'green';

	ctx.moveTo(coord.x, coord.y);

	getPosition(event);

	ctx.lineTo(coord.x, coord.y);
	ctx.stroke();
};

window.addEventListener('load', () => {
	document.addEventListener('mousedown', startPainting);
	document.addEventListener('mouseup', stopPainting);
	document.addEventListener('mousemove', sketch);
});
