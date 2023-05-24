'use strict';

const canvas = document.querySelector('.canv');
const ctx = canvas.getContext('2d');

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
	const startGameDiv = document.querySelector('.start-game');
	startGameDiv.classList.add('hidden');
}

function beginClientSketchingFlow(socketMessage) {
	console.table(socketMessage);

	// initialise the time at which this word expires
	const currentWordExpiresAt = new Date(
		socketMessage.currWordExpiresAt
	).getTime();

	// start timer for the word expiry
	const wordExpiryTimerId = setInterval(async () => {
		const timeLeftDiv = document.querySelector('.time-left-for-word');
		timeLeftDiv.classList.remove('hidden');

		const secondsLeft = getSecondsLeftFrom(currentWordExpiresAt);
		timeLeftDiv.querySelector('span').textContent = secondsLeft;

		if (secondsLeft <= 0) {
			clearAllIntervals(wordExpiryTimerId);
			console.log('timer for word cleared');
		}
	}, 1000);

	// for enabling drawing access if clientId matches
	if (clientId === socketMessage.currSketcherId) {
		paintUtils.isAllowedToPaint = true;

		// display the word by unhiding the painter-utils div
		document.querySelector('.painter-utils').classList.remove('hidden');
		document.querySelector('.your-word').textContent = socketMessage.currWord;

		// add EL for clearing the canvas
		document
			.querySelector('.clear-canvas')
			.addEventListener('click', requestCanvasClear);
	} else {
		paintUtils.isAllowedToPaint = false;
		document.querySelector('.painter-utils').classList.add('hidden');
		document.querySelector('.your-word').textContent = '';

		document
			.querySelector('.clear-canvas')
			.removeEventListener('click', requestCanvasClear);
	}

	return wordExpiryTimerId;
}
