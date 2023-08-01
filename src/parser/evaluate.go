package parser

func (parser *PromptParser) evaluatePromptContents(contents []*prompt, currentWeight float64, weightMultiplier float64, evaluated *ParsedPrompt) {
	for _, content := range contents {
		switch content.kind {
		case positiveWeight:
			parser.evaluatePromptContents(content.contents, currentWeight*weightMultiplier, weightMultiplier, evaluated)
		case negativeWeight:
			parser.evaluatePromptContents(content.contents, currentWeight/weightMultiplier, weightMultiplier, evaluated)
		case customWeight:
			parser.evaluatePromptContents(content.contents, currentWeight*content.weight, weightMultiplier, evaluated)
		case lora:
			mutiplier := content.multiplier
			if mutiplier == 0 {
				mutiplier = 1
			}
			evaluated.Loras = append(evaluated.Loras, &PromptModel{Filename: content.filename, Multiplier: mutiplier})
		case hypernet:
			mutiplier := content.multiplier
			if mutiplier == 0 {
				mutiplier = 1
			}
			evaluated.Hypernets = append(evaluated.Hypernets, &PromptModel{Filename: content.filename, Multiplier: mutiplier})
		default:
			evaluated.Tags = append(evaluated.Tags, &PromptTag{Tag: content.name, Weight: currentWeight})
		}
	}
}

func (parser *PromptParser) evaluate(prompt *prompt) *ParsedPrompt {
	currentWeight, weightMultiplier := 1.0, 1.1
	evaluated := &ParsedPrompt{}
	parser.evaluatePromptContents(prompt.contents, currentWeight, weightMultiplier, evaluated)

	return evaluated
}
