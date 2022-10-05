package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViper(t *testing.T) {
	c, err := LoadConfig("..")
	assert.NotEmpty(t, c.DBDriver)
	assert.NoError(t, err)
	assert.NotEmpty(t, c.Duration)
}
