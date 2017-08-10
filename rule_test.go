package secbench

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var (
	ModeMap = map[string]int{
		"NOTE": 0,
		"INFO": 1,
		"PASS": 2,
		"WARN": 3,
	}
)

func TestNewRule(t *testing.T) {
	r := NewRule("4.1", "Do something", "PASS")
	assert.Equal(t, "4.1", r.Num)
	assert.Equal(t, "Do something", r.Description)
	assert.Equal(t, PASS, r.CurrentMode)
	assert.Equal(t, NIL, r.DesiredMode)
}

func TestRule_AddCurrentMode(t *testing.T) {
	r := NewRule("4.1", "Do something", "PASS")
	r.AddCurrentMode("WARN")
	assert.Equal(t, "4.1", r.Num)
	assert.Equal(t, "Do something", r.Description)
	assert.Equal(t, WARN, r.CurrentMode)
	assert.Equal(t, NIL, r.DesiredMode)
}

func TestRule_AddDescription(t *testing.T) {
	r := NewRule("4.1", "Do something", "PASS")
	r.AddDescription("Do something else")
	assert.Equal(t, "4.1", r.Num)
	assert.Equal(t, "Do something else", r.Description)
	assert.Equal(t, PASS, r.CurrentMode)
	assert.Equal(t, NIL, r.DesiredMode)
}

func TestRule_String(t *testing.T) {
	r := NewRule("4.1", "Do something", "PASS")
	exp := "4.1  | Cur:PASS  || Do something"
	assert.Equal(t, exp, r.String())
}

func TestModeToInt(t *testing.T) {
	exp := ModeMap
	for m,i := range exp {
		assert.Equal(t, i, ModeToInt(m))
	}
}

func TestModeToStr(t *testing.T) {
	exp := ModeMap
	for m,i := range exp {
		assert.Equal(t, m, ModeToStr(i))
	}
}

func TestRule_Skip(t *testing.T) {
	r := NewRule("4.1", "Do something", "PASS")
	c := map[string]string{"modes-show": "PASS"}
	assert.False(t, r.Skip(c), "Should be shown as PASS is include")
	c["modes-show"] = "WARM"
	assert.True(t, r.Skip(c), "Should be skipped as PASS is not include")
}
