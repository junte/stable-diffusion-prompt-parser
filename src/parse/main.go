package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/junte/stable-diffusion-prompt-parser/src/parser"
)

func main() {

	args := os.Args
	if len(args) > 1 {
		parser := parser.NewPromptParser()
		parsed, err := parser.ParsePrompt(args[1])
		if err != nil {
			fmt.Println(err)
		} else {
			marshalled, err := json.MarshalIndent(parsed, "", "  ")
			if err == nil {
				fmt.Println(string(marshalled))
			}
		}
	} else {
		fmt.Println("Please provide prompt to parse.")
	}
}
