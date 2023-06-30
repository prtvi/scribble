'use strict';

const color = document.querySelector('.avatar .color');
const eyes = document.querySelector('.avatar .eyes');
const mouth = document.querySelector('.avatar .mouth');

const scaleBy = 3;
const offset = 48 * scaleBy;
const rows = 10;

const positions = {
	color: { x: 0, y: 0 },
	eyes: { x: 0, y: 0 },
	mouth: { x: 0, y: 0 },
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

	positions[name].x = getCurrPosition(elem.style.backgroundPositionX) + 1;

	if (positions[name].x >= rows) {
		positions[name].y += 1;
		positions[name].x = 0;
	}

	if (
		positions[name].x > boundaries[name].x &&
		positions[name].y === boundaries[name].y
	) {
		positions[name].x = 0;
		positions[name].y = 0;
	}

	setBgPosition(elem, positions[name].x, positions[name].y);
}

function leftEL(e) {
	const name = e.currentTarget.name;
	const elem = document.querySelector(`.${name}`);

	if (elem.style.backgroundPositionX === '')
		elem.style.backgroundPositionX = '0px';

	positions[name].x = getCurrPosition(elem.style.backgroundPositionX) - 1;

	if (positions[name].x < 0 && positions[name].y > 0) {
		positions[name].y -= 1;
		positions[name].x = rows - 1;
	}

	if (positions[name].x < 0 && positions[name].y === 0) {
		positions[name].x = boundaries[name].x;
		positions[name].y = boundaries[name].y;
	}

	setBgPosition(elem, positions[name].x, positions[name].y);
}

function randomizeAvatar() {
	const coords = getRandomizedAvatarCoords();

	positions.color = coords.color;
	positions.eyes = coords.eyes;
	positions.mouth = coords.mouth;

	setBgPosition(color, coords.color.x, coords.color.y);
	setBgPosition(eyes, coords.eyes.x, coords.eyes.y);
	setBgPosition(mouth, coords.mouth.x, coords.mouth.y);
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

// how to play slideshow

let slideIndex = 0;

const slides = document.getElementsByClassName('slides');
const dots = document.getElementsByClassName('dot');

showSlides();

for (let i = 0; i < dots.length; i++) {
	dots[i].addEventListener('click', e => {
		const id = e.currentTarget.id;
		slideIndex = +id.slice(id.length - 1);

		showCurrSlide(slideIndex);
	});
}

function showCurrSlide(idx) {
	for (let i = 0; i < slides.length; i++) {
		slides[i].style.display = 'none';
		dots[i].classList.remove('active');
	}

	slides[idx].style.display = 'block';
	dots[idx].classList.add('active');
}

function showSlides() {
	showCurrSlide(slideIndex);

	slideIndex += 1;
	if (slideIndex >= slides.length) slideIndex = 0;

	setTimeout(showSlides, 2500);
}
