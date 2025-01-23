/* Создать Go-функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы/руны.

Например:
"a4bc2d5e" => "aaaabccddddde"
"abcd" => "abcd"
"45" => "" (некорректная строка)
"" => ""
*/

package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"unicode"
)

func Unpack(s string) (string, error) {
	var result []rune
	var lastRune rune
	var haveLast bool

	for _, r := range s {
		if unicode.IsDigit(r) {
			if !haveLast {
				return "", errors.New("строка начинается с цифры или содержит неправильный формат")
			}
			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}
			for i := 0; i < count; i++ {
				result = append(result, lastRune)
			}
			haveLast = false
		} else {
			if haveLast {
				result = append(result, lastRune)
			}
			lastRune = r
			haveLast = true
		}
	}

	if haveLast {
		result = append(result, lastRune)
	}

	return string(result), nil
}

func main() {
	examples := []string{
		"a4bc2d5e",
		"abcd",
		"45",
		"",
		"a0b1c3",
	}

	for _, ex := range examples {
		unpacked, err := Unpack(ex)
		if err != nil {
			log.Printf("Unpack(%q) => ERROR: %v\n", ex, err)
		} else {
			fmt.Printf("Unpack(%q) => %q\n", ex, unpacked)
		}
	}
}
