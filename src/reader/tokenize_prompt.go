package reader

import "strings"

func addTokens(tokens *[]string, input *string, start *int, end *int) {
	if *end <= len(*input) && *start < *end {
		*tokens = append(*tokens, strings.Trim((*input)[*start:*end], " "))
	}
}

func tokenizeModel(tokens *[]string, input *string, start *int, end *int) {
	for {
		if *end >= len(*input) {
			break
		}

		char := (*input)[*end]
		switch char {
		case '<', ':', '>':
			addTokens(tokens, input, start, end)
			*tokens = append(*tokens, string(char))
			*start = *end + 1
		}

		if char == '>' {
			break
		}

		*end++
	}
}

func tokenizeInput(input string) (tokens []string) {
	var current int
	var index int
	for index = 0; index < len(input); index++ {
		char := input[index]
		switch char {
		case '(', ')', '[', ']', ':', ',', '|':
			addTokens(&tokens, &input, &current, &index)
			tokens = append(tokens, string(char))
			current = index + 1
		case '<':
			tokenizeModel(&tokens, &input, &current, &index)
		case ' ':
			addTokens(&tokens, &input, &current, &index)
			current = index + 1
		}
	}

	addTokens(&tokens, &input, &current, &index)

	return tokens
}
