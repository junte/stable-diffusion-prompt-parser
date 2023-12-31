package parser

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/junte/stable-diffusion-prompt-parser/src/reader"
)

const DefaultMultiplier = 0.5
const DefaultWeight = 1

func (parser *PromptParser) escapeToken(token string) string {
	return regexp.MustCompile(`\\(.)`).ReplaceAllString(token, "$1")
}

func (parser *PromptParser) parseTagPrompt(reader *reader.TokenReader) (*prompt, error) {
	tokens := []string{}
	invalidTokens := []string{"(", ")", "[", "]", "<", ">", ":", ",", "|", ""}
	for {
		token := reader.GetToken()

		if slices.Contains(invalidTokens, token) {
			break
		}

		tokens = append(tokens, token)
		reader.NextToken()
	}

	if len(tokens) == 0 {
		return &prompt{}, errors.New("tag expected")
	}

	for i, token := range tokens {
		tokens[i] = parser.escapeToken(token)
	}

	return &prompt{
		kind:   tag,
		name:   strings.Join(tokens, " "),
		tokens: tokens,
	}, nil
}

func (parser *PromptParser) parsePositivePrompt(reader *reader.TokenReader) (*prompt, error) {
	reader.NextToken()
	contents, err := parser.parsePromptContents(reader, false)
	if err != nil {
		return &prompt{}, fmt.Errorf("%v", err)
	}

	for {
		tokens, err := reader.GetMultipleTokens(2)
		if err != nil {
			break
		}

		colon, number := tokens[0], tokens[1]
		if colon != ":" {
			break
		}

		if _, err := strconv.ParseFloat(number, 64); err == nil {
			break
		}

		reader.NextToken()

		token := reader.GetToken()
		token = strings.TrimSpace(token)
		token = strings.ReplaceAll(token, " ", "")
		token = strings.ReplaceAll(token, ",", ".")
		matches := regexp.MustCompile(`(\d+)[\.\, ]+(\d+)`).FindStringSubmatch(token)
		if len(matches) > 0 {
			token = matches[1] + "." + matches[2]
		}

		if weight, err := strconv.ParseFloat(token, 64); err == nil {
			reader.NextToken()
			return &prompt{
				kind:     customWeight,
				weight:   weight,
				contents: contents,
			}, nil
		}

		// RECOVER: (a:b)
		newContents, err := parser.parsePromptContents(reader, false)
		if err != nil {
			return &prompt{}, fmt.Errorf("%v", err)
		}

		contents = append(contents, newContents...)
	}

	if reader.GetToken() == ":" {
		reader.NextToken()
		weight, err := parser.parseNumber(reader, "weight")
		if err != nil {
			return &prompt{}, fmt.Errorf("%v", err)
		}

		if reader.GetToken() == "," {
			// RECOVER: :1,)
			reader.NextToken()
		}

		if reader.GetToken() == ")" || reader.GetToken() == "}" /* RECOVER: Expected ) but } found */ || reader.GetToken() == "" /* RECOVER: missing ) */ {
			reader.NextToken()
		}

		return &prompt{
			kind:     customWeight,
			weight:   weight,
			contents: contents,
		}, nil
	}

	if reader.GetToken() == ")" || reader.GetToken() == "}" /* RECOVER: Expected ) but } found */ || reader.GetToken() == "" /* RECOVER: missing ) */ {
		reader.NextToken()
	}

	return &prompt{
		kind:     positiveWeight,
		contents: contents,
	}, nil
}

func (parser *PromptParser) parseNegativePrompt(reader *reader.TokenReader) (*prompt, error) {
	reader.NextToken()
	contents, err := parser.parsePromptContents(reader, false)
	if err != nil {
		return &prompt{}, fmt.Errorf("%v", err)
	}

	if reader.GetToken() == "]" || reader.GetToken() == "}" || reader.GetToken() == "" {
		reader.NextToken()
	}

	return &prompt{
		kind:     negativeWeight,
		contents: contents,
	}, nil
}

func (parser *PromptParser) parseContentToken(reader *reader.TokenReader, name string) (string, error) {
	token := reader.GetToken()
	switch token {
	case "(", ")", "[", "]", "<", ">", ":", ",", "":
		return "", fmt.Errorf(fmt.Sprintf("%s expected", name))
	default:
		return parser.escapeToken(token), nil
	}
}

func (parser *PromptParser) parseFilename(reader *reader.TokenReader) (filename string, err error) {
	filename, err = parser.parseContentToken(reader, "filename")
	if err != nil {
		return "", err
	}

	reader.NextToken()

	return filename, nil
}

