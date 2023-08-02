package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/junte/stable-diffusion-prompt-parser/src/parser"
)

type Output struct {
	Evaluated  *parser.ParsedPrompt `json:"evaluated"`
	Beautified string               `json:"beautified"`
}

func toIndentedJson(i *Output, prefix string, indent string) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent(prefix, indent)
	err := encoder.Encode(i)

	return bytes.TrimRight(buffer.Bytes(), "\n"), err
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	input := scanner.Text()
	parser := parser.NewPromptParser()

	parsed, err := parser.ParsePrompt(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	beautified, err := parser.BeautifyPrompt(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	output := Output{
		Evaluated:  parsed,
		Beautified: beautified,
	}

	marshalled, err := toIndentedJson(&output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, string(marshalled))
}
