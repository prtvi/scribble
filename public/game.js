'use strict';

// -------------------------------- UTILS --------------------------------

function log(...args) {
	if (allowLogs) console.log(...args);
}

function wait(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

function clearAllIntervals(...ids) {
	ids.forEach(i => clearInterval(i));
}

/**
 * Get time left from the input time in seconds
 * @param {Date.getTime} futureTime
 * @returns time left in seconds from the given time
 */
function getSecondsLeftFrom(futureTime) {
	const now = new Date().getTime();
	const diff = futureTime - now;
	return Math.round(diff / 1000);
}

// -------------------------------- AVATAR UTILS --------------------------------

/**
 * Sets the background position of the given element
 * @param {HTMLElement} element element whose backgroud position will be set
 * @param {Number} x x coordinate on the atlas image
 * @param {Number} y y coordinate on the atlas image
 * @param {Number} scale set the scale
 */
function setBgPosition(element, x, y, scale) {
	const offset = scale * 48;
	element.style.backgroundPositionX = `-${x * offset}px`;
	element.style.backgroundPositionY = `-${y * offset}px`;
}

/**
 * Returns the current position of window over the avatar atlas
 * @param {String} pos element position as a string: 34px
 * @param {*} scale scale of avatar
 * @returns current position of window over the avatar atlas
 */
function getCurrPosition(pos, scale) {
	const offset = scale * 48;
	const lastIdx = pos.lastIndexOf('px');
	return Math.abs(+pos.slice(0, lastIdx)) / offset;
}

/**
 * Returns a random element from input array
 * @param {Array} arr input array
 * @returns a random element from input array
 */
function getRandomValue(arr) {
	return arr[Math.floor(Math.random() * arr.length)];
}

/**
 * Returns randomised avatar configuration
 * @returns randomised avatar configuration
 */
function getRandomizedAvatarCoords() {
	return {
		color: getRandomValue(validCoordsForAvatarAtlas.color),
		eyes: getRandomValue(validCoordsForAvatarAtlas.eyes),
		mouth: getRandomValue(validCoordsForAvatarAtlas.mouth),
	};
}

/**
 * Initialises valid co-ordinates over which the avatar can be rendered
 * @param {Object} prop one of the three properties of the avatar - mouth, color and eyes
 * @returns valid co-ordinates over which the avatar can be chosen over
 */
function initValidCoordsForAvatarAtlas(prop) {
	const rows = 10;
	const columns = 10;
	let flag = false;

	const coords = [];

	for (let row = 0; row < rows; row++) {
		for (let col = 0; col < columns; col++) {
			coords.push({ x: col, y: row });

			if (col === prop.x && row === prop.y) {
				flag = true;
				break;
			}
		}

		if (flag) break;
	}

	return coords;
}

/**
 * Sets the avatarConfig.isOwner field as true if the url contains isOwner=true
 */
function setIfOwner() {
	const urlParams = new URLSearchParams(location.search);
	if (urlParams.get('isOwner') === 'true') avatarConfig.isOwner = true;
}

/**
 * Renders row-avatars on home page
 */
function renderRowAvatars() {
	const avatarRow = document.querySelector('.avatar-row');
	const randomWinner = Math.round(Math.random() * 7);

	for (let i = 0; i < 8; i++) {
		const coords = getRandomizedAvatarCoords();

		const ac = {
			color: coords.color,
			eyes: coords.eyes,
			mouth: coords.mouth,
			isOwner: false,
			isCrowned: false,
		};
		if (randomWinner === i) ac.isCrowned = true;

		avatarRow.appendChild(getAvatarDom(ac, rowAvatarScale, 'row-avatar'));
	}
}

/**
 * Renders customizable avatar and adds event listeners for the same
 */
function renderCustomisableAvatar() {
	randomizeAvatar();

	const colorLeft = document.querySelector('.avc-btn.color-left');
	const colorRight = document.querySelector('.avc-btn.color-right');
	colorLeft.name = 'color';
	colorRight.name = 'color';

	const eyesLeft = document.querySelector('.avc-btn.eyes-left');
	const eyesRight = document.querySelector('.avc-btn.eyes-right');
	eyesLeft.name = 'eyes';
	eyesRight.name = 'eyes';

	const mouthLeft = document.querySelector('.avc-btn.mouth-left');
	const mouthRight = document.querySelector('.avc-btn.mouth-right');
	mouthLeft.name = 'mouth';
	mouthRight.name = 'mouth';

	document
		.querySelector('.randomize')
		.addEventListener('click', randomizeAvatar);

	[colorLeft, eyesLeft, mouthLeft].forEach(ele =>
		ele.addEventListener('click', leftElForCustomizableAvatar)
	);

	[colorRight, eyesRight, mouthRight].forEach(ele =>
		ele.addEventListener('click', rightElForCustomizableAvatar)
	);
}

/**
 * Event listener for customising avatar, left direction
 * @param {Event} e click event for customising avatar
 */
function leftElForCustomizableAvatar(e) {
	const rows = 10;
	const name = e.currentTarget.name;
	const elem = document.querySelector(`.avatar.custom-avatar .${name}`);

	if (elem.style.backgroundPositionX === '')
		elem.style.backgroundPositionX = '0px';

	avatarConfig[name].x =
		getCurrPosition(
			elem.style.backgroundPositionX,
			customizableAvatarScale
		) - 1;

	if (avatarConfig[name].x < 0 && avatarConfig[name].y > 0) {
		avatarConfig[name].y -= 1;
		avatarConfig[name].x = rows - 1;
	}

	if (avatarConfig[name].x < 0 && avatarConfig[name].y === 0) {
		avatarConfig[name].x = boundariesForAvatarAtlas[name].x;
		avatarConfig[name].y = boundariesForAvatarAtlas[name].y;
	}

	setBgPosition(
		elem,
		avatarConfig[name].x,
		avatarConfig[name].y,
		customizableAvatarScale
	);
	setIfOwner();
	saveToLocalStorage('avatarConfig', avatarConfig);
}

/**
 * Event listener for customising avatar, right direction
 * @param {Event} e click event for customising avatar
 */
function rightElForCustomizableAvatar(e) {
	const rows = 10;
	const name = e.currentTarget.name;
	const elem = document.querySelector(`.avatar.custom-avatar .${name}`);

	if (elem.style.backgroundPositionX === '')
		elem.style.backgroundPositionX = '0px';

	avatarConfig[name].x =
		getCurrPosition(
			elem.style.backgroundPositionX,
			customizableAvatarScale
		) + 1;

	if (avatarConfig[name].x >= rows) {
		avatarConfig[name].y += 1;
		avatarConfig[name].x = 0;
	}

	if (
		avatarConfig[name].x > boundariesForAvatarAtlas[name].x &&
		avatarConfig[name].y === boundariesForAvatarAtlas[name].y
	) {
		avatarConfig[name].x = 0;
		avatarConfig[name].y = 0;
	}

	setBgPosition(
		elem,
		avatarConfig[name].x,
		avatarConfig[name].y,
		customizableAvatarScale
	);
	setIfOwner();
	saveToLocalStorage('avatarConfig', avatarConfig);
}

/**
 * Randomise the customisable avatar
 */
function randomizeAvatar() {
	const coords = getRandomizedAvatarCoords();

	avatarConfig.color = coords.color;
	avatarConfig.eyes = coords.eyes;
	avatarConfig.mouth = coords.mouth;

	const color = document.querySelector('.avatar.custom-avatar .color');
	const eyes = document.querySelector('.avatar.custom-avatar .eyes');
	const mouth = document.querySelector('.avatar.custom-avatar .mouth');

	setBgPosition(
		color,
		coords.color.x,
		coords.color.y,
		customizableAvatarScale
	);
	setBgPosition(eyes, coords.eyes.x, coords.eyes.y, customizableAvatarScale);
	setBgPosition(
		mouth,
		coords.mouth.x,
		coords.mouth.y,
		customizableAvatarScale
	);

	setIfOwner();
	saveToLocalStorage('avatarConfig', avatarConfig);
}

/**
 * Generates the DOM for an avatar
 * @param {AvatarConfig Object} avatarConfig avatarConfig for rendering an avatar
 * @param {...String} classNames classnames to be added to avatar
 * @returns the DOM for the avatar
 */
function getAvatarDom(avatarConfig, scale, ...classNames) {
	const playerAvatar = document.createElement('div');
	playerAvatar.classList.add('avatar', ...classNames);

	const pColor = document.createElement('div');
	pColor.classList.add('color');
	setBgPosition(pColor, avatarConfig.color.x, avatarConfig.color.y, scale);

	const pEyes = document.createElement('div');
	pEyes.classList.add('eyes');
	setBgPosition(pEyes, avatarConfig.eyes.x, avatarConfig.eyes.y, scale);

	const pMouth = document.createElement('div');
	pMouth.classList.add('mouth');
	setBgPosition(pMouth, avatarConfig.mouth.x, avatarConfig.mouth.y, scale);

	// TODO: add this to the beginning of player card
	// const pOwner = document.createElement('div');
	// pOwner.classList.add('owner');
	// if (avatarConfig.isOwner) pOwner.classList.add('active');

	const pCrowned = document.createElement('div');
	pCrowned.classList.add('crowned');
	if (avatarConfig.isCrowned) pCrowned.classList.add('active');

	playerAvatar.appendChild(pColor);
	playerAvatar.appendChild(pEyes);
	playerAvatar.appendChild(pMouth);
	// playerAvatar.appendChild(pOwner);
	playerAvatar.appendChild(pCrowned);

	return playerAvatar;
}

/**
 * Generates DOM for displaying player stats
 * @param {Object} playerInfo object defining the player stats
 * @param {Number} iteration #number on the UI
 * @returns player DOM with rank, name, score and avatar
 */
function getPlayerCardDom(playerInfo, iteration) {
	// player div
	const playerCard = document.createElement('div');
	playerCard.classList.add('player-card');

	// player num span
	const playerNum = document.createElement('span');
	playerNum.classList.add('num');
	playerNum.textContent = `#${iteration + 1}`;

	// name span
	const playerName = document.createElement('span');
	playerName.classList.add('name');

	playerName.textContent = playerInfo.name;
	if (clientId === playerInfo.id) {
		playerName.classList.add('self');
		playerName.textContent += ' (you)';
	}

	// score span
	const playerScore = document.createElement('span');
	playerScore.classList.add('score');
	playerScore.textContent = `${playerInfo.score} points`;

	// player name and score div
	const playerNameAndScore = document.createElement('div');
	playerNameAndScore.appendChild(playerName);
	playerNameAndScore.appendChild(playerScore);

	// append everything to player div
	playerCard.appendChild(playerNum);
	playerCard.appendChild(playerNameAndScore);

	if (playerInfo.isSketching) {
		const playerIsSketching = document.createElement('div');
		playerIsSketching.classList.add('player-sketching');
		const isSketchingImg = document.createElement('img');
		isSketchingImg.src = 'public/assets/images/pen.gif';
		isSketchingImg.width = 36 * playerCardAvatarScale;
		playerIsSketching.appendChild(isSketchingImg);
		playerCard.appendChild(playerIsSketching);
	}

	// player avatar
	const playerAvatar = getAvatarDom(
		playerInfo.avatarConfig,
		playerCardAvatarScale,
		'player-card-avatar'
	);
	playerCard.appendChild(playerAvatar);

	return playerCard;
}

/**
 * Sends the chat message to the server
 * @param {Event} e event that triggers sending text message
 * @returns void
 */
function sendChatMsgBtnEL(e) {
	e.preventDefault();

	const msg = document.querySelector('.msg').value.trim();
	if (msg.length === 0 || msg === '') return;

	const socketMsg = {
		type: 3,
		content: msg,
		clientName,
		clientId,
		poolId,
	};

	sendViaSocket(socketMsg);
}

/**
 * Returns the HTML DOM for the overlay, for text only
 * @param {String} overlayHeading text that will be shown on the overlay
 * @returns String, html dom
 */
function getOverlayContentDomWithHeading(overlayHeading) {
	const overlayContent = document.createElement('div');
	overlayContent.classList.add('overlay-content');

	const heading = document.createElement('p');
	heading.classList.add('overlay-heading');
	heading.textContent = overlayHeading;

	const headingContainer = document.createElement('div');
	headingContainer.appendChild(heading);

	overlayContent.appendChild(headingContainer);

	return overlayContent;
}

/**
 * Renders the overlay with the given HTML string
 * @param {HTMLElement} dom the HTML that will be shown on the overlay
 */
function displayOverlay(dom) {
	// display overlay after some delay to render fade in animation
	// if event listeners are to be added to the given html, then use the same timeout to attach the event listeners

	overlay.style.opacity = 0;
	setTimeout(function () {
		overlay.innerHTML = '';
		overlay.appendChild(dom);
		overlay.style.display = 'flex';
		overlay.style.opacity = 1;
		adjustOverlay();
	}, overlayFadeInAnimationDuration);
}

/**
 * Hides the overlay
 */
function hideOverlay() {
	// render hiding animation using timeout
	overlay.style.opacity = 1;
	setTimeout(() => {
		overlay.style.opacity = 0;
		overlay.innerHTML = '';
	}, overlayFadeInAnimationDuration);

	// change overlay display property to none after the animation
	setTimeout(() => (overlay.style.display = 'none'), 1000);
}

/**
 * Adjusts the overlay position, currently used at scroll event
 */
function adjustOverlay() {
	const cc = document.querySelector('.canvas-container');
	overlay.style.top = `${cc.offsetTop}px`;
	overlay.style.height = `${cc.offsetHeight}px`;
	overlay.style.width = `${cc.offsetWidth}px`;
	overlay.style.left = `${cc.offsetLeft}px`;
}

/**
 * Sets the local storage with the given values
 * @param {String} key unique key to set the local storage
 * @param {any} value any value to be stored
 */
function saveToLocalStorage(key, value) {
	window.localStorage.setItem(key, JSON.stringify(value));
}

/**
 * Returns the locally stored content based on the key provided
 * @param {String} key unique key to get the locally saved content
 * @returns string value of stored content
 */
function getFromLocalStorage(key) {
	return window.localStorage.getItem(key);
}

/**
 * Disable player sketching by hiding the paint utils and removing the event listeners
 */
function disableSketching() {
	// get the paint utils, clear canvav btn & undo btn
	const paintUtilsDiv = document.querySelector('.paint-utils');
	const clearCanvasBtn = document.querySelector('.pu.clear');
	const undoBtn = document.querySelector('.pu.undo');

	// get the strokes and brush Stroke Selected divs
	const strokes = document.querySelector('.strokes');
	const strokeSelected = document.querySelector('.pu.stroke-selected');

	// get the colors and color selected divs
	const colors = document.querySelector('.colors');
	const colorSelected = document.querySelector('.pu.color-selected');

	// hide paint utils div and disable painting access
	paintUtils.isAllowedToPaint = false;
	paintUtilsDiv.classList.add('hidden');

	// remove ELs for clear canvas & undo btn
	clearCanvasBtn.removeEventListener('click', requestCanvasClear);
	undoBtn.removeEventListener('click', undo);

	// remove ELs for selecting brush stroke
	strokes.removeEventListener('click', selectStrokeEL);
	strokeSelected.removeEventListener('click', openStrokeSelectOptions);

	// remove ELs for selecting brush color
	colors.removeEventListener('click', selectColorEL);
	colorSelected.removeEventListener('click', openColorSelectOptions);

	// clear the paths
	paintUtils.points = [];
	paintUtils.paths = [];

	// undo overflow:hidden
	document.querySelector('*').style.overflow = '';
}

/**
 * EVENT: 86
 * Begins the player sketching process by starting the timer
 * @param {Object} socketMessage
 * @returns Number, the timer id from the setInterval function
 */
function beginClientSketchingFlowInit(socketMessage) {
	// initialise the time at which this word expires
	const currentWordExpiresAt = new Date(
		socketMessage.currWordExpiresAt
	).getTime();

	const timeLeftSpan = document.querySelector('.timer span');
	setGbTimerStat(timeForEachWordInSeconds);
	return runTimer(timeLeftSpan, currentWordExpiresAt);
}

/**
 * Runs a timer
 * @param {HTMLElement} timerElement element on which the timer will display the time left
 * @param {Date.getTime} timeoutAt timeStamp at which the timer runs out
 * @returns Timer ID
 */
function runTimer(timerElement, timeoutAt) {
	const countdownTimer = setInterval(function () {
		const secondsLeft = getSecondsLeftFrom(timeoutAt);
		if (secondsLeft > -1) timerElement.textContent = `${secondsLeft}s`;
		else clearInterval(countdownTimer);
	}, 1000);
	return countdownTimer;
}

/**
 * Remove event listeners on game start and hide joining link btn
 */
function removeEventListenersOnGameStart() {
	const isOwner = JSON.parse(getFromLocalStorage('avatarConfig')).isOwner;
	if (isOwner) {
		const startGameBtn = document.querySelector('.start-game-btn');
		startGameBtn && startGameBtn.removeEventListener('click', startGameEl);
	}

	document
		.querySelector('.joining-link-btn')
		.removeEventListener('click', copyJoiningLinkEL);

	document.querySelector('.joining-link-div').classList.add('hidden');
}

/**
 * Initialise all necessary event listeners - chat, drawing, resize canvas, joining link, start game btn, adjust overlay on scroll and modal
 */
function initGlobalEventListeners() {
	// chat
	document
		.querySelector('.send-msg')
		.addEventListener('click', sendChatMsgBtnEL);

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
		canvas.addEventListener('mousedown', startPainting);
		canvas.addEventListener('touchstart', startPainting);

		canvas.addEventListener('mouseup', stopPainting);
		canvas.addEventListener('touchend', stopPainting);

		canvas.addEventListener('mousemove', paint);
		canvas.addEventListener('touchmove', paint);

		repositionElements();
		adjustOverlay();
	});

	// resize canvas on window resize
	window.addEventListener('resize', function () {
		const { w, h } = getCanvasSize();
		canvas.width = w;
		canvas.height = h;

		const cc = document.querySelector('.canvas-container');
		cc.style.width = `${w}px`;
		cc.style.height = `${h}px`;

		adjustOverlay();
		repositionElements();
	});

	// copy joining link
	document
		.querySelector('.joining-link-btn')
		.addEventListener('click', copyJoiningLinkEL);

	// add event listener to start game button to start game
	const isOwner = JSON.parse(getFromLocalStorage('avatarConfig')).isOwner;
	const startGameBtn = document.querySelector('.start-game-btn');

	if (isOwner && startGameBtn)
		startGameBtn.addEventListener('click', startGameEl);
	else if (startGameBtn) startGameBtn.classList.add('hidden');

	// adjust overlay position on scroll
	window.addEventListener('scroll', adjustOverlay);

	// modal
	const modal = document.getElementById('modal');
	document
		.querySelector('.close-modal')
		.addEventListener('click', () => (modal.style.display = 'none'));

	window.addEventListener('click', e => {
		if (e.target === modal && modal.style.display != 'none')
			modal.style.display = 'none';
	});
}

