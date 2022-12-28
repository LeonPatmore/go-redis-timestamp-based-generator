package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	list := []string{"123", "1", "12345"}
	lengths := Map(list, func(item string) int { return len(item) })

	assert.Len(t, lengths, 3)
	assert.Equal(t, 3, lengths[0])
	assert.Equal(t, 1, lengths[1])
	assert.Equal(t, 5, lengths[2])
}
