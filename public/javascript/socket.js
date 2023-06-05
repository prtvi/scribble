'use strict';

function initSocket() {
	// initialises socket connection and adds corresponding function handlers to the socket

	const wsUrl = `ws://${getDomain()}/ws?poolId=${poolId}&clientId=${clientId}&clientName=${clientName}&clientColor=${clientColor}`;
	const socket = new WebSocket(wsUrl);

	socket.onopen = () => {
		// on socket open success, get all clients and render them on UI
		console.log('Socket successfully connected!');
	};

	socket.onmessage = socketOnMessage;
	socket.onclose = socketOnClose;
	socket.onerror = error => console.log('Socket error', error);

	function getDomain() {
		// extract domain from url
		const url = window.location.href;
		const fi = url.indexOf('/');
		const li = url.lastIndexOf('/');
		const domain = url.slice(fi + 2, li);

		return domain;
	}

	function socketOnMessage(message) {
		// runs when a message is received on the socket conn, runs the corresponding functions depending on message type

		// parse json string into json object
		const socketMessage = JSON.parse(message.data);

		if (socketMessage.type !== 4)
			console.log(socketMessage.type, socketMessage.typeStr);

		switch (socketMessage.type) {
			case 1:
				if (socketMessage.clientId === clientId)
					// if the current clientId and the clientId from response match then
					appendChatMsgToDOM(
						`You joined the pool as ${socketMessage.clientName}!`,
						''
					);
				else
					appendChatMsgToDOM(
						`${socketMessage.clientName} has joined the pool!`,
						''
					);
				break;

			case 2:
				appendChatMsgToDOM(
					`${socketMessage.clientName} has left the pool!`,
					''
				);
				break;

			case 3:
				appendChatMsgToDOM(
					`${socketMessage.clientName}: ${socketMessage.content}`,
					''
				);
				break;

			case 31:
				appendChatMsgToDOM(socketMessage.content, '#00ff00');
				break;

			case 32:
				appendChatMsgToDOM(socketMessage.content, '#ffa500');
				break;

			case 4:
				displayImgOnCanvas(socketMessage);
				break;

			case 5:
				clearCanvas();
				break;

			case 6:
				renderClients(socketMessage.content);
				break;

			case 7:
				startGame(socketMessage);
				break;

			case 71:
				renderRoundDetails(socketMessage);
				break;

			case 8:
				wordExpiryTimerIdG = beginClientSketchingFlow(socketMessage);
				break;

			case 81:
				disableSketching(socketMessage);
				break;

			case 9:
				displayScores(socketMessage);
				break;

			default:
				break;
		}
	}

	function socketOnClose() {
		// on socket conn close, stop all timer or intervals
		console.log('Socket connection closed, stopping timers and timeouts!');
		clearAllIntervals(wordExpiryTimerIdG, startGameTimerId);
	}

	return socket;
}

function sendViaSocket(socketMsg) {
	/*  socket.readyState: int
			0 - connecting
			1 - open
			2 - closing
			3 - closed
	*/

	if (socket.readyState === socket.OPEN) socket.send(JSON.stringify(socketMsg));
	else {
		console.log(
			'0: connecting | 1: open | 2: closing | 3: closed, current state:',
			socket.readyState
		);

		clearAllIntervals(wordExpiryTimerIdG, startGameTimerId);
	}
}