/**
 * Alter the game bar UI, change timer icon, hide round icon and alter round number text
 * @param {Boolean} landscape true if in landscape mode
 */
function alterGameBar(landscape) {
	const timerIcon = document.querySelector('.gb-item.timer img');
	const roundIcon = document.querySelector('.gb-item.round img');
	const roundInfoText = document.querySelector('.round span.img-sub');
	const pretext = roundInfoText.querySelector('.pretext');
	const seperator = roundInfoText.querySelector('.separator');

	if (landscape) {
		roundIcon.classList.add('hidden');
		timerIcon.setAttribute('src', 'public/assets/images/clock.gif');
		pretext.textContent = 'Round';
		seperator.textContent = 'of';
	} else {
		roundIcon.classList.remove('hidden');
		timerIcon.setAttribute('src', 'public/assets/images/setting_2.gif');
		pretext.textContent = '';
		seperator.textContent = '/';
	}
}

/**
 * Reposition/hide/unhide/set css properties of elements on window resize
 */
function repositionElements() {
	const cc = document.querySelector('.canvas-container');
	const main = document.querySelector('div.main');
	const pcContainer = document.querySelector('.pc-container');
	const pcMessages = document.querySelector('.pc-container .messages');
	const pu = document.querySelector('.paint-utils');
	const cb = document.querySelector('.chat-box');
	const msgInput = document.querySelector('.msg');

	const { w, h } = getCanvasSize();

	// 1 -> move canvas-container element to-and-fro from main to pc-container depending on window size
	// 2 -> move chat-box from main to pc-container to-and-fro
	// 3 -> hide the send msg btn
	// 4 -> set messages-div max height
	// 5 -> set msg input width
	// 6 -> move paint utils to the bottom
	// 7 -> set width of paint utils
	// 8 -> alter game bar UI

	if (window.innerWidth > 768) {
		landscapeOrientation = true;

		pcContainer.insertBefore(cc, pcMessages); // 1
		pcContainer.appendChild(cb); // 2
		cb.querySelector('button').classList.add('hidden'); // 3
		pcMessages.style.maxHeight = `${h - 30}px`; // 4
		msgInput.style.width = `${pcMessages.offsetWidth - 30}px`; // 5
		main.appendChild(pu); // 6
		pu.style.width = `${w}px`; // 7
		alterGameBar(landscapeOrientation); // 8
	} else {
		landscapeOrientation = false;

		main.insertBefore(cc, pu); // 1
		main.appendChild(cb); // 2
		cb.querySelector('button').classList.remove('hidden'); // 3
		pcMessages.style.maxHeight = ``; // 4
		msgInput.style.width = ``; // 5
		main.insertBefore(pu, pcContainer); // 6
		pu.style.width = ``; // 7
		alterGameBar(landscapeOrientation); // 8
	}
}

