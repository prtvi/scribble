package utils

import "scribble/model"

var reset string = "\033[0m"
var colorMap map[string]string = map[string]string{
	"reset": reset,

	// colors
	"black":  "\033[0;30m",
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"purple": "\033[35m",
	"cyan":   "\033[36m",
	"white":  "\033[37m",

	// underline
	"blackU":  "\033[4;30m",
	"redU":    "\033[4;31m",
	"greenU":  "\033[4;32m",
	"yellowU": "\033[4;33m",
	"blueU":   "\033[4;34m",
	"purpleU": "\033[4;35m",
	"cyanU":   "\033[4;36m",
	"whiteU":  "\033[4;37m",

	// backgrounds
	"blackBg":  "\033[40m",
	"redBg":    "\033[41m",
	"greenBg":  "\033[42m",
	"yellowBg": "\033[43m",
	"blueBg":   "\033[44m",
	"purpleBg": "\033[45m",
	"cyanBg":   "\033[46m",
	"whiteBg":  "\033[47m",
}

var COLORS_FOR_DRAWING = []string{"#fff", "#c1c1c1", "#ef100c", "#ff7002", "#ffe300", "#00cc00", "#00ff92", "#00b2fe", "#231fd3", "#a302ba", "#df69a7", "#ffad8e", "#a1512d", "#000", "#505050", "#750b07", "#c23a00", "#e8a200", "#004618", "#06785d", "#00579e", "#0f0665", "#550069", "#886353", "#cc784c", "#64300d"}

var WORDS = []string{"hammer", "goggles", "tiger", "candle", "hair", "glass", "perfume", "hanger", "bowl", "ball", "poop", "radio", "panda", "shovel", "cap", "glasses", "mirror", "newspaper", "nail", "belt", "truck", "table", "marker", "sign", "button", "frog", "games", "card", "fork", "cucumber", "suitcase", "mug", "bag", "ring", "tissue", "dictionary", "sponge", "seat", "cord", "lamp", "box", "crayons", "hat", "flashlight", "bottle", "milk", "note", "carrots", "computer", "tire", "earrings", "pan", "locket", "email", "watch", "brush", "toothbrush", "flag", "freezer", "dice", "pool", "puddle", "pencil", "scissors", "spatula", "hook", "toy", "paperclip", "necktie", "door", "chalk", "mousepad", "mountain", "water", "car", "tree", "clock", "shirt", "lollipop", "soap", "umbrella", "pillow", "shoes", "headphones", "bulb", "fan", "fruit", "grass", "coffee", "phone", "notepad", "string", "comb", "eraser", "stamp", "bell", "ice", "spoon", "chicken", "egg", "painting", "notebook", "purse", "sword", "book", "desk", "plant", "tie", "candy", "lighter", "oil", "flowers", "racket", "wrench", "bracelet", "duck", "light", "globe", "house", "necklace", "keyboard", "screw", "jar", "wire", "camera", "paper", "marble", "mask", "deodorant", "pants", "socks", "bookmark", "toothpaste", "stickers", "shawl", "dove", "cork", "needle", "conditioner", "rat", "monitor", "stick", "key", "bananas", "television", "toilet", "chain", "microphone", "lightbulb", "lipstick", "shoelace", "zebra", "whistle", "controller", "bed", "map", "helmet", "paintbrush", "magnet", "boat", "wallet", "sofa", "mop", "lemon", "dart", "seed", "zip", "feather", "turtle"}

var AboutText = []string{"scribble is a free online multiplayer drawing and guessing pictionary game.", "A normal game consists of a few rounds, where every round a player has to draw their chosen word and others have to guess it to gain points!", "The person with the most points at the end of the game, will then be crowned as the winner!"}

var HowToSlides = []string{"When it's your turn, choose a word you want to draw!",
	"Try to draw your choosen word! No spelling!",
	"Let other players try to guess your drawn word!",
	"When it's not your turn, try to guess what other players are drawing!",
	"Score the most points and be crowned the winner at the end!"}

var FormParams = []model.CreateFormParam{
	{ID: "players", Label: "Players", ImgIdx: 1, Desc: "Number of players in the room",
		Options: []model.FormOption{
			{Value: "2", Label: "2"},
			{Value: "3", Label: "3"},
			{Value: "4", Label: "4"},
			{Value: "5", Label: "5", Selected: true},
			{Value: "6", Label: "6"},
			{Value: "7", Label: "7"},
			{Value: "8", Label: "8"},
			{Value: "9", Label: "9"},
			{Value: "10", Label: "10"}}},

	{ID: "drawTime", Label: "Draw time", ImgIdx: 2, Desc: "Number of seconds each player gets to sketch",
		Options: []model.FormOption{
			{Value: "15", Label: "15"},
			{Value: "20", Label: "20"},
			{Value: "40", Label: "40"},
			{Value: "50", Label: "50"},
			{Value: "60", Label: "60"},
			{Value: "70", Label: "70"},
			{Value: "80", Label: "80", Selected: true},
			{Value: "90", Label: "90"},
			{Value: "100", Label: "100"},
			{Value: "120", Label: "120"},
			{Value: "150", Label: "150"},
			{Value: "180", Label: "180"},
			{Value: "210", Label: "210"},
			{Value: "240", Label: "240"}}},

	{ID: "rounds", Label: "Rounds", ImgIdx: 3, Desc: "Number of rounds",
		Options: []model.FormOption{
			{Value: "2", Label: "2"},
			{Value: "3", Label: "3", Selected: true},
			{Value: "4", Label: "4"},
			{Value: "5", Label: "5"},
			{Value: "6", Label: "6"},
			{Value: "7", Label: "7"},
			{Value: "8", Label: "8"},
			{Value: "9", Label: "9"},
			{Value: "10", Label: "10"}}},

	{ID: "wordMode", Label: "Word mode", ImgIdx: 4, Desc: "Word mode, Normal: display number of characters in word, also display hints. Hidden: do not reveal the number of characters, no hints. Combination: combination of two words separated by '+'",
		Options: []model.FormOption{
			{Value: "normal", Label: "Normal", Selected: true},
			{Value: "hidden", Label: "Hidden"},
			{Value: "combination", Label: "Combination"}}},

	{ID: "wordCount", Label: "Word count", ImgIdx: 5, Desc: "Number of words the sketcher gets to choose from to sketch",
		Options: []model.FormOption{
			{Value: "1", Label: "1"},
			{Value: "2", Label: "2"},
			{Value: "3", Label: "3", Selected: true},
			{Value: "4", Label: "4"},
			{Value: "5", Label: "5"}}},

	{ID: "hints", Label: "Hints", ImgIdx: 6, Desc: "Number of characters in the word to be revealed as hints",
		Options: []model.FormOption{
			{Value: "1", Label: "1"},
			{Value: "2", Label: "2", Selected: true},
			{Value: "3", Label: "3"},
			{Value: "4", Label: "4"},
			{Value: "5", Label: "5"}}},
}
