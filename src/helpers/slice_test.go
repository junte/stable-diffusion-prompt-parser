package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConatins(t *testing.T) {
	haystack := []string{"abc", "xyz"}
	assert.True(t, Contains(haystack, "abc"))
	assert.False(t, Contains(haystack, "mno"))
}
