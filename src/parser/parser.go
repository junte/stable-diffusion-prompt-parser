package parser

type PromptParser struct{}

func NewPromptParser() *PromptParser {
	return &PromptParser{}
}

func (parser *PromptParser) ParsePrompt(input string) (*ParsedPrompt, error) {
	prompt, err := parser.parse(input)
	if err != nil {
		return &ParsedPrompt{}, err
	}

	return parser.evaluate(prompt), nil
}

func (parser *PromptParser) BeautifyPrompt(input string) (string, error) {
	prompt, err := parser.parse(input)
	if err != nil {
		return "", err
	}

	return parser.toString(prompt), nil
}
