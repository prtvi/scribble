'use strict';

// canvas, canvas ctx and overlay init
const { canvas, ctx, overlay } = initCanvasAndOverlay();

// utils for painting on canvas
const paintUtils = {
	coords: { x: 0, y: 0 },
	color: `#${clientColor}`,
	isPainting: false,
	hasGameStarted: false,
	isAllowedToPaint: false,
};

// init socket connection and check game begin status
const socket = initSocket();
const startGameTimerId = gameStartTimer();

let wordExpiryTimerIdG;

// chat
document.querySelector('.send-msg').addEventListener('click', sendChatMsgBtnEL);

// show number of characters typed in chat box
document.querySelector('.msg').addEventListener('input', function (e) {
	document.querySelector('.input-wrapper span').textContent =
		e.target.value.length;
});

// event listeners for drawing
window.addEventListener('load', () => {
	document.addEventListener('mousedown', startPainting);
	document.addEventListener('mouseup', stopPainting);
	document.addEventListener('mousemove', paint);
});

// copy joining link
document
	.querySelector('.joining-link')
	.addEventListener('click', () => navigator.clipboard.writeText(joiningLink));

// add event listener to start game button to start game
document.querySelector('.start-game-btn').addEventListener('click', () => {
	const socketMsg = {
		type: 7,
		typeStr: 'start_game',
		clientId,
		clientName,
		poolId,
	};

	sendViaSocket(socketMsg);
});