/**
 * Event listener for copying joining link
 */
function copyJoiningLinkEL() {
	navigator.clipboard.writeText(`${window.location.origin}${joiningLink}`);
	appendChatMsgToDOM('Copied to clipboard!', '#0043ff');
}

/**
 * Event listener for starting the game
 */
function startGameEl() {
	const socketMsg = {
		type: 7,
		clientId,
		clientName,
		poolId,
	};

	sendViaSocket(socketMsg);
}

/**
 * Event listener to choose the word and send the selected word to the server
 * @param {Event} e click event on choosing the word
 * @returns void
 */
function wordChooseEL(e) {
	const chosenWord = e.target.textContent.trim();
	if (!e.currentTarget.words.includes(chosenWord)) return;

	// passing parameters to addEvenListener function by assigning the "this" element a new parameter, say words, then using it in the event handler like e.currentTarget.words

	const socketMsg = {
		type: 34,
		content: chosenWord,
		clientName,
		clientId,
		poolId,
	};

	sendViaSocket(socketMsg);
}

/**
 * Event listener to select the brush stroke width
 * @param {Event} e click event
 * @returns void
 */
function selectStrokeEL(e) {
	// if the target does not contain the class 'stroke' then return
	if (!e.target.classList.contains('stroke')) return;

	// remove the 'active' class from all the strokes and then add the class to the clicked/selected stroke
	const strokeElems = document.getElementsByClassName('stroke');
	Array.from(strokeElems).forEach(ele => ele.classList.remove('active'));
	e.target.classList.add('active');

	// get the size of the stroke and set the lineWidth in the paintUtils object
	const size = +e.target.id.slice(1);
	paintUtils.lineWidth = size * 2;

	// selected stroke width being shown on this element
	const viewImg = document.querySelector('.stroke-selected img');
	viewImg.style.width = `${size * minBrushStrokeSizeForImg}px`;
	viewImg.style.height = `${size * minBrushStrokeSizeForImg}px`;

	// add a small timeout to add the hidden class
	setTimeout(
		() =>
			document
				.querySelector('.stroke-select-options')
				.classList.add('hidden'),
		100
	);
}

