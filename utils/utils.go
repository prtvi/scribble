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

func GenerateUUID() string {
	return uuid.New().String()[0:8]
}

func FormatTimeLong(t time.Time) string {
	return t.Format(time.RFC3339Nano) // 2021-12-12T12:23:34.002342369
}

func GetTimeString(t time.Time) string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
}

func GetSecondsLeftFrom(t time.Time) int {
	return int(t.Sub(time.Now()).Seconds())
}

func CalcScore(scoreForCorrectGuess, currRound int, currWordExpiresAt time.Time) int {
	return scoreForCorrectGuess*currRound*GetDiffBetweenTimesInSeconds(time.Now(), currWordExpiresAt) + scoreForCorrectGuess
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

func GetDiffBetweenTimesInSeconds(t1, t2 time.Time) int {
	return int(math.Abs(t1.Sub(t2).Seconds()))
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
