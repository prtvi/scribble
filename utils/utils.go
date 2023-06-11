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

var WORDS = []string{"mountain", "water", "ball", "bottle", "car", "tree", "clock", "shirt", "lollipop", "soap", "bag", "umbrella", "pillow", "Shoes", "headphones", "bulb", "fan", "fruit", "grass", "coffee", "phone"}

func GenerateUUID() string {
	return uuid.New().String()
}

func FormatTimeLong(t time.Time) string {
	// 2021-12-12T12:23:34.002342369
	return t.Format(time.RFC3339Nano)
}

func GetSecondsLeftFrom(t time.Time) int {
	return int(t.Sub(time.Now()).Seconds())
}

func GetRandomWord(arr []string) string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int() % len(arr)

	return arr[n]
}

func Get3RandomWords(arr []string) []string {
	ret := make([]string, 3)

	for i := 0; i < 3; i++ {
		ret[i] = GetRandomWord(arr)
	}

	return ret
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
		fmt.Println("Error loading .env file")
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
	c, ok := colorMap[color]
	if ok {
		return c
	} else {
		return reset
	}
}

func Cp(color string, message ...string) {
	fmt.Printf("%s%s%s\n", getColor(color), strings.Join(message, " "), reset)
}

func Cs(color string, message ...string) string {
	return fmt.Sprintf("%s%s%s", getColor(color), strings.Join(message, " "), reset)
}