/**
 * Event listener to open the stroke select menu popup
 */
function openStrokeSelectOptions() {
	// remove the hidden class
	document.querySelector('.stroke-select-options').classList.remove('hidden');

	// add EL to the strokes elements and check for click events
	const strokes = document.querySelector('.strokes');
	strokes.addEventListener('click', selectStrokeEL);
}

/**
 * Event listener to select the brush color
 * @param {Event} e click event for selecting brush color
 * @returns void
 */
function selectColorEL(e) {
	// if the target does not contain the class 'color' then return
	if (!e.target.classList.contains('color')) return;

	// remove the 'active' class from all the colors and then add the class to the clicked/selected color
	const allColors = document.getElementsByClassName('color');
	Array.from(allColors).forEach(ele => ele.classList.remove('active'));
	e.target.classList.add('active');

	// get the selected color and set the strokeStyle in the paintUtils object
	const selectedColor = e.target.style.backgroundColor;
	paintUtils.strokeStyle = selectedColor;

	// selected color being shown on this element
	const viewColor = document.querySelector('.color-selected');
	viewColor.style.backgroundColor = selectedColor;

	if (!landscapeOrientation) {
		// add a small timeout to add the hidden class, ! only in portrait mode
		setTimeout(
			() =>
				document
					.querySelector('.color-select-options')
					.classList.add('hidden'),
			100
		);
	}
}

/**
 * Event listener to open the color selection menu
 */
function openColorSelectOptions() {
	// remove the hidden class
	const colorSelectOptions = document.querySelector('.color-select-options');
	colorSelectOptions.classList.remove('hidden');

	// add EL to the colors elements and check for click events
	const colors = document.querySelector('.colors');
	colors.addEventListener('click', selectColorEL);
}

/**
 * Sets the word status on the game bar
 * @param {String} status display the status on game bar
 * @param {String} content display the content on game bar
 */
function setGbWordStatus(status, content) {
	const word = document.querySelector('.word');
	word.querySelector('span.status').textContent = status;
	word.querySelector('span.content').textContent = content;
}

/**
 * Sets the number of rounds on the game bar
 * @param {Number} num current round number
 */
function setGbRoundNum(num) {
	document.querySelector('.round span.curr-round').textContent = num;
}

/**
 * Sets the number of seconds left on the timer display element
 * @param {Number} seconds seconds left
 */
function setGbTimerStat(seconds) {
	document.querySelector('.timer span').textContent = `${seconds}s`;
}

// -------------------------------- CANVAS --------------------------------

/**
 * Initialises the canvas element, the ctx and the overlay and returns the same
 * @returns Object { canvas, ctx, overlay }
 */
