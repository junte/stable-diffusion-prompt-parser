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
	Tag    string  `json:"tag"`
	Weight float64 `json:"weight"`
}

type PromptModel struct {
	Filename   string  `json:"filename"`
	Multiplier float64 `json:"multiplier"`
}

type ParsedPrompt struct {
	Tags      []*PromptTag   `json:"tags"`
	Loras     []*PromptModel `json:"loras"`
	Hypernets []*PromptModel `json:"hypernets"`
}
