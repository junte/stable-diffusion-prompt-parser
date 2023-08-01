package main

import (
	"fmt"
	"os"

	"github.com/junte/stable-diffusion-prompt-parser/src/parser"
)

func main() {
	args := os.Args
	if len(args) > 1 {
		parser := parser.NewPromptParser()
		beautified, err := parser.BeautifyPrompt(args[1])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(beautified)
		}
	} else {
		fmt.Println("Please provide prompt to beautify.")
	}
}
