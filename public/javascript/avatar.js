'use strict';

const color = document.querySelector('.avatar .color');
const eyes = document.querySelector('.avatar .eyes');
const mouth = document.querySelector('.avatar .mouth');

const scaleAvatarBy = 3;
const offset = 48 * scaleAvatarBy;
const rows = 10;

const avatarConfig = {
	color: { x: 0, y: 0 },
	eyes: { x: 0, y: 0 },
	mouth: { x: 0, y: 0 },
	isOwner: false,
	isCrowned: false,
};

const boundaries = {
	color: { x: 5, y: 2 },
	eyes: { x: 6, y: 5 },
	mouth: { x: 0, y: 5 },
};

const validCoords = {
	color: initValidCoords(boundaries.color),
	eyes: initValidCoords(boundaries.eyes),
	mouth: initValidCoords(boundaries.mouth),
};

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

document.querySelector('.randomize').addEventListener('click', randomizeAvatar);

[colorLeft, eyesLeft, mouthLeft].forEach(ele =>
	ele.addEventListener('click', leftEL)
);

[colorRight, eyesRight, mouthRight].forEach(ele =>
	ele.addEventListener('click', rightEL)
);

function setBgPosition(element, x, y) {
	element.style.backgroundPositionX = `-${x * offset}px`;
	element.style.backgroundPositionY = `-${y * offset}px`;
}

function getCurrPosition(pos) {
	const lastIdx = pos.lastIndexOf('px');
	return Math.abs(+pos.slice(0, lastIdx)) / offset;
}

function rightEL(e) {
	const name = e.currentTarget.name;
	const elem = document.querySelector(`.${name}`);

	if (elem.style.backgroundPositionX === '')
		elem.style.backgroundPositionX = '0px';

	avatarConfig[name].x = getCurrPosition(elem.style.backgroundPositionX) + 1;

	if (avatarConfig[name].x >= rows) {
		avatarConfig[name].y += 1;
		avatarConfig[name].x = 0;
	}

	if (
		avatarConfig[name].x > boundaries[name].x &&
		avatarConfig[name].y === boundaries[name].y
	) {
		avatarConfig[name].x = 0;
		avatarConfig[name].y = 0;
	}

	setBgPosition(elem, avatarConfig[name].x, avatarConfig[name].y);
	setIfOwner();
	saveToLocalStorage('avatarConfig', avatarConfig);
}

function leftEL(e) {
	const name = e.currentTarget.name;
	const elem = document.querySelector(`.${name}`);

	if (elem.style.backgroundPositionX === '')
		elem.style.backgroundPositionX = '0px';

	avatarConfig[name].x = getCurrPosition(elem.style.backgroundPositionX) - 1;

	if (avatarConfig[name].x < 0 && avatarConfig[name].y > 0) {
		avatarConfig[name].y -= 1;
		avatarConfig[name].x = rows - 1;
	}

	if (avatarConfig[name].x < 0 && avatarConfig[name].y === 0) {
		avatarConfig[name].x = boundaries[name].x;
		avatarConfig[name].y = boundaries[name].y;
	}

	setBgPosition(elem, avatarConfig[name].x, avatarConfig[name].y);
	setIfOwner();
	saveToLocalStorage('avatarConfig', avatarConfig);
}

function randomizeAvatar() {
	const coords = getRandomizedAvatarCoords();

	avatarConfig.color = coords.color;
	avatarConfig.eyes = coords.eyes;
	avatarConfig.mouth = coords.mouth;

	setBgPosition(color, coords.color.x, coords.color.y);
	setBgPosition(eyes, coords.eyes.x, coords.eyes.y);
	setBgPosition(mouth, coords.mouth.x, coords.mouth.y);

	setIfOwner();
	saveToLocalStorage('avatarConfig', avatarConfig);
}

function getRandomValue(arr) {
	return arr[Math.floor(Math.random() * arr.length)];
}

function getRandomizedAvatarCoords() {
	return {
		color: getRandomValue(validCoords.color),
		eyes: getRandomValue(validCoords.eyes),
		mouth: getRandomValue(validCoords.mouth),
	};
}

function initValidCoords(prop) {
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

function setIfOwner() {
	const urlParams = new URLSearchParams(location.search);
	if (urlParams.get('isOwner') === 'true') avatarConfig.isOwner = true;
}

function saveToLocalStorage(key, value) {
	window.localStorage.setItem(key, JSON.stringify(value));
}

function getFromLocalStorage(key) {
	return window.localStorage.getItem(key);
}