function initCanvasAndOverlay() {
	const canvas = document.querySelector('.canv');
	const ctx = canvas.getContext('2d');

	// set canvas dimensions
	const { w, h } = getCanvasSize();
	canvas.width = w;
	canvas.height = h;

	// set the canvas container dimensions
	const cc = document.querySelector('.canvas-container');
	cc.style.width = `${w}px`;
	cc.style.height = `${h}px`;

	// set the overlay dimensions
	const overlay = document.querySelector('#overlay');
	overlay.style.top = `${cc.offsetTop}px`;
	overlay.style.height = `${cc.offsetHeight}px`;

	return [canvas, ctx, overlay];
}

/**
 * Calculates the canvas size for the current window size and returns the same
 * @returns Object, { w: width, h: height }
 */
function getCanvasSize() {
	const windowWidth = window.innerWidth;
	const gb = document.querySelector('.component.game-bar').offsetWidth;

	let w = 0;
	if (windowWidth >= 768) w = gb * 0.6;
	else w = gb;

	const cw = w;
	const ch = cw / 1.5;

	return { w: cw, h: ch };
}

/**
 * Returns the mouse position on the canvas
 * @param {Event} event mousemove event
 * @returns Object containing the coordinates of the mouse
 */
function getMousePos(event) {
	// https://stackoverflow.com/questions/60688935/my-canvas-drawing-app-wont-work-on-mobile

	const clientRect = canvas.getBoundingClientRect();

	if (event.type === 'touchmove' || event.type === 'touchstart') {
		return {
			x: Math.round(event.touches[0].clientX - clientRect.left),
			y: Math.round(event.touches[0].clientY - clientRect.top),
		};
	} else if (event.type === 'mousemove' || event.type === 'mousedown') {
		return {
			x: Math.round(event.clientX - clientRect.left),
			y: Math.round(event.clientY - clientRect.top),
		};
	}
}

/**
 * Starts the drawing of the sketch on the canvas
 * @param {Event} event mousedown event
 */
function startPainting(event) {
	// set painting as true
	paintUtils.isPainting = true;

	// set the prevMouse and get the current mouse position
	paintUtils.prevMouse = {
		x: paintUtils.mouse.x,
		y: paintUtils.mouse.y,
	};
	paintUtils.mouse = getMousePos(event);

	// clear the points array to encorporate new points for a new path
	paintUtils.points = [];

	// push the coordinates, lineWidth and the strokeStyle into this array
	paintUtils.points.push({
		coords: paintUtils.mouse,
		lineWidth: paintUtils.lineWidth,
		strokeStyle: paintUtils.strokeStyle,
	});
}

/**
 * Stop drawing on the canvas, mouseup event
 */
function stopPainting() {
	// set painting as false
	paintUtils.isPainting = false;

	// push the points array into the paths array
	paintUtils.paths.push(paintUtils.points);
}

/**
 * Draws the sketch on the canvas
 * @param {Event} event mousemove event
 * @returns void
 */
async function paint(event) {
	if (!paintUtils.isPainting) return;
	if (!paintUtils.hasGameStarted) return;
	if (!paintUtils.isAllowedToPaint) return;

	// set the drawing configs
	ctx.lineWidth = paintUtils.lineWidth;
	ctx.lineCap = paintUtils.lineCap;
	ctx.lineJoin = paintUtils.lineJoin;
	ctx.strokeStyle = paintUtils.strokeStyle;

	// set the prevMouse and get the current mouse position
	paintUtils.prevMouse = {
		x: paintUtils.mouse.x,
		y: paintUtils.mouse.y,
	};
	paintUtils.mouse = getMousePos(event);

	// push the coordinates, lineWidth and the strokeStyle into the points array
	paintUtils.points.push({
		coords: paintUtils.mouse,
		lineWidth: paintUtils.lineWidth,
		strokeStyle: paintUtils.strokeStyle,
	});

	// draw the path
	ctx.beginPath();
	ctx.moveTo(paintUtils.prevMouse.x, paintUtils.prevMouse.y);
	ctx.lineTo(paintUtils.mouse.x, paintUtils.mouse.y);
	ctx.stroke();

	// send this data to the server
	await wait(500);
	sendImgData();
}

/**
 * Used to redraw the paths from the paths array
 */
function drawPaths() {
	// first clear the canvas
	clearCanvas();

	// loop through each of the paths
	paintUtils.paths.forEach(path => {
		if (path.length === 0) return;

		// for each path, get the ctx config for drawing the specific styles for the strokes/paths
		ctx.lineWidth = path[0].lineWidth;
		ctx.strokeStyle = path[0].strokeStyle;

		// draw this path
		ctx.beginPath();
		ctx.moveTo(path[0].coords.x, path[0].coords.y);

		for (let i = 1; i < path.length; i++)
			ctx.lineTo(path[i].coords.x, path[i].coords.y);

		ctx.stroke();
	});
}

/**
 * Undo on canvas event listener
 */
function undo() {
	// remove the last path/stroke from the paths array
	paintUtils.paths.splice(-1, 1);
	// redraw the paths
	drawPaths();
	// send the new canvas data to the server
	sendImgDataForUndoAction();
}

/**
 * Clear canvas event listener
 */
function requestCanvasClear() {
	// clear canvas and request clear canvas on rest of the clients
	clearCanvas();

	// clear the arrays
	paintUtils.points = [];
	paintUtils.paths = [];

	const socketMsg = {
		type: 5,
		clientId,
		clientName,
		poolId,
	};

	sendViaSocket(socketMsg);
}

/**
 * Send the canvas data to the server
 */
function sendImgData() {
	// called by paint function
	const socketMsg = {
		type: 4,
		content: String(canvas.toDataURL('img/png')),
		clientName,
		clientId,
		poolId,
	};

	sendViaSocket(socketMsg);
}

/**
 * Send the canvas data to server after the undo action
 */
function sendImgDataForUndoAction() {
	// called by undo function
	const socketMsg = {
		type: 41,
		content: String(canvas.toDataURL('img/png')),
		clientName,
		clientId,
		poolId,
	};

	sendViaSocket(socketMsg);
}

// -------------------------------- ON MESSAGE HANDLERS --------------------------------

/**
 * EVENT: 1, 2, 3, 31, 312, 313
 * Appends the HTML string message into the messages box on the UI
 * @param {String} msg HTML formatted string
 * @param {String} formatColor hex code for message that will be displayed
 * @returns void
 */
