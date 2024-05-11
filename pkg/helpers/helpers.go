package helpers

import (
	"strconv"
	"strings"
	"time"
)

func ParsePositiveInt(line string) (int, bool) {
	number, err := strconv.Atoi(line)
	if err != nil || number < 0 {
		return 0, false
	}

	return number, true
}

func ParseTime(str string) (time.Time, error) {
	return time.Parse("15:04", str)
}

func IsAllowedChars(line string, allowedChars string) bool {
	for _, char := range line {
		if !strings.ContainsRune(allowedChars, char) {
			return false
		}
	}
	return true
}
