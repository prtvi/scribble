'use strict';

console.log('game will be started by', clientName);

// add event listener to start game button to start game
const startGameBtn = document.querySelector('.start-game-btn');
startGameBtn.addEventListener('click', requestStartGameEL);

// start game after this timeout
const startGameAfterTimeoutId = setTimeout(
	requestStartGameEL,
	getSecondsLeftFrom(gameStartTime) * 1000
);

function requestStartGameEL() {
	// runs when the game starts, makes socket conn call to server to start the game
	// clear the countdown timers
	clearAllIntervals(startGameAfterTimeoutId);

	// generate response and send
	const responseMsg = {
		type: 7,
		content: 'start the game bro!',
		clientId,
		clientName,
		poolId,
	};

	sendViaSocket(responseMsg);
}
