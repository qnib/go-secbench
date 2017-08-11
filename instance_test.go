package secbench

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInstance_String(t *testing.T) {
	i := NewInstance("WARN", "Have done something")
	exp := "     | WARN  || Have done something"
	assert.Equal(t, exp, i.String())

}
