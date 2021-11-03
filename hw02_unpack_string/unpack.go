package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

const (
	Nil       = 48
	Backslash = 92
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder

	runes := []rune(input)

	for i := 0; i < len(runes); i++ {
		if runes[i] == Backslash {
			stringAdd, indexAdd, err := processSlash(runes, i)
			if err != nil {
				return "", ErrInvalidString
			}

			result.WriteString(stringAdd)
			i += indexAdd
			continue
		}

		if !isDigit(runes[i]) {
			stringAdd, indexAdd := processString(runes, i)

			result.WriteString(stringAdd)
			i += indexAdd
			continue
		}

		if isDigit(runes[i]) {
			return "", ErrInvalidString
		}
	}

	return result.String(), nil
}

func isDigit(r rune) bool {
	_, err := strconv.Atoi(string(r))

	return err == nil
}

func processString(input []rune, index int) (string, int) {
	currentString := string(input[index])

	if index+1 == len(input) {
		return currentString, 0
	}

	nextRune := input[index+1]

	if !isDigit(nextRune) {
		return currentString, 0
	}

	if nextRune == Nil {
		return "", 1
	}

	num, _ := strconv.Atoi(string(nextRune))
	stringToAdd := strings.Repeat(currentString, num)
	return stringToAdd, 1
}

func processSlash(input []rune, index int) (string, int, error) {
	currentString := string(input[index])

	if index+1 == len(input) {
		return "", 0, ErrInvalidString
	}

	nextRune := input[index+1]

	if isDigit(nextRune) && index+2 == len(input) {
		nextString := string(input[index+1])
		return nextString, 1, nil
	}

	if nextRune == Backslash && index+2 == len(input) {
		return `\`, 0, nil
	}

	nextNextRune := input[index+2]

	if isDigit(nextRune) && !isDigit(nextNextRune) {
		nextString := string(input[index+1])
		return nextString, 1, nil
	}

	if isDigit(nextRune) && isDigit(nextNextRune) {
		nextString := string(input[index+1])

		num, _ := strconv.Atoi(string(nextNextRune))
		stringToAdd := strings.Repeat(nextString, num)
		return stringToAdd, 2, nil
	}

	if nextRune == Backslash && !isDigit(nextNextRune) {
		return `\`, 1, nil
	}

	num, _ := strconv.Atoi(string(nextNextRune))
	stringToAdd := strings.Repeat(currentString, num)
	return stringToAdd, 2, nil
}
