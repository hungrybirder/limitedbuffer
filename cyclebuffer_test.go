package limitedbuffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCycleBuffer(t *testing.T) {
	assert := assert.New(t)
	const bufSize = 8
	cb := NewCycleBuffer(bufSize)
	status := cb.Status()
	assert.Equal(bufSize, status.Capacity())
}
