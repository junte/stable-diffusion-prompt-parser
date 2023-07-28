package reader

import (
	"strings"
)

func tokenizeInput(input string) (tokens []Token) {
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

	add := func(char string) {
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
			add(string(char))
		} else {
			switch char {
			case '\\':
				escaping = true
				add(string(char))
			case '(', ')', '[', ']', '<', '>', ':', ',', '|':
				submit()
				add(string(char))
				submit()
			case ' ':
				skip()
				submit()
			default:
				add(string(char))
			}
		}
	}

	submit()
	return tokens
}
