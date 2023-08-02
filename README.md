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
```bash
$ make build-mac

## Supported commands
# make build-linux
# make build-linux64
# make build-mac
# make build-mac32
```

### Usage
```bash
$ ./bin/parse < <(echo "landscape from the Moon, (realistic, detailed:1.5), <lora:file>, <hypernet:file:1.5>")

# or

$ ./bin/parse 
$ "landscape,,,, moon, ( realistic,detailed:1, 5), <hypernet:file:1. 5>"
```