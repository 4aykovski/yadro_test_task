package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/4aykovski/yadro_test_task/internal/app"
	"github.com/4aykovski/yadro_test_task/pkg/helpers"
)

func main() {
	if len(os.Args) < 2 {
		slog.Error("arguments is not specified")
		os.Exit(0)
	}

	if len(os.Args) > 2 {
		slog.Error("too many arguments")
		os.Exit(0)
	}

	filename := os.Args[1]

	data, err := readFileByLine(filename)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(0)
	}

	if line, err := validateInput(data); err != nil {
		if line == "" {
			slog.Error(err.Error())
			os.Exit(0)
		}

		fmt.Println(line)
		os.Exit(0)
	}

	if err = app.Run(data); err != nil {
		slog.Error(err.Error())
		os.Exit(0)
	}
}

func readFileByLine(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("can't open file: %w", err)
	}
	defer file.Close()

	data := make([]string, 0, 64)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		data = append(data, line)
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("can't read file: %w", err)
	}

	return data, nil
}

func validateInput(data []string) (string, error) {
	if len(data) < 4 {
		return "", fmt.Errorf("data is too small")
	}

	firstLine := data[0]
	tablesCount, ok := helpers.ParsePositiveInt(firstLine)
	if !ok {
		return firstLine, fmt.Errorf("third line is not valid")
	}

	secondLine := data[1]
	if !isSecondLineValid(secondLine) {
		return secondLine, fmt.Errorf("second line is not valid")
	}

	thirdLine := data[2]
	if _, ok := helpers.ParsePositiveInt(thirdLine); !ok {
		return thirdLine, fmt.Errorf("third line is not valid")
	}

	for _, line := range data[3:] {
		if !isLineValid(line, tablesCount) {
			return line, fmt.Errorf("line is not valid")
		}
	}

	return "", nil
}

func isSecondLineValid(secondLine string) bool {
	split := strings.Split(secondLine, " ")

	if len(split) != 2 {
		return false
	}

	for _, str := range split {
		_, err := helpers.ParseTime(str)
		if err != nil {
			return false
		}
	}

	return true
}

func isLineValid(line string, tablesCount int) bool {
	const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789_-"
	var threeLengthTypes = []string{"1", "3", "4"}
	var fourLengthTypes = []string{"2"}

	split := strings.Split(line, " ")
	length := len(split)
	switch length {
	case 4:
		table, ok := helpers.ParsePositiveInt(split[3])
		if !ok || table > tablesCount {
			return false
		}

		if !slices.Contains(fourLengthTypes, split[1]) {
			return false
		}

		fallthrough
	case 3:
		if length == 3 && !slices.Contains(threeLengthTypes, split[1]) {
			return false
		}

		_, err := helpers.ParseTime(split[0])
		if err != nil {
			return false
		}

		if _, ok := helpers.ParsePositiveInt(split[1]); !ok {
			return false
		}

		if !helpers.IsAllowedChars(split[2], allowedChars) {
			return false
		}
	default:
		return false
	}

	return true
}
