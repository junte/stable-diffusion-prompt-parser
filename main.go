package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/junte/stable-diffusion-prompt-parser/src/parser"
)

type Output struct {
	Evaluated  *parser.ParsedPrompt `json:"evaluated"`
	Beautified string               `json:"beautified"`
	Cleaned    string               `json:"cleaned"`
}

func toIndentedJson(output *Output, prefix string, indent string) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent(prefix, indent)

	// intialize slices to get rid of nulls in json
	if output.Evaluated.Tags == nil {
		output.Evaluated.Tags = make([]*parser.PromptTag, 0)
	}

	if output.Evaluated.Hypernets == nil {
		output.Evaluated.Hypernets = make([]*parser.PromptModel, 0)
	}

	if output.Evaluated.Loras == nil {
		output.Evaluated.Loras = make([]*parser.PromptModel, 0)
	}

	err := encoder.Encode(output)

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

	regex := regexp.MustCompile(` ?<.*> ?`)
	output := Output{
		Evaluated:  parsed,
		Beautified: beautified,
		Cleaned:    regex.ReplaceAllString(beautified, ", "),
	}

	marshalled, err := toIndentedJson(&output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, string(marshalled))
}