func (parser *PromptParser) parseNumber(reader *reader.TokenReader, name string) (number float64, err error) {
	isInt := func(s string) bool {
		i, err := strconv.ParseInt(s, 10, 64)
		return err == nil && strconv.FormatInt(i, 10) == s
	}

	// RECOVER: 1,5 (means 1.5)
	tokens, err := reader.GetMultipleTokens(3)
	if err == nil {
		first, second, third := tokens[0], tokens[1], tokens[2]
		if isInt(first) && second == "," && isInt(third) {
			reader.NextToken()
			reader.NextToken()
			reader.NextToken()
			number, err = strconv.ParseFloat(fmt.Sprintf("%s.%s", first, third), 64)
			if err != nil {
				return 0, fmt.Errorf("%v", err)
			}

			return number, nil
		}
	}

	// RECOVER: 1. 5 (means 1.5)
	tokens, err = reader.GetMultipleTokens(2)
	if err == nil {
		first, second := tokens[0], tokens[1]
		if strings.HasSuffix(first, ".") && isInt(first[:len(first)-1]) && isInt(second) {
			reader.NextToken()
			reader.NextToken()
			number, err = strconv.ParseFloat(fmt.Sprintf("%s%s", first, second), 64)
			if err != nil {
				return 0, fmt.Errorf("%v", err)
			}

			return number, nil
		}
	}

	switch reader.GetToken() {
	case ")", ">", ":":
		switch name {
		case "multiplier":
			return DefaultMultiplier, nil
		case "weight":
			return DefaultWeight, nil
		}
	}

	token, err := parser.parseContentToken(reader, name)
	if err != nil {
		return 0, err
	}

	token = strings.TrimSpace(token)
	token = strings.ReplaceAll(token, " ", "")
	token = strings.ReplaceAll(token, ",", ".")
	matches := regexp.MustCompile(`(\d+)[\.\, ]+(\d+)`).FindStringSubmatch(token)
	if len(matches) > 0 {
		token = matches[1] + "." + matches[2]
	}

	number, err = strconv.ParseFloat(token, 64)
	if err != nil {
		switch name {
		case "multiplier":
			number = DefaultMultiplier
		case "weight":
			number = DefaultWeight
		}
	}

	reader.NextToken()

	return number, nil
}

func (parser *PromptParser) parseAnglePrompt(reader *reader.TokenReader, kind string) (*prompt, error) {
	reader.NextToken()
	reader.NextToken()
	if reader.GetToken() != ":" {
		return &prompt{}, errors.New(": expected")
	}

	reader.NextToken()
	filename, err := parser.parseFilename(reader)
	if err != nil {
		return &prompt{}, err
	}

	if reader.GetToken() == ":" {
		reader.NextToken()
		multiplier, err := parser.parseNumber(reader, "multiplier")
		if err != nil {
			return &prompt{}, err
		}

		if reader.GetToken() == ">" {
			reader.NextToken()
			return &prompt{
				kind:       kind,
				filename:   filename,
				multiplier: multiplier,
			}, nil
		}
	}

	if reader.GetToken() == ">" {
		reader.NextToken()
		return &prompt{
			kind:       kind,
			filename:   filename,
			multiplier: DefaultMultiplier,
		}, nil
	}
	reader.NextToken()

	return &prompt{}, errors.New("> expected")
}

func (parser *PromptParser) parsePromptContent(reader *reader.TokenReader, topLevel bool) (*prompt, error) {
	token := reader.GetToken()
	switch token {
	case "(":
		return parser.parsePositivePrompt(reader)
	case "[":
		return parser.parseNegativePrompt(reader)
	case "<":
		tokens, err := reader.GetMultipleTokens(2)
		if err == nil {
			// RECOVER: (topLevel === false) A <a:b:c> cannot be nested in other prompt
			modelName := tokens[1]
			switch modelName {
			case lora, hypernet:
				return parser.parseAnglePrompt(reader, modelName)
			default:
				// RECOVER: unknown model name
				reader.NextToken()
				return &prompt{}, nil
			}
		}
		reader.NextToken()
		return &prompt{}, nil
	case ",":
		return &prompt{}, errors.New("prompt expected")
	default:
		tagPrompt, err := parser.parseTagPrompt(reader)
		if err != nil {
			return &prompt{}, err
		}

		if topLevel && reader.GetToken() == ":" {
			tokens, err := reader.GetMultipleTokens(2)
			if err == nil {
				numberText := tokens[1]
				f, err := strconv.ParseFloat(numberText, 64)
				if err == nil && !math.IsNaN(f) {
					reader.NextToken()
					weight, err := parser.parseNumber(reader, "weight")
					if err == nil {
						return &prompt{
							kind:     customWeight,
							weight:   weight,
							contents: []*prompt{tagPrompt},
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

func (parser *PromptParser) parsePromptContents(reader *reader.TokenReader, topLevel bool) (contents []*prompt, err error) {
	for {
		token := reader.GetToken()
		switch token {
		case ",":
			reader.NextToken()
		case ":", ")", "]", ">", "|", "":
			return contents, nil
		default:
			content, err := parser.parsePromptContent(reader, topLevel)
			if err != nil {
				return nil, fmt.Errorf("%v", err)
			}

			if !reflect.DeepEqual(prompt{}, *content) {
				contents = append(contents, content)
			}
		}
	}
}

func (parser *PromptParser) parse(input string) (*prompt, error) {
	input = strings.ReplaceAll(input, "|", ",")
	prompt := &prompt{}
	reader := reader.NewTokenReader(input)

	for {
		switch reader.GetToken() {
		case ")", "]", ">", ":":
			reader.NextToken()
			continue
		case "":
			return prompt, nil
		default:
			contents, err := parser.parsePromptContents(reader, true)
			if err != nil {
				return prompt, fmt.Errorf("%v", err)
			}

			prompt.contents = append(prompt.contents, contents...)
		}
	}
}
