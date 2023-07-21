package main

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Prompt struct {
	kind           string
	name           string
	filename       string
	multiplier     float64
	multiplierText string
	weight         float64
	weightText     string
	tokens         []string
	contents       []Prompt
}

func escapeToken(token string) string {
	return regexp.MustCompile(`\\(.)`).ReplaceAllString(token, "$1")
}

func parseTagPrompt(reader *TokenReader) (Prompt, error) {
	tokens := []string{}

parsing:
	for {
		token := reader.getToken()
		switch token {
		case "(", ")", "[", "]", "<", ">", ":", ",", "|", "":
			break parsing
		default:
			tokens = append(tokens, token)
			reader.nextToken()
		}
	}

	if len(tokens) == 0 {
		return Prompt{}, errors.New("tag expected")
	} else {
		for i, token := range tokens {
			tokens[i] = escapeToken(token)
		}

		return Prompt{
			kind:   "tag",
			name:   strings.Join(tokens, " "),
			tokens: tokens,
		}, nil
	}
}

func parsePositivePrompt(reader *TokenReader) (Prompt, error) {
	reader.nextToken()

	contents, err := parsePromptContents(reader, false)
	if err != nil {
		return Prompt{}, fmt.Errorf("%v", err)
	}

	for {
		tokens, err := reader.getMultipleTokens(2)
		if err != nil {
			break
		} else {
			colon, number := tokens[0], tokens[1]
			if colon != ":" {
				break
			}
			if _, err := strconv.ParseFloat(number, 64); err == nil {
				break
			}

			// RECOVER: (a:b)
			reader.nextToken()
			newContents, err := parsePromptContents(reader, false)
			if err != nil {
				return Prompt{}, fmt.Errorf("%v", err)
			}
			contents = append(contents, newContents...)
		}
	}

	if reader.getToken() == ":" {
		reader.nextToken()
		weight, weightText, err := parseNumber(reader, "Weight")
		if err != nil {
			return Prompt{}, fmt.Errorf("%v", err)
		}
		if reader.getToken() == "," {
			// RECOVER: :1,)
			reader.nextToken()
		}

		if reader.getToken() == ")" || reader.getToken() == "}" /* RECOVER: Expected ) but } found */ || reader.getToken() == "" /* RECOVER: missing ) */ {
			reader.nextToken()
		}

		return Prompt{
			kind:       "ew",
			weight:     weight,
			weightText: weightText,
			contents:   contents,
		}, nil
	}

	if reader.getToken() == ")" || reader.getToken() == "}" /* RECOVER: Expected ) but } found */ || reader.getToken() == "" /* RECOVER: missing ) */ {
		reader.nextToken()
	}

	return Prompt{
		kind:     "pw",
		contents: contents,
	}, nil
}

func parseNegativePrompt(reader *TokenReader) (Prompt, error) {
	reader.nextToken()
	contents, err := parsePromptContents(reader, false)

	if err != nil {
		return Prompt{}, fmt.Errorf("%v", err)
	}

	if reader.getToken() == "]" || reader.getToken() == "}" || reader.getToken() == "" {
		reader.nextToken()
	}

	return Prompt{
		kind:     "nw",
		contents: contents,
	}, nil
}

func parseContentToken(reader *TokenReader, name string) (string, error) {
	token := reader.getToken()
	switch token {
	case "(", ")", "[", "]", "<", ">", ":", ",", "":
		return "", fmt.Errorf(fmt.Sprintf("%s expected", name))
	default:
		return escapeToken(token), nil
	}
}

func parseFilename(reader *TokenReader) (string, error) {
	result, err := parseContentToken(reader, "Filename")

	if err != nil {
		return "", err
	}

	reader.nextToken()
	return result, nil
}

func parseNumber(reader *TokenReader, name string) (float64, string, error) {
	isInt := func(s string) bool {
		i, err := strconv.ParseInt(s, 10, 64)
		return err == nil && strconv.FormatInt(i, 10) == s
	}

	// RECOVER: 1,5 (means 1.5)
	tokens, err := reader.getMultipleTokens(3)
	if err == nil {
		first, second, third := tokens[0], tokens[1], tokens[2]
		if isInt(first) && second == "," && isInt(third) {
			reader.nextToken()
			reader.nextToken()
			reader.nextToken()
			numberText := fmt.Sprintf("%s.%s", first, third)
			f, err := strconv.ParseFloat(numberText, 64)
			if err != nil {
				return 0, "", fmt.Errorf("%v", err)
			}

			return f, numberText, nil
		}
	}

	// RECOVER: 1. 5 (means 1.5)
	tokens, err = reader.getMultipleTokens(2)
	if err == nil {
		first, second := tokens[0], tokens[1]
		if strings.HasSuffix(first, ".") && isInt(first[:len(first)-1]) && isInt(second) {
			reader.nextToken()
			reader.nextToken()
			numberText := fmt.Sprintf("%s%s", first, second)
			f, err := strconv.ParseFloat(numberText, 64)
			if err != nil {
				return 0, "", fmt.Errorf("%v", err)
			}

			return f, numberText, nil
		}
	}

	numberText, err := parseContentToken(reader, name)
	if err != nil {
		return 0, "", err
	}

	f, err := strconv.ParseFloat(numberText, 64)
	if err != nil {
		return 0, "", fmt.Errorf(fmt.Sprintf("Incorrect %s format: %s", strings.ToLower(name), numberText))
	}
	reader.nextToken()

	return f, numberText, nil
}

