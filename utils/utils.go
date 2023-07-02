package utils

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var COLORS = []string{"36fdc3", "ff2200", "90c335", "d17161", "a16014", "2f38a0", "11ea10", "9e5df3", "87425b", "180dab", "91ff00", "00ffc8", "00fff2", "ff5100", "ffe100", "ddff00", "c8ff00", "ece8f8", "8fb803", "fac8d2", "9a2eaf", "c099a6", "fb1974", "cc4e0b", "8fb288", "4a2073", "9484e0", "0a2980", "399299", "e0f066", "159988", "309c0d", "2ce997", "bb0b39", "d5860f", "38204e", "583e3c", "f4f4f4", "1c542d", "8673a1"}

var WORDS = []string{"hammer", "goggles", "tiger", "candle", "hair", "glass", "perfume", "hanger", "bowl", "ball", "poop", "radio", "panda", "shovel", "cap", "glasses", "mirror", "newspaper", "nail", "belt", "truck", "table", "marker", "sign", "button", "frog", "games", "card", "fork", "cucumber", "suitcase", "mug", "bag", "ring", "tissue", "dictionary", "sponge", "seat", "cord", "lamp", "box", "crayons", "hat", "flashlight", "bottle", "milk", "note", "carrots", "computer", "tire", "earrings", "pan", "locket", "email", "watch", "brush", "toothbrush", "flag", "freezer", "dice", "pool", "puddle", "pencil", "scissors", "spatula", "hook", "toy", "paperclip", "necktie", "door", "chalk", "mousepad", "mountain", "water", "car", "tree", "clock", "shirt", "lollipop", "soap", "umbrella", "pillow", "shoes", "headphones", "bulb", "fan", "fruit", "grass", "coffee", "phone", "notepad", "string", "comb", "eraser", "stamp", "bell", "ice", "spoon", "chicken", "egg", "painting", "notebook", "purse", "sword", "book", "desk", "plant", "tie", "candy", "lighter", "oil", "flowers", "racket", "wrench", "bracelet", "duck", "light", "globe", "house", "necklace", "keyboard", "screw", "jar", "wire", "camera", "paper", "marble", "mask", "deodorant", "pants", "socks", "bookmark", "toothpaste", "stickers", "shawl", "dove", "cork", "needle", "conditioner", "rat", "monitor", "stick", "key", "bananas", "television", "toilet", "chain", "microphone", "lightbulb", "lipstick", "shoelace", "zebra", "whistle", "controller", "bed", "map", "helmet", "paintbrush", "magnet", "boat", "wallet", "sofa", "mop", "lemon", "dart", "seed", "zip", "feather", "turtle"}

func GenerateUUID() string {
	return uuid.New().String()
}

func FormatTimeLong(t time.Time) string {
	// 2021-12-12T12:23:34.002342369
	return t.Format(time.RFC3339Nano)
}

func GetTimeString(t time.Time) string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
}

func GetSecondsLeftFrom(t time.Time) int {
	return int(t.Sub(time.Now()).Seconds())
}

func GetRandomWord(arr []string) string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int() % len(arr)

	return arr[n]
}

func GetNrandomWords(arr []string, n int) []string {
	ret := make([]string, n)

	for i := 0; i < n; i++ {
		ret[i] = GetRandomWord(arr)
	}

	return ret
}

func SplitIntoWords(s string) []string {
	arr := strings.Split(s, ",")
	trimmed := make([]string, 0)

	for _, val := range arr {
		trimmedValue := strings.Trim(val, " ")

		if len(trimmedValue) > 0 {
			trimmed = append(trimmed, trimmedValue)
		}
	}

	return trimmed
}

func GetDiffBetweenTimesInSeconds(t1, t2 time.Time) float64 {
	return math.Abs(t1.Sub(t2).Seconds())
}

func ShuffleList(list []string) []string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })

	return list
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		Cp("redBg", "Error loading .env file")
	}
}

func DurationToSeconds(t time.Duration) int {
	return int(t.Seconds())
}

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

func getColor(color string) string {
	if c, ok := colorMap[color]; ok {
		return c
	}
	return reset
}

func Cp(color string, message ...string) {
	fmt.Printf("%s%s%s\n", getColor(color), strings.Join(message, " "), reset)
}

func Cs(color string, message ...string) string {
	return fmt.Sprintf("%s%s%s", getColor(color), strings.Join(message, " "), reset)
}
