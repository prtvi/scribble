package utils

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func LoadAndGetEnv() bool {
	err := godotenv.Load(".env")
	if err != nil {
		Cp("redBg", "Error loading .env file")
	}

	env := os.Getenv("ENV")
	if env == "" || env == "PROD" {
		return false
	}

	return true
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

func Cp(color string, message ...string) {
	fmt.Printf("%s%s%s\n", getColor(color), strings.Join(message, " "), reset)
}

func Cs(color string, message ...string) string {
	return fmt.Sprintf("%s%s%s", getColor(color), strings.Join(message, " "), reset)
}
