package src

type PromptTag struct {
	tag    string
	weight float64
}

type PromptModel struct {
	filename   string
	multiplier float64
}

type ParsedPrompt struct {
	tags      []PromptTag
	loras     []PromptModel
	hypernets []PromptModel
}

func evaluatePromptContents(contents []Prompt, currentAttention float64, defaultAttention float64, evaluated *ParsedPrompt) {
	for _, content := range contents {
		switch content.kind {
		case "pw":
			evaluatePromptContents(content.contents, currentAttention*defaultAttention, defaultAttention, evaluated)
		case "nw":
			evaluatePromptContents(content.contents, currentAttention/defaultAttention, defaultAttention, evaluated)
		case "ew":
			evaluatePromptContents(content.contents, currentAttention*content.weight, defaultAttention, evaluated)
		case "lora":
			mutiplier := content.multiplier
			if mutiplier == 0 {
				mutiplier = 1
			}
			evaluated.loras = append(evaluated.loras, PromptModel{filename: content.filename, multiplier: mutiplier})
		case "hypernet":
			mutiplier := content.multiplier
			if mutiplier == 0 {
				mutiplier = 1
			}
			evaluated.hypernets = append(evaluated.hypernets, PromptModel{filename: content.filename, multiplier: mutiplier})
		default:
			evaluated.tags = append(evaluated.tags, PromptTag{tag: content.name, weight: currentAttention})
		}
	}
}

type PromptParser struct{}

func (parser *PromptParser) ParsePrompt(prompt string) (results ParsedPrompt, err error) {
	return
}

func NewPromptParser() *PromptParser {
	return &PromptParser{}
}
func EvaluatePrompt(prompt Prompt, defaultAttention float64) ParsedPrompt {
	if defaultAttention == 0 {
		defaultAttention = 1.1
	}

	evaluated := ParsedPrompt{
		tags:      []PromptTag{},
		loras:     []PromptModel{},
		hypernets: []PromptModel{},
	}

	evaluatePromptContents(prompt.contents, 1, defaultAttention, &evaluated)

	return evaluated
}
