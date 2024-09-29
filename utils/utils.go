package utils

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func IsProdEnv() bool {
	return GetEnvVar("ENV") == "PROD"
}

func GetEnvVar(key string) string {
	return os.Getenv(key)
}

func GetPort() (port string) {
	if !IsProdEnv() {
		DEBUG = true
		port = "1323"
		Cp("greenBg", "----------- DEV/DEBUG ENV -----------")
		return
	}

	DEBUG = false
	port = GetEnvVar("PORT")
	if port == "" {
		port = "1323"
	}

	Cp("redBg", "----------- PROD ENV -----------")
	return
}

func logToFile(content string) {
	f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		fmt.Println(err)
	}
}

func GenerateUUID() string {
	return uuid.New().String()[0:8]
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

// time utils

func Sleep(d time.Duration) {
	time.Sleep(d)
}

func SleepWithInterrupt(d time.Duration, stop chan bool) bool {
	// this func can be used to sleep for d duration, with an interuppt if any to stop this sleep
	// to achieve this interrupt before timeout, pass a channel bool, which will be used to break this timeout
	// this chan needs to be used to pass acknowledgement for stopping this timeout
	// returns boolean whether the timeout was interrupted or not, if interrupted then returns true

	select {
	case <-stop:
		return true
	case <-time.After(d):
		return false
	}
}

func FormatTimeLong(t time.Time) string {
	return t.Format(time.RFC3339Nano) // 2021-12-12T12:23:34.002342369
}

func GetTimeString(t time.Time) string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
}

func GetSecondsLeftFrom(t time.Time) int {
	return int(time.Until(t).Seconds())
}

func DurationToSeconds(t time.Duration) int {
	return int(t.Seconds())
}

func GetDiffBetweenTimesInSeconds(t1, t2 time.Time) int {
	return int(math.Abs(t1.Sub(t2).Seconds()))
}

// game logic/calculations

func CalcScore(scoreForCorrectGuess, currRound int, currWordExpiresAt time.Time) int {
	return scoreForCorrectGuess*currRound*GetDiffBetweenTimesInSeconds(time.Now(), currWordExpiresAt) + scoreForCorrectGuess
}

func CalculateMaxHintsAllowedForWord(currWord string, nHintsPref int) int {
	maxHintsAllowed := len(currWord) / 2
	if nHintsPref <= maxHintsAllowed {
		maxHintsAllowed = nHintsPref
	}

	return maxHintsAllowed
}

func GetHintString(word, char, hintString string) string {
	for i, c := range word {
		charString := string(c)
		if charString == char && string(hintString[i]) == "_" {
			hintString = hintString[:i] + charString + hintString[i+1:]
			break
		}
	}

	return hintString
}

func SelectWordsForPool(globalWords, inputWords []string, useCustomWordsOnly bool) []string {
	if len(inputWords) == 0 {
		return globalWords
	}

	if len(inputWords) >= 10 && useCustomWordsOnly {
		return inputWords
	}

	// inputWords share, 2/3
	// totalWords,  t
	// inputWords,  i = (2/3) * t
	//              t = (3/2) * i
	// globalWords, g = (1/3) * t
	// solve for g wrt to i, that's the number of words that will be fetched from the global array
	// here, g = i/2

	inputShare := 0.66
	globalWordsShare := int(math.Round((1 - inputShare) * (float64(len(inputWords)) / float64(inputShare))))

	result := make([]string, 0)
	result = append(result, inputWords...)
	result = append(result, GetNrandomWords(globalWords, globalWordsShare)...)

	return ShuffleList(result)
}

// randomise

func GetRandomItem(arr []string) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rand.Int() % len(arr)
	return arr[n]
}

func GetRandomItemWithIdx(arr []string) (string, int) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rand.Int() % len(arr)
	return arr[n], n
}

func GetNrandomWords(arr []string, n int) []string {
	ret := make([]string, n)
	for i := 0; i < n; i++ {
		ret[i] = GetRandomItem(arr)
	}
	return ret
}

func PickRandomCharacter(chars [](string)) ([]string, string) {
	charPicked, idx := GetRandomItemWithIdx(chars)
	chars = append(chars[:idx], chars[idx+1:]...)
	return chars, charPicked
}

func ShuffleList(list []string) []string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })
	return list
}

// color printing

func getColor(color string) string {
	if c, ok := colorMap[color]; ok {
		return c
	}
	return reset
}

// solids: "black", "red", "green", "yellow", "blue", "purple", "cyan", "white"
// underline: "blackU", "redU", "greenU", "yellowU", "blueU", "purpleU", "cyanU", "whiteU"
// backgrounds: "blackBg", "redBg", "greenBg", "yellowBg", "blueBg", "purpleBg", "cyanBg", "whiteBg"
func Cp(color string, message ...any) {
	msg := ""
	for _, m := range message {
		msg += fmt.Sprintf("%+v ", m)
	}
	if len(msg) > 1 {
		msg = msg[:len(msg)-1]
	}

	if IsProdEnv() {
		fmt.Println(msg)
	} else {
		fmt.Printf("%s%s%s\n", getColor(color), msg, reset)
		logToFile(fmt.Sprintf("%s: %s\n", FormatTimeLong(time.Now())[:19], msg))
	}
}

func Cs(color string, message ...string) string {
	return fmt.Sprintf("%s%s%s", getColor(color), strings.Join(message, " "), reset)
}
