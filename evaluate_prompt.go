package main

type EvaluatedTag struct {
	tag    string
	weight float64
}

type EvaluatedAngle struct {
	filename   string
	multiplier float64
}

type EvaluatedPrompt struct {
	tags      []EvaluatedTag
	loras     []EvaluatedAngle
	hypernets []EvaluatedAngle
}

func evaluatePromptContents(contents []Prompt, currentAttention float64, defaultAttention float64, evaluated *EvaluatedPrompt) {
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
			evaluated.loras = append(evaluated.loras, EvaluatedAngle{filename: content.filename, multiplier: mutiplier})
		case "hypernet":
			mutiplier := content.multiplier
			if mutiplier == 0 {
				mutiplier = 1
			}
			evaluated.hypernets = append(evaluated.hypernets, EvaluatedAngle{filename: content.filename, multiplier: mutiplier})
		default:
			evaluated.tags = append(evaluated.tags, EvaluatedTag{tag: content.name, weight: currentAttention})
		}
	}
}

func EvaluatePrompt(prompt Prompt, defaultAttention float64) EvaluatedPrompt {
	if defaultAttention == 0 {
		defaultAttention = 1.1
	}

	evaluated := EvaluatedPrompt{
		tags:      []EvaluatedTag{},
		loras:     []EvaluatedAngle{},
		hypernets: []EvaluatedAngle{},
	}

	evaluatePromptContents(prompt.contents, 1, defaultAttention, &evaluated)

	return evaluated
}
