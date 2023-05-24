'use strict';

console.log('game will be started by me');

// add event listener to start game button to start game
// makes socket conn call to server to start the game
document.querySelector('.start-game-btn').addEventListener('click', () => {
	const socketMsg = {
		type: 7,
		typeStr: 'start_game',
		clientId,
		clientName,
		poolId,
	};

	console.log('requesting game to start');
	sendViaSocket(socketMsg);
});
