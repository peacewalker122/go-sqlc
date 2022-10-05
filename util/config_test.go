package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const path = "/home/servumtopia/CODE/GO/sqlc"

func TestViper(t *testing.T) {
	c, err := LoadConfig(path)
	assert.NotEmpty(t, c.DBDriver)
	assert.NoError(t, err)
	assert.NotEmpty(t, c.Duration)
}
