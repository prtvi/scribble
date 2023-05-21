'use strict';

console.log('game will be started by', clientName);

// add event listener to start game button to start game
const startGameBtn = document.querySelector('.start-game-btn');
startGameBtn.addEventListener('click', () => {
	// runs when the game starts, makes socket conn call to server to start the game
	// generate response and send

	const responseMsg = {
		type: 7,
		content: 'start the game bro!',
		clientId,
		clientName,
		poolId,
	};

	sendViaSocket(responseMsg);
});
