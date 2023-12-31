package parser

import (
	"fmt"
	"regexp"
)

func truncateZero(input string) string {
	if input[0] == '0' {
		return input[1:]
	}

	return input
}

func (parser *PromptParser) contentsToString(contents []*prompt) (result string) {
	var lastPromptIsTag bool

	for _, content := range contents {
		if content.kind == tag && lastPromptIsTag {
			result += ", "
		} else if result != "" {
			result += " "
		}

		switch content.kind {
		case positiveWeight:
			result += "(" + parser.contentsToString(content.contents) + ")"
		case negativeWeight:
			result += "[" + parser.contentsToString(content.contents) + "]"
		case customWeight:
			result += "(" + parser.contentsToString(content.contents)
			if content.weight != 0 {
				result += ":" + truncateZero(fmt.Sprintf("%v", content.weight))
			}
			result += ")"
		case lora, hypernet:
			result += "<" + content.kind + ":" + content.filename
			if content.multiplier != 0 {
				result += ":" + truncateZero(fmt.Sprintf("%v", content.multiplier))
			}
			result += ">"
		default:
			result += content.name
		}

		lastPromptIsTag = content.kind == tag
	}

	return result
}

func (parser *PromptParser) toString(prompt *prompt) string {
	result := parser.contentsToString(prompt.contents)

	regex := regexp.MustCompile(`( [<(\[])`)
	result = regex.ReplaceAllString(result, ",$1")

	regex = regexp.MustCompile(`([>)\]] )`)
	result = regex.ReplaceAllStringFunc(result, func(entry string) string {
		return entry[:1] + ", "
	})

	return result
}
