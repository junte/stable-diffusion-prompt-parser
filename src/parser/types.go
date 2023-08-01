package parser

type prompt struct {
	kind       string
	name       string
	filename   string
	multiplier float64
	weight     float64
	tokens     []string
	contents   []*prompt
}

type PromptTag struct {
	Tag    string
	Weight float64
}

type PromptModel struct {
	Filename   string
	Multiplier float64
}

type ParsedPrompt struct {
	Tags      []*PromptTag
	Loras     []*PromptModel
	Hypernets []*PromptModel
}
