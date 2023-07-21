package main

import (
	"strings"
)

type Token struct {
	value    string
	position int
}

func TokenizePrompt(input string) []Token {
	var tokens []Token
	var current string
	var startPosition, stopPosition int
	var escaping bool

	submit := func() {
		if current != "" {
			tokens = append(tokens, Token{current, startPosition})
			current = ""
		}
		startPosition = stopPosition
	}

	skip := func() {
		stopPosition++
	}

	append := func(char string) {
		current += char
		skip()
	}

	input = strings.ReplaceAll(input, "，", ",")
	input = strings.ReplaceAll(input, "：", ":")
	input = strings.ReplaceAll(input, "（", "(")
	input = strings.ReplaceAll(input, "）", ")")

	for _, char := range input {
		if escaping {
			escaping = false
			append(string(char))
		} else {
			switch char {
			case '\\':
				escaping = true
				append(string(char))
			case '(', ')', '[', ']', '<', '>', ':', ',', '|':
				submit()
				append(string(char))
				submit()
			case ' ':
				skip()
				submit()
			default:
				append(string(char))
			}
		}
	}

	submit()
	return tokens
}
