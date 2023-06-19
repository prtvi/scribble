'use strict';

function initSocket() {
	// initialises socket connection and adds corresponding function handlers to the socket

	const wsUrl = `ws://${getDomain()}/ws?poolId=${poolId}&clientId=${clientId}&clientName=${clientName}&clientColor=${clientColor}`;

	const socket = new WebSocket(wsUrl);

	socket.onopen = () => console.log('Socket successfully connected!');
	socket.onerror = error => console.log('Socket error', error);
	socket.onmessage = socketOnMessage;
	socket.onclose = socketOnClose;

	return socket;
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
					`You joined the room as <strong>${socketMessage.clientName}</strong>!`,
					''
				);
			else
				appendChatMsgToDOM(
					`<strong>${socketMessage.clientName}</strong> has joined the room!`,
					''
				);
			break;

		case 2:
			appendChatMsgToDOM(
				`<strong>${socketMessage.clientName}</strong> has left the room!`,
				''
			);
			break;

		case 3:
			appendChatMsgToDOM(
				`<strong>${socketMessage.clientName}</strong>: ${socketMessage.content}`,
				''
			);
			break;

		case 31:
			appendChatMsgToDOM(
				`${socketMessage.clientName} guessed the word!`,
				'#00ff00'
			);
			break;

		case 312:
			appendChatMsgToDOM(`Naughty @${socketMessage.clientName}`, '#ff0000');
			break;

		case 32:
			revealWordOnOverlayAndChat(socketMessage);
			break;

		case 33:
			showWordToChoose(socketMessage);
			break;

		case 35:
			showChoosingWordOnOverlay(socketMessage);
			break;

		case 4:
			displayImgOnCanvas(socketMessage);
			break;

		case 5:
		case 51:
			clearCanvas();
			break;

		case 6:
			renderClients(socketMessage.content);
			break;

		case 70:
			startGame(socketMessage);
			break;

		case 71:
			renderRoundDetails(socketMessage);
			break;

		case 8:
			wordExpiryTimer = beginClientSketchingFlow(socketMessage);
			break;

		case 88:
			wordExpiryTimer = showClientDrawing(socketMessage);
			break;

		case 81:
			disableSketchingTurnOver();
			break;

		case 82:
			showTimeUp();
			break;

		case 83:
			disableSketchingAllGuessed();
			break;

		case 84:
			showAllHaveGuessed();
			break;

		case 9:
			displayScores(socketMessage);
			break;

		case 10:
			makeMessageTypeMapGlobal(socketMessage);
			break;

		default:
			break;
	}
}

function socketOnClose() {
	// on socket conn close, stop all timer or intervals
	console.log('Socket connection closed, stopping timers and timeouts!');
	clearAllIntervals(startGameTimerId);
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

		clearAllIntervals(startGameTimerId);
	}
}