function appendChatMsgToDOM(msg, formatColor) {
	if (msg.length === 0 || msg === '') return;

	const messagesDiv = document.querySelector('.messages');

	// create a new message container
	const newMsgDiv = document.createElement('div');
	newMsgDiv.classList.add('message');

	// message text
	const text = document.createElement('span');
	text.style.color = formatColor || '#1d1d1f'; // f5f5f7

	newMsgDiv.style.backgroundColor = `${formatColor}20`;
	text.innerHTML = msg;

	// append the text into message container
	newMsgDiv.appendChild(text);

	// append this new message into the messages container
	messagesDiv.appendChild(newMsgDiv);

	// bring this new message into view
	newMsgDiv.scrollIntoView(false);

	// clear the message box
	document.querySelector('.msg').value = '';
	document.querySelector('.input-wrapper span').textContent = 0;
}

/**
 * EVENT: 32
 * Reveals the word on the overlay, chat and the game bar
 * @param {Object} socketMessage
 */
function revealWordOnOverlayAndChat(socketMessage) {
	const message = `The word was '${socketMessage.content}'`;
	displayOverlay(getOverlayContentDomWithHeading(message));
	appendChatMsgToDOM(message, '#ffa500');
	setGbWordStatus('The word was', socketMessage.content);
}

/**
 * EVENT: 33
 * Show the sketcher the words they can choose
 * @param {Object} socketMessage
 */
function showWordToChoose(socketMessage) {
	// parse the words
	const words = JSON.parse(socketMessage.content);

	const overlayContent = getOverlayContentDomWithHeading(
		'Your turn, choose a word to draw!'
	);

	const wordOptionsContainer = document.createElement('div');
	wordOptionsContainer.classList.add('word-options');

	words.forEach(w => {
		const wordOptionContainer = document.createElement('div');
		wordOptionContainer.classList.add('word-option');

		const wordOption = document.createElement('span');
		wordOption.textContent = w;

		wordOptionContainer.appendChild(wordOption);
		wordOptionsContainer.appendChild(wordOptionContainer);
	});

	const timerContainer = document.createElement('div');
	const wordChooseTimerDiv = document.createElement('div');
	wordChooseTimerDiv.classList.add('word-choose-timer');

	const timeLeftSpan = document.createElement('span');
	timeLeftSpan.textContent = `${timeForChoosingWordInSeconds}s`;

	wordChooseTimerDiv.appendChild(timeLeftSpan);
	timerContainer.appendChild(wordChooseTimerDiv);

	// append to overlayContent
	overlayContent.appendChild(wordOptionsContainer);
	overlayContent.appendChild(timerContainer);

	// to access the words array outside this function scope, attach the words array to the element on which the EL is attached, like in the next line
	wordOptionsContainer.words = words;
	wordOptionsContainer.addEventListener('click', wordChooseEL);

	displayOverlay(overlayContent);

	setTimeout(() => {
		const timeoutAt = new Date(socketMessage.timeoutAfter).getTime();
		runTimer(timeLeftSpan, timeoutAt);
	}, overlayFadeInAnimationDuration + 50);
}

/**
 * EVENT: 35
 * Show to non-sketchers, that the current sketcher is choosing a word
 * @param {Object} socketMessage
 */
function showChoosingWordOnOverlay(socketMessage) {
	const overlayContent = getOverlayContentDomWithHeading(
		`${socketMessage.currSketcherName} is choosing a word!`
	);

	const ac = JSON.parse(socketMessage.content);
	const avatar = getAvatarDom(
		ac,
		scoreCardAvatarScale,
		'score-card-avatar',
		'single-avatar-on-overlay'
	);
	overlayContent.insertBefore(avatar, overlayContent.firstChild);

	displayOverlay(overlayContent);
}

/**
 * EVENT: 4
 * Display the image data on the canvas for non-sketchers
 * @param {Object} socketMessage
 */
function displayImgOnCanvas(socketMessage) {
	// display image data on canvas
	const img = new Image();
	// scale up/down canvas data based on current canvas size using outer bounds
	img.onload = () => ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
	img.setAttribute('src', socketMessage.content);
}

/**
 * EVENT: 41
 * Display the image data on canvas for non-sketchers after clearing the canvas
 * @param {Object} socketMessage
 */
function displayUndoCanvas(socketMessage) {
	clearCanvas(); // clearing the canvas for this case is important or else the new data will simply overwrite instead of clearing and redrawing
	displayImgOnCanvas(socketMessage);
}

/**
 * EVENT: 5, 51
 * Clear canvas
 */
function clearCanvas() {
	ctx.clearRect(0, 0, canvas.width, canvas.height);
}

/**
 * EVENT: 6
 * Re-renders all the players on the UI
 * @param {Object} allClients player info list containing info for rendering the players
 * @returns void
 */
function renderClients(allClients) {
	if (allClients.length === 0) return;

	const membersDiv = document.querySelector('.player-cards');
	// clear the existing dom
	membersDiv.innerHTML = '';

	// parse json into array
	allClients = JSON.parse(allClients);

	// append each player into the members div
	allClients.forEach((n, i) =>
		membersDiv.appendChild(getPlayerCardDom(n, i))
	);
}

/**
 * EVENT: 70
 * Starts the game, received ack from the server
 * @param {Object} socketMessage
 * @returns void
 */
function startGame(socketMessage) {
	if (!socketMessage.success) return;

	// flag game started
	paintUtils.hasGameStarted = true;
	log('game started');

	// remove event listeners
	removeEventListenersOnGameStart();

	// display game started overlay
	displayOverlay(getOverlayContentDomWithHeading(socketMessage.content));

	// show game started on game bar
	setGbWordStatus(socketMessage.content, socketMessage.content);

	if (socketMessage.midGameJoinee) hideOverlay();
}

/**
 * EVENT: 71
 * Show round details on overlay and the game bar
 * @param {Object} socketMessage
 */
function renderRoundDetails(socketMessage) {
	setGbRoundNum(socketMessage.currRound);

	if (!socketMessage.midGameJoinee)
		displayOverlay(
			getOverlayContentDomWithHeading(`Round ${socketMessage.currRound}`)
		);
}

/**
 * EVENT: 8
 * Begins the player's sketching flow, called for sketcher
 * @param {Object} socketMessage
 * @returns timer id for word expiry
 */
function beginClientSketchingFlow(socketMessage) {
	// hide the overlay and remove the word choosing ELs
	hideOverlay();

	// start the timer
	const wordExpiryCountdown = beginClientSketchingFlowInit(socketMessage);

	// get elements to add ELs and allow drawing
	const paintUtilsDiv = document.querySelector('.paint-utils');
	const clearCanvasBtn = document.querySelector('.pu.clear');
	const undoBtn = document.querySelector('.pu.undo');
	const colorSelected = document.querySelector('.pu.color-selected');
	const strokeSelected = document.querySelector('.pu.stroke-selected');

	// display paint utils div and add ELs
	paintUtilsDiv.classList.remove('hidden');

	// always show all colors in landscape mode
	const colorSelectOptions = document.querySelector('.color-select-options');
	if (landscapeOrientation) colorSelectOptions.classList.remove('hidden');
	else colorSelectOptions.classList.add('hidden');

	clearCanvasBtn.addEventListener('click', requestCanvasClear);
	undoBtn.addEventListener('click', undo);
	colorSelected.addEventListener('click', openColorSelectOptions);
	strokeSelected.addEventListener('click', openStrokeSelectOptions);

	// enable painting
	paintUtils.isAllowedToPaint = true;

	// to avoid scrolling on canvas on mobile devices
	if (!landscapeOrientation)
		document.querySelector('*').style.overflow = 'hidden';

	// display the word to be sketched
	setGbWordStatus('Draw this!', socketMessage.currWord);

	return wordExpiryCountdown;
}

