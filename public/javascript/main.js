'use strict';

// canvas, canvas ctx and overlay init
const { canvas, ctx, overlay } = initCanvasAndOverlay();

// utils for painting on canvas
const paintUtils = {
	coords: { x: 0, y: 0 },
	isPainting: false,
	hasGameStarted: false,
	isAllowedToPaint: false,
};

let messageTypeMap;
let timeForEachWordInSeconds;
let timeForChoosingWordInSeconds;
let wordExpiryTimer;
let allowLogs;

// init socket connection and check game begin status
const socket = initSocket();
initGlobalEventListeners();
