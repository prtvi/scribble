package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

var COLORS = []string{"36fdc3", "ece8f8", "90c335", "d17161", "a16014", "2f38a0", "11ea10", "9e5df3", "87425b", "180dab", "91ff00", "00ffc8", "00fff2", "ff5100", "ffe100", "ddff00", "c8ff00", "ff2200", "8FB803", "FAC8D2", "9A2EAF", "C099A6", "FB1974", "CC4E0B", "8FB288", "4A2073", "9484E0", "0A2980", "399299", "E0F066", "159988", "309C0D", "2CE997", "BB0B39", "D5860F", "38204E", "583E3C", "F4F4F4", "1C542D", "8673A1"}

var WORDS = []string{"Mountain", "Water", "Ball", "Bottle", "Car", "Tree", "Clock", "T-shirt", "Lollipop", "Soap", "Bag", "Umbrella", "Pillow", "Shoes", "Headphones", "Bulb", "Fan", "Fruit", "Grass", "Coffee", "Phone"}

func GenerateUUID() string {
	return uuid.New().String()
}

func FormatTimeLong(t time.Time) string {
	// 2023-09-19 23:34:09
	return t.String()[0:19]
}

func GetRandomWord() string {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(WORDS)

	return WORDS[n]
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
