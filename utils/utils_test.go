package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnName(t *testing.T) {
	assert.Equal(t, `A`, ColumnName(0))
	assert.Equal(t, `Z`, ColumnName(25))
	assert.Equal(t, `AA`, ColumnName(26))
	assert.Equal(t, `AZ`, ColumnName(51))
	assert.Equal(t, `BA`, ColumnName(52))
}