func parseAnglePrompt(reader *TokenReader, kind string) (Prompt, error) {
	reader.nextToken()
	reader.nextToken()
	if reader.getToken() != ":" {
		return Prompt{}, errors.New(": expected")
	}
	reader.nextToken()

	filename, err := parseFilename(reader)
	if err != nil {
		return Prompt{}, err
	}

	if reader.getToken() == ":" {
		reader.nextToken()
		multiplier, multiplierText, err := parseNumber(reader, "Multiplier")
		if err != nil {
			return Prompt{}, err
		}

		if reader.getToken() == ">" {
			reader.nextToken()
			return Prompt{
				kind:           kind,
				filename:       filename,
				multiplier:     multiplier,
				multiplierText: multiplierText,
			}, nil
		}
	}

	if reader.getToken() == ">" {
		reader.nextToken()
		return Prompt{
			kind:     kind,
			filename: filename,
		}, nil
	}
	reader.nextToken()

	return Prompt{}, errors.New("> expected")
}

func parsePromptContent(reader *TokenReader, topLevel bool) (Prompt, error) {
	token := reader.getToken()
	switch token {
	case "(":
		return parsePositivePrompt(reader)
	case "[":
		return parseNegativePrompt(reader)
	case "<":
		tokens, err := reader.getMultipleTokens(2)
		if err == nil {
			// RECOVER: (topLevel === false) A <a:b:c> cannot be nested in other prompt
			modelName := tokens[1]
			switch modelName {
			case "lora", "hypernet":
				return parseAnglePrompt(reader, modelName)
			default:
				// RECOVER: unknown model name
				reader.nextToken()
				return Prompt{}, nil
			}
		} else {
			reader.nextToken()
			return Prompt{}, nil
		}
	case ",":
		return Prompt{}, errors.New("Prompt expected")
	default:
		tagPrompt, err := parseTagPrompt(reader)
		if err != nil {
			return Prompt{}, err
		}

		if topLevel && reader.getToken() == ":" {
			tokens, err := reader.getMultipleTokens(2)
			if err == nil {
				numberText := tokens[1]
				f, err := strconv.ParseFloat(numberText, 64)
				if err == nil && !math.IsNaN(f) {
					reader.nextToken()
					weight, weightText, err := parseNumber(reader, "Weight")
					if err == nil {
						return Prompt{
							kind:       "ew",
							weight:     weight,
							weightText: weightText,
							contents:   []Prompt{tagPrompt},
						}, nil
					}
				}
			}
			return tagPrompt, nil
		} else {
			return tagPrompt, nil
		}
	}
}

func parsePromptContents(reader *TokenReader, topLevel bool) ([]Prompt, error) {
	contents := make([]Prompt, 0)
	for {
		token := reader.getToken()
		switch token {
		case ",":
			reader.nextToken()
		case ":", ")", "]", ">", "|", "":
			return contents, nil
		default:
			content, err := parsePromptContent(reader, topLevel)
			if err != nil {
				return nil, fmt.Errorf("%v", err)
			}

			if !reflect.DeepEqual(content, Prompt{}) {
				contents = append(contents, content)
			}
		}
	}
}

func ParsePrompt(input string) ([]Prompt, error) {
	reader := NewTokenReader(TokenizePrompt(input))
	prompts := make([]Prompt, 0)
	newPrompt := false

parsing:
	for {
		newPrompt = false

		for {
			token := reader.getToken()
			switch token {
			case "|":
				newPrompt = true
				reader.nextToken()
			case ")", "]", ">", ":":
				reader.nextToken()
				continue parsing
			case "":
				return prompts, nil
			}
			break
		}

		contents, err := parsePromptContents(reader, true)
		if err != nil {
			return []Prompt{}, fmt.Errorf("%v", err)
		}

		if newPrompt || len(prompts) == 0 {
			prompts = append(prompts, Prompt{
				kind:     "prompt",
				contents: contents,
			})
		} else {
			prompts[len(prompts)-1].contents = append(prompts[len(prompts)-1].contents, contents...)
		}
	}
}
