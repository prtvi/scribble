'use strict';

// 1, 2, 3, 31, 32
function appendChatMsgToDOM(msg, formatColor) {
	// adds the msg into the DOM

	if (msg.length === 0 || msg === '') return;

	const messagesDiv = document.querySelector('.messages');

	const msgDiv = document.createElement('div');
	msgDiv.classList.add('message');

	const text = document.createElement('span');
	text.textContent = msg;
	text.style.color = formatColor || '#000';

	msgDiv.appendChild(text);
	messagesDiv.appendChild(msgDiv);

	msgDiv.scrollIntoView();

	document.querySelector('.msg').value = '';
	document.querySelector('.input-wrapper span').textContent = 0;
}

// 33
function showWordToChoose(socketMessage) {
	if (clientId === socketMessage.currSketcherId) {
		const words = JSON.parse(socketMessage.content);

		let html = `<div><p>Choose a word to draw</p>`;
		words.forEach(w => (html += `<span class="word-option">${w}</span>`));
		html += `</div>`;

		overlay.innerHTML = html;
		displayOverlay();

		overlay.querySelector('div').addEventListener('click', function (e) {
			const chosenWord = e.target.textContent.trim();
			if (!words.includes(chosenWord)) return;

			const socketMsg = {
				type: 34,
				typeStr: 'chosen_word',
				content: chosenWord,
				clientName,
				clientId,
				poolId,
			};

			sendViaSocket(socketMsg);
		});
	} else {
		overlay.innerHTML = `<div>${socketMessage.currSketcherName} is choosing a word!</div>`;
		displayOverlay();
	}
}

// 4
function displayImgOnCanvas(socketMessage) {
	// display image data on canvas

	if (clientId === socketMessage.currSketcherId) return;

	var img = new Image();
	// scale up/down canvas data based on current canvas size using outer bounds
	img.onload = () => ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
	img.setAttribute('src', socketMessage.content);
}

// 5
function clearCanvas() {
	hideOverlay();
	ctx.clearRect(0, 0, canvas.width, canvas.height);
}

// 6
function renderClients(allClients) {
	// called when the socket conn receives a message from server as type 6

	if (allClients.length === 0) return;

	const membersDiv = document.querySelector('.members');
	membersDiv.innerHTML = '';

	// parse array of objects into json
	allClients = JSON.parse(allClients);

	// render
	allClients.forEach((n, i) => membersDiv.appendChild(getClientNameDiv(n, i)));
}

// 7
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

// 71
function renderRoundDetails(socketMessage) {
	document.querySelector(
		'.round span'
	).textContent = `Round: ${socketMessage.currRound}`;

	overlay.innerHTML = `<div>Round: ${socketMessage.currRound}</div>`;
	displayOverlay();
}

// 8
function beginClientSketchingFlow(socketMessage) {
	hideOverlay();

	// initialise the time at which this word expires
	const currentWordExpiresAt = new Date(
		socketMessage.currWordExpiresAt
	).getTime();

	// start timer for the word expiry
	const wordExpiryTimerId = setInterval(() => {
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
		// show word length
		word.textContent = `${socketMessage.currWord.length} characters`;
	}

	return wordExpiryTimerId;
}

// 81
function disableSketching(socketMessage) {
	if (clientId !== socketMessage.currSketcherId) return;

	const painterUtilsDiv = document.querySelector('.painter-utils');
	const clearCanvasBtn = document.querySelector('.clear-canvas');

	paintUtils.isAllowedToPaint = false;

	// display painter utils div and remove EL
	painterUtilsDiv.classList.add('hidden');
	clearCanvasBtn.removeEventListener('click', requestCanvasClear);
}

// 9
function displayScores(socketMessage) {
	const dataArr = JSON.parse(socketMessage.content);

	let html = `<div> <table>
	<tr>
		<th>Name</th>
		<th>Score</th>
	</tr>`;
	dataArr.forEach(
		item => (html += `<tr><td>${item.name}</td><td>${item.score}</td></tr>`)
	);
	html += `</table> </div>`;

	overlay.innerHTML = html;
	displayOverlay();

	clearAllIntervals(wordExpiryTimerIdG);
}
