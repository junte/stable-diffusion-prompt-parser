package src

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenReader(t *testing.T) {
	tokenizedInput := TokenizePrompt("(abc:1.5)")
	reader := NewTokenReader(tokenizedInput)

	assert.Equal(t, "(", reader.getToken())
	reader.nextToken()
	assert.Equal(t, "abc", reader.getToken())

	result, err := reader.getMultipleTokens(3)
	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"abc", ":", "1.5"}, result)
}
