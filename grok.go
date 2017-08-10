package secbench

import (
	"reflect"
	"github.com/vjeantet/grok"
	"github.com/pkg/errors"

)

var (
	pattern = `^\[%{LEVEL}\]\s+(%{RULE}|%{RESULT})`
)

type Grok struct {
	grok *grok.Grok
}

func NewGrok() Grok {
	gr, _ := grok.New()
	gr, _ = grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	// Parse message
	gr.AddPattern("LEVEL", `(\w+)`)
	gr.AddPattern("RULE_INT", `(\d+)`)
	gr.AddPattern("RULE_FLOAT", `(\d+)?(\.\d+)`)
	gr.AddPattern("RAW_MSG", `.*`)
	gr.AddPattern("RULE1", `\[%{LEVEL:user}\]\s+%{RULE_INT:num}\s+\-\s+%{RAW_MSG:rule}`)
	gr.AddPattern("RULE2", `\[%{LEVEL:user}\]\s+%{RULE_FLOAT:num}\s+\-\s+%{RAW_MSG:rule}`)
	gr.AddPattern("MSG", `\[%{LEVEL:level}\]\s+\*\s+%{RAW_MSG:msg}`)
	g := Grok{
		grok: gr,
	}
	return g
}

func (g *Grok) ParseMsg(str string) (map[string]string,  error) {
	res, err := g.grok.Parse("%{MSG}", str)
	if err != nil {
		return nil, err
	}
	keys := reflect.ValueOf(res).MapKeys()
	if len(keys) == 0 {
		err = errors.New("Empty result")
	}
	return res, err
}
func (g *Grok) ParseRule1(str string) (map[string]string,  error) {
	res, err := g.grok.Parse("%{RULE1}", str)
	if err != nil {
		return nil, err
	}
	keys := reflect.ValueOf(res).MapKeys()
	if len(keys) == 0 {
		err = errors.New("Empty result")
	}
	return res, err
}
func (g *Grok) ParseRule2(str string) (map[string]string,  error) {
	res, err := g.grok.Parse("%{RULE2}", str)
	if err != nil {
		return nil, err
	}
	keys := reflect.ValueOf(res).MapKeys()
	if len(keys) == 0 {
		err = errors.New("Empty result")
	}
	return res, err
}

func (g *Grok) ParseRules(str string) (map[string]string,  error) {
	msg, err := g.ParseRule1(str)
	if err != nil {
		return g.ParseRule2(str)
	}
	return msg, err

}
func (g *Grok) ParseLine(str string) (map[string]string,  error) {
	msg, err := g.ParseMsg(str)
	if err != nil {
		return g.ParseRules(str)
	}
	return msg, err

}