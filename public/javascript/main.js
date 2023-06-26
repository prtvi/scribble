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

let messageTypeMap;
let timeForEachWordInSeconds;
let timeForChoosingWordInSeconds;
let wordExpiryTimer;

// if overlay is active, disable scroll
if (overlay.style.display !== 'none') {
	document.documentElement.style.overflow = 'hidden';
	document.body.scroll = 'no';
}

// chat
document.querySelector('.send-msg').addEventListener('click', sendChatMsgBtnEL);

// show number of characters typed in chat box
const lenIndicator = document.querySelector('.input-wrapper span');
document
	.querySelector('.msg')
	.addEventListener(
		'input',
		e => (lenIndicator.textContent = e.target.value.length)
	);

// event listeners for drawing
window.addEventListener('load', () => {
	document.addEventListener('mousedown', startPainting);
	document.addEventListener('mouseup', stopPainting);
	document.addEventListener('mousemove', paint);
});

// copy joining link
document
	.querySelector('.joining-link-btn')
	.addEventListener('click', copyJoiningLinkEL);

// add event listener to start game button to start game
document
	.querySelector('.start-game-btn')
	.addEventListener('click', startGameEl);
