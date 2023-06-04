'use strict';

const { canvas, ctx, overlay } = initCanvasAndOverlay();

// init socket connection and check game begin status
const socket = initSocket();
const startGameTimerId = gameStartTimer();

let wordExpiryTimerIdG;

function gameStartTimer() {
	// start game countdown to show user how much time is left for game to start
	return setInterval(
		() =>
			(document.querySelector('.loading').textContent =
				getSecondsLeftFrom(gameStartTime)),
		1000
	);
}

function startGame(socketMessage) {
	// called when socket receives message from server with type as 6
	if (!socketMessage.success) return;

	console.log('game started by server');

	paintUtils.hasGameStarted = true;
	clearAllIntervals(startGameTimerId);

	// hide the div and toggle paintUtils.has Game Started
	hideOverlay();
	document.querySelector('.joining-link-div').classList.add('hidden');
}

function beginClientSketchingFlow(socketMessage) {
	// initialise the time at which this word expires
	const currentWordExpiresAt = new Date(
		socketMessage.currWordExpiresAt
	).getTime();

	// start timer for the word expiry
	const wordExpiryTimerId = setInterval(async () => {
		const timeLeftDiv = document.querySelector('.time-left-for-word span');

		const secondsLeft = getSecondsLeftFrom(currentWordExpiresAt);
		timeLeftDiv.textContent = `Time: ${secondsLeft} seconds`;

		if (secondsLeft <= 0) clearAllIntervals(wordExpiryTimerId);
	}, 1000);

	const word = document.querySelector('.word span');
	const painterUtilsDiv = document.querySelector('.painter-utils');
	const clearCanvasBtn = document.querySelector('.clear-canvas');

	// for enabling drawing access if clientId matches
	if (clientId === socketMessage.currSketcherId) {
		paintUtils.isAllowedToPaint = true;

		// display the word
		word.textContent = socketMessage.currWord;

		// display painter utils div and add EL for clearing the canvas
		painterUtilsDiv.classList.remove('hidden');
		clearCanvasBtn.addEventListener('click', requestCanvasClear);
	} else {
		paintUtils.isAllowedToPaint = false;

		// show word length
		word.textContent = socketMessage.currWord.length;

		// display painter utils div and remove EL
		painterUtilsDiv.classList.add('hidden');
		clearCanvasBtn.removeEventListener('click', requestCanvasClear);
	}

	return wordExpiryTimerId;
}
