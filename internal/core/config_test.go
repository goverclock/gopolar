package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestNewConfig(t *testing.T) {
	assert := assert.New(t)
	cfg := NewConfig()
	assert.Equal("/Users/goverclock/.config/gopolar/gopolar.toml", cfg.filePath)
}