/**
 * EVENT: 87
 * Show all non-sketchers that the sketcher is now drawing
 * @param {Object} socketMessage
 */
function showSketcherBeginDrawing(socketMessage) {
	const overlayContent = getOverlayContentDomWithHeading(
		`${socketMessage.currSketcherName} is now drawing!`
	);

	const ac = JSON.parse(socketMessage.content);
	const avatar = getAvatarDom(
		ac,
		scoreCardAvatarScale,
		'score-card-avatar',
		'single-avatar-on-overlay'
	);
	overlayContent.insertBefore(avatar, overlayContent.firstChild);

	displayOverlay(overlayContent);

	setTimeout(hideOverlay, 2000);
}

/**
 * EVENT: 88
 * Begins the flow for non-sketchers when someone else is starting to sketch
 * @param {Object} socketMessage
 * @returns timer id for word expiry
 */
function showSketcherIsDrawing(socketMessage) {
	// begin the timer
	const wordExpiryCountdown = beginClientSketchingFlowInit(socketMessage);

	let text = 'The word is hidden';
	if (socketMessage.wordMode === 'normal') {
		text = '';
		for (let i = 0; i < socketMessage.currWordLen; i++) text += '_ ';
		text = text + text.length / 2;
	}

	// show the word to be guessed - stats on the game bar
	setGbWordStatus('Guess this!', text);

	return wordExpiryCountdown;
}

/**
 * EVENT: 89
 * Display the hint string with space separator
 * @param {Object} socketMessage
 */
function displayHintString(socketMessage) {
	const hintString = socketMessage.content;

	let strToDisplay = '';
	for (let i = 0; i < hintString.length; i++)
		strToDisplay += hintString.at(i) + ' ';

	setGbWordStatus('Guess this!', strToDisplay + strToDisplay.length / 2);
}

/**
 * EVENT: 81
 * Turn over event for sketcher, disable all sketching and display time up on overlay
 */
function disableSketchingTurnOver() {
	clearAllIntervals(wordExpiryTimer);
	disableSketching();
	showTimeUp();
}

/**
 * EVENT: 82
 * Turn over event for non-sketcher, display time up on overlay
 */
function showTimeUp() {
	setGbTimerStat(0);
	clearAllIntervals(wordExpiryTimer);
	displayOverlay(getOverlayContentDomWithHeading('Time up!'));
}

/**
 * EVENT: 83
 * Disable sketching for sketcher, all guessed
 */
function disableSketchingAllGuessed() {
	clearAllIntervals(wordExpiryTimer);
	disableSketching();
	showAllHaveGuessed();
}

/**
 * EVENT: 84
 * Turn over for non-sketchers, all guessed, show everyone guessed on overlay
 */
function showAllHaveGuessed() {
	setGbTimerStat(0);
	clearAllIntervals(wordExpiryTimer);
	displayOverlay(
		getOverlayContentDomWithHeading('Everyone guessed the word!')
	);
}

/**
 * EVENT: 9
 * Render final score on overlay using table and show go home button
 * @param {Object} socketMessage
 */
function displayScores(socketMessage) {
	const dataArr = JSON.parse(socketMessage.content);

	const overlayContent = getOverlayContentDomWithHeading('Game over!');

	const table = document.createElement('table');
	const tr = document.createElement('tr');
	const avatar = document.createElement('th');
	const name = document.createElement('th');
	const score = document.createElement('th');
	name.textContent = 'Name';
	score.textContent = 'Score';

	tr.appendChild(avatar);
	tr.appendChild(name);
	tr.appendChild(score);
	table.appendChild(tr);

	dataArr.forEach(d => {
		const tr = document.createElement('tr');

		const avatar = document.createElement('th');
		const name = document.createElement('th');
		const score = document.createElement('th');

		avatar.appendChild(
			getAvatarDom(
				d.avatarConfig,
				scoreCardAvatarScale,
				'score-card-avatar'
			)
		);
		name.textContent = d.name;
		score.textContent = d.score;

		tr.appendChild(avatar);
		tr.appendChild(name);
		tr.appendChild(score);

		table.appendChild(tr);
	});

	// append table
	overlayContent.appendChild(table);

	displayOverlay(overlayContent);
	appendChatMsgToDOM('Game over!', '#ff0000');

	concludeGame();
}

/**
 * EVENT: 10
 * Init the message type map and other variables from server config
 * @param {Object} socketMessage
 */
function makeMessageTypeMapGlobal(socketMessage) {
	const content = JSON.parse(socketMessage.content);

	timeForEachWordInSeconds = content.timeForEachWordInSeconds;
	timeForChoosingWordInSeconds = content.timeForChoosingWordInSeconds;
	allowLogs = content.printLogs;
	closeSocketConnInSeconds = content.closeSocketConnInSeconds * 1000;

	const m = content.messageTypeMap;
	const keys = Object.keys(m);
	messageTypeMap = new Map();

	keys.forEach(k => messageTypeMap.set(Number(k), m[k]));
}

function concludeGame() {
	// render go home button
	setTimeout(() => {
		const btnDiv = document.querySelector('.joining-link-div');
		btnDiv.classList.remove('hidden');

		const btn = document.querySelector('.joining-link-btn');
		btn.textContent = 'Play again!';
		btn.addEventListener('click', () => {
			closeSocketConn('game over, clicked on play again!');
			window.location.href = homeRoute;
		});
	}, 2000);

	setTimeout(() => {
		closeSocketConn('game over, closing connection after timeout');
	}, closeSocketConnInSeconds);
}

// -------------------------------- SOCKET --------------------------------

/**
 * Initialises a socket connection to the server adds corresponding event handlers to the socket
 * @returns {socket} socket connection
 */
