# stable-diffusion-prompt-parser

Stable Diffusion Prompt Parser

## Installation

```bash
$ go get -u github.com/junte/stable-diffusion-prompt-parser 
```

## Supported prompt syntax
- tags with default weight (`dog`, `dog,cat`)
- tags with increased weight (`(dog)`, `((dog))`)
- tags with decreased weight (`[dog]`, `[[dog]]`)
- tags with custom weight (`(dog:1.5)`, `(cat:0.5)`)
- lora models with default and custom multiplier (`<lora:filename>`, `<lora:filename:1.5>`)
- hypernet models with default and custom multiplier (`<hypernet:filename>`, `<hypernet:filename:1.5>`)

## Examples

### Parse prompt

``` go
package main

import "github.com/junte/stable-diffusion-prompt-parser/src/parser"

func main() {
    prompt := "landscape from the Moon, (realistic, detailed:1.5), <lora:file>, <hypernet:file:1.5>"
    parser := parser.NewPromptParser()
    parsed, err := parser.ParsePrompt(prompt)
}
```
parsed:

```json
{
  "Tags": [
    {
      "Tag": "landscape from the Moon",
      "Weight": 1
    },
    {
      "Tag": "realistic",
      "Weight": 1.5
    },
    {
      "Tag": "detailed",
      "Weight": 1.5
    }
  ],
  "Loras": [
    {
      "Filename": "file",
      "Multiplier": 1
    }
  ],
  "Hypernets": [
    {
      "Filename": "file",
      "Multiplier": 1.5
    }
  ]
}
```

### Beautify prompt

```go
package main

import "github.com/junte/stable-diffusion-prompt-parser/src/parser"

func main() {
    prompt := "landscape,,,, moon, ( realistic,detailed:1, 5), <hypernet:file:1. 5>"
    parser := parser.NewPromptParser()
    beautified, err := parser.BeautifyPrompt(prompt)
}
```
beautified:
```
landscape, moon (realistic, detailed:1.5) <hypernet:file:1.5>
```

## Build
Use following make rules for build binary and run 
```bash
$ make build
$ ./bin/linux-amd64/parse < <(echo "landscape from the Moon")
$ ./bin/mac-arm64/parse < <(echo "landscape from the Moon")
```
Or use following command to build binary for desirable platform (use valid combinations of $GOOS and $GOARCH from here: https://go.dev/doc/install/source#environment)

```bash
$ GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o bin/parse main.go
```
