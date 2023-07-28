package parser

import "fmt"

func (parser *PromptParser) contentsToString(contents []Prompt) (result string) {
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
				result += ":" + fmt.Sprintf("%v", content.weight)
			}
			result += ")"
		case lora, hypernet:
			result += "<" + content.kind + ":" + content.filename
			if content.multiplier != 0 {
				result += ":" + fmt.Sprintf("%v", content.multiplier)
			}
			result += ">"
		default:
			result += content.name
		}

		lastPromptIsTag = content.kind == tag
	}

	return result
}

func (parser *PromptParser) toString(prompt *Prompt) string {
	return parser.contentsToString(prompt.contents)
}
