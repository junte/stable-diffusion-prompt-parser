package src

import "fmt"

func printPromptContents(contents []Prompt) string {
	printed := ""
	var lastPromptIsTag bool

	for _, content := range contents {
		if content.kind == "tag" && lastPromptIsTag {
			printed += ", "
		} else if printed != "" {
			printed += " "
		}

		switch content.kind {
		case "pw":
			printed += "(" + printPromptContents(content.contents) + ")"
		case "nw":
			printed += "[" + printPromptContents(content.contents) + "]"
		case "ew":
			weightText := content.weightText
			if weightText == "" {
				weightText = fmt.Sprintf("%v", content.weight)
			}
			printed += "(" + printPromptContents(content.contents) + ":" + weightText + ")"
		case "lora", "hypernet":
			multiplierText := content.multiplierText
			if multiplierText == "" {
				multiplierText = fmt.Sprintf("%v", content.multiplier)
			}
			printed += "<" + content.kind + ":" + content.filename + ":" + multiplierText + ">"
		default:
			printed += content.name
		}
		lastPromptIsTag = content.kind == "tag"
	}

	return printed
}

func PrintPrompt(prompt Prompt) string {
	return printPromptContents(prompt.contents)
}