function initSocket() {
	// get the avatar config from the local storage
	const avatarConfig = getFromLocalStorage('avatarConfig');

	let wsProtocol = 'ws';
	if (window.location.protocol === 'https:') wsProtocol = 'wss';

	// construct the web socket url with the required params
	const wsUrl = `${wsProtocol}://${window.location.host}/scribble/ws?poolId=${poolId}&clientId=${clientId}&clientName=${clientName}&avatarConfig=${avatarConfig}`;

	// make the connection
	const socket = new WebSocket(wsUrl);

	// attach event handlers to the socket
	socket.onopen = () => log('Socket successfully connected!');
	socket.onerror = error => log('Socket error', error);
	socket.onmessage = socketOnMessage;
	socket.onclose = socketOnClose;

	return socket;
}

/**
 * Socket onmessage handler
 * @param {any} message raw message received from server
 */
function socketOnMessage(message) {
	// runs when a message is received on the socket conn, runs the corresponding functions depending on message type

	// parse json string into json object
	const socketMessage = JSON.parse(message.data);

	if (socketMessage.type !== 4)
		log(
			socketMessage.type,
			messageTypeMap && messageTypeMap.get(socketMessage.type)
		);

	switch (socketMessage.type) {
		case 1:
			if (socketMessage.clientId === clientId)
				// if the current clientId and the clientId from response match then
				appendChatMsgToDOM(
					`You joined the room as <strong>${socketMessage.clientName}</strong>!`,
					''
				);
			else
				appendChatMsgToDOM(
					`<strong>${socketMessage.clientName}</strong> has joined the room!`,
					''
				);
			break;

		case 2:
			appendChatMsgToDOM(
				`<strong>${socketMessage.clientName}</strong> has left the room!`,
				''
			);
			break;

		case 3:
			appendChatMsgToDOM(
				`<strong>${socketMessage.clientName}</strong>: ${socketMessage.content}`,
				''
			);
			break;

		case 31:
			appendChatMsgToDOM(
				`<strong>${socketMessage.clientName}</strong> guessed the word!`,
				'#00ff00'
			);
			break;

		case 312:
			appendChatMsgToDOM(
				`Naughty <strong>@${socketMessage.clientName}</strong>`,
				'#ff0000'
			);
			break;

		case 313:
			appendChatMsgToDOM(
				`You can't reveal the word <strong>@${socketMessage.clientName}</strong>`,
				'#ff7f00'
			);
			break;

		case 32:
			revealWordOnOverlayAndChat(socketMessage);
			break;

		case 33:
			showWordToChoose(socketMessage);
			break;

		case 35:
			showChoosingWordOnOverlay(socketMessage);
			break;

		case 4:
			displayImgOnCanvas(socketMessage);
			break;

		case 41:
			displayUndoCanvas(socketMessage);
			break;

		case 5:
		case 51:
			clearCanvas();
			break;

		case 6:
			renderClients(socketMessage.content);
			break;

		case 69:
			appendChatMsgToDOM(
				'You need at least two players to start the game',
				'#457ef4'
			);
			break;

		case 70:
			startGame(socketMessage);
			break;

		case 71:
			renderRoundDetails(socketMessage);
			break;

		case 8:
			wordExpiryTimer = beginClientSketchingFlow(socketMessage);
			break;

		case 86:
			wordExpiryTimer = beginClientSketchingFlowInit(socketMessage);
			break;

		case 87:
			showSketcherBeginDrawing(socketMessage);
			break;

		case 88:
			wordExpiryTimer = showSketcherIsDrawing(socketMessage);
			break;

		case 89:
			displayHintString(socketMessage);
			break;

		case 81:
			disableSketchingTurnOver();
			break;

		case 82:
			showTimeUp();
			break;

		case 83:
			disableSketchingAllGuessed();
			break;

		case 84:
			showAllHaveGuessed();
			break;

		case 9:
			displayScores(socketMessage);
			break;

		case 10:
			makeMessageTypeMapGlobal(socketMessage);
			break;

		default:
			break;
	}
}

/**
 * Socket onclose handler
 */
function socketOnClose() {
	// on socket conn close, stop all timer or intervals
	log('Socket connection closed, stopping timers and timeouts!');
	clearAllIntervals(wordExpiryTimer);

	// display connection lost on the modal
	document.getElementById('modal').style.display = 'flex';
}

function closeSocketConn(reason) {
	socket.close(1000, reason);
}

/**
 * Send the socketMsg to the server if connected, else show disconnected prompt
 * @param {Object} socketMsg message to be sent to the server
 */
function sendViaSocket(socketMsg) {
	/*  socket.readyState: int
			0 - connecting
			1 - open
			2 - closing
			3 - closed
	*/

	// if socket is in open state then send the message
	if (socket.readyState === socket.OPEN)
		socket.send(JSON.stringify(socketMsg));
	else {
		// clear any intervals and show connection lost
		log('socket current state:', socket.readyState);
		clearAllIntervals(wordExpiryTimer);

		// display connection lost on the modal
		document.getElementById('modal').style.display = 'flex';
	}
}

// -------------------------------- MAIN --------------------------------

const overlayFadeInAnimationDuration = 300;
const minBrushStrokeSizeForImg = 6;

const customizableAvatarScale = 3;
const rowAvatarScale = 1;
const playerCardAvatarScale = 0.5;
const scoreCardAvatarScale = 0.75;

// utils for painting on canvas
const paintUtils = {
	isPainting: false,
	hasGameStarted: false,
	isAllowedToPaint: false,
	points: [],
	paths: [],
	mouse: { x: 0, y: 0 },
	prevMouse: { x: 0, y: 0 },
	lineWidth: 2,
	lineCap: 'round',
	lineJoin: 'round',
	strokeStyle: '#000',
};

let messageTypeMap,
	timeForEachWordInSeconds,
	timeForChoosingWordInSeconds,
	wordExpiryTimer,
	allowLogs,
	closeSocketConnInSeconds,
	landscapeOrientation = false;

const boundariesForAvatarAtlas = {
	color: { x: 5, y: 2 },
	eyes: { x: 6, y: 5 },
	mouth: { x: 0, y: 5 },
};

const validCoordsForAvatarAtlas = {
	color: initValidCoordsForAvatarAtlas(boundariesForAvatarAtlas.color),
	eyes: initValidCoordsForAvatarAtlas(boundariesForAvatarAtlas.eyes),
	mouth: initValidCoordsForAvatarAtlas(boundariesForAvatarAtlas.mouth),
};

const avatarConfig = {
	color: { x: 0, y: 0 },
	eyes: { x: 0, y: 0 },
	mouth: { x: 0, y: 0 },
	isOwner: false,
	isCrowned: false,
};

let canvas, ctx, overlay, socket;

switch (renderTemplateName) {
	case 'game':
		[canvas, ctx, overlay] = initCanvasAndOverlay();
		socket = initSocket();
		initGlobalEventListeners();
		break;

	case 'home':
	case 'join':
		renderCustomisableAvatar();
		renderRowAvatars();
		break;

	default:
		renderRowAvatars();
		break;
}
