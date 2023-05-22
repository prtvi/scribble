'use strict';

console.log('game will be started by me');

// add event listener to start game button to start game
// makes socket conn call to server to start the game
document.querySelector('.start-game-btn').addEventListener('click', () => {
	const responseMsg = {
		type: 7,
		content: 'start the game bro!',
		clientId,
		clientName,
		poolId,
	};

	console.log('requesting game to start');
	sendViaSocket(responseMsg);
});
