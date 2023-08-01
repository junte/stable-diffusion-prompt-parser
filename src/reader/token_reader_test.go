package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenReader(t *testing.T) {
	reader := NewTokenReader("(abc:1.5)")

	assert.Equal(t, "(", reader.GetToken())
	reader.NextToken()
	assert.Equal(t, "abc", reader.GetToken())

	result, err := reader.GetMultipleTokens(3)
	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"abc", ":", "1.5"}, result)

	_, err = reader.GetMultipleTokens(10)
	assert.EqualError(t, err, "count out of range")
}
