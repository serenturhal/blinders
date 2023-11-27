package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLanguage(t *testing.T) {
	unknown := NewLanguage("invalid", "unknown")
	assert.Equal(t, "unknown", unknown.Name)
	assert.Equal(t, "invalid", unknown.Code)
}

func TestLevel(t *testing.T) {
	const (
		beginner     = iota
		intermediate = iota
		advanced     = iota
		unknow       = iota
	)

	assert.Equal(t, Level(beginner).String(), Beginner.String())
	assert.Equal(t, Level(intermediate).String(), Intermediate.String())
	assert.Equal(t, Level(advanced).String(), Advanced.String())
	assert.Equal(t, Level(unknow).String(), "Unknown")
}
