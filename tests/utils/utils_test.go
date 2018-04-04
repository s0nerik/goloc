package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/s0nerik/goloc/utils"
)

func TestColumnName(t *testing.T) {
	assert.Equal(t, `A`, utils.ColumnName(0))
	assert.Equal(t, `Z`, utils.ColumnName(25))
	assert.Equal(t, `AA`, utils.ColumnName(26))
	assert.Equal(t, `AZ`, utils.ColumnName(51))
	assert.Equal(t, `BA`, utils.ColumnName(52))
}