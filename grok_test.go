package secbench

import (
	"testing"
	"github.com/stretchr/testify/assert"
)
const (
	r1 = "[INFO] 4 - Container Images and Build File"
	r2 = "[WARN] 4.1  - Ensure a user for the container has been created"
	r3 = "[NOTE] 4.2  - Ensure that containers use trusted base images"
	m1 = "[WARN]      * Running as root: gcollect_agent.2lbtasd6ltpfthys0mpvqng36.q1t9wm4xlykmcvruzme9bncfm"
)

var (
	rules = []string{r1,r2,r3}
	msgs = []string{m1,
		"[WARN]      * Running as root: dviz_client.1.65tgygcncsttortakexp152n8",
		"[WARN]      * Running as root: gocd_agent.4.iqupn0jnkvl44rrz7kcb5wh34",
		"[WARN]      * Running as root: gocd_magent.1.n723u0e7u4kuyy7vcmbac14xb",
		"[WARN]      * Running as root: gocd_server.1.nvcaygx5mtl3niavhq3g7pc74",
	}
	fail = "asdasd"
)

func TestGrok_ParseMsg(t *testing.T) {
	exp := map[string]string{"mode":"WARN", "msg":"Running as root: gcollect_agent.2lbtasd6ltpfthys0mpvqng36.q1t9wm4xlykmcvruzme9bncfm"}
	g := NewGrok()
	got, err := g.parseMsg(m1)
	assert.NoError(t, err, "Should go through")
	assert.Equal(t, exp, got)
}


func TestGrok_ParseRule1(t *testing.T) {
	exp := map[string]string{"mode":"INFO", "num":"4", "rule":"Container Images and Build File"}
	g := NewGrok()
	got, err := g.parseRule1(r1)
	assert.NoError(t, err, "Should go through")
	assert.Equal(t, exp, got)
}

func TestGrok_ParseRule2(t *testing.T) {
	exp := map[string]string{"mode":"WARN", "num":"4.1", "rule":"Ensure a user for the container has been created"}
	g := NewGrok()
	got, err := g.parseRule2(r2)
	assert.NoError(t, err, "Should go through")
	assert.Equal(t, exp, got)
}

func TestGrok_ParseMsgs(t *testing.T) {
	g := NewGrok()
	for _, msg := range msgs {
		res, err := g.parseMsg(msg)
		assert.NoError(t, err, "Should go through")
		_, isMsg := res["msg"]
		_, isRule := res["rule"]
		assert.True(t, isMsg || isRule, "Should be either")
	}
}

func TestGrok_ParseRules(t *testing.T) {
	g := NewGrok()
	for _, rule := range rules {
		res, err := g.parseRules(rule)
		assert.NoError(t, err, "Should go through")
		_, isMsg := res["msg"]
		_, isRule := res["rule"]
		assert.True(t, isMsg || isRule, "Should be either")
	}
}

func TestGrok_ParseLines(t *testing.T) {
	g := NewGrok()
	tests := []string{}
	tests = append(tests, msgs...)
	tests = append(tests, rules...)
	for _, test := range tests {
		res, err := g.ParseLine(test)
		assert.NoError(t, err, "Should go through")
		_, isMsg := res["msg"]
		_, isRule := res["rule"]
		assert.True(t, isMsg || isRule, "Should be either")
	}
	_, err := g.ParseLine(fail)
	assert.Error(t, err, "Should fail through")

}
