package secbench

import (
	"fmt"
	"strings"
	"regexp"
)

const (
	NIL = -1
	NOTE = 0
	INFO = 1
	PASS = 2
	WARN = 3
)

func ModeToStr(m int) string {
	switch m {
	case NOTE:
		return "NOTE"
	case INFO:
		return "INFO"
	case PASS:
		return "PASS"
	case WARN:
		return "WARN"
	}
	return "<nil>"
}

func ModeToInt(m string) int {
	switch m {
	case "NOTE":
		return NOTE
	case "INFO":
		return INFO
	case "PASS":
		return PASS
	case "WARN":
		return WARN
	}
	return NIL
}

type Rule struct {
	Num 		string
	Description string
	DesiredMode int
	CurrentMode int
	Instances	[]Instance
	Pass		bool
}

func NewRule(num, desc, mode string) Rule {
	return Rule{
		Num: num,
		Description: desc,
		CurrentMode: ModeToInt(mode),
		DesiredMode: NIL,
		Instances: []Instance{},
	}
}

func (r *Rule) String() string {
	res := []string{}
	res = append(res, fmt.Sprintf("%-4s | %-5s || %s", r.Num, ModeToStr(r.CurrentMode), r.Description))
	for _, inst := range r.Instances {
		res = append(res, inst.String())
	}
	return strings.Join(res, "\n")

}

func (r *Rule) AddCurrentMode(current string) {
	r.CurrentMode = ModeToInt(current)
}

func (r *Rule) AddDescription(desc string) {
	r.Description = desc
}

func (r Rule) AddInstance(inst Instance) Rule {
	r.Instances = append(r.Instances, inst)
	return r
}

func (r *Rule) Skip(cfg map[string]string) bool {
	modes, ok := cfg["modes-ignore"]
	if !ok {
		return false
	}
	for _, mode := range strings.Split(modes, ",") {
		if mode == ModeToStr(r.CurrentMode) {
			return true
		}
	}
	if cfg["skip-empty-rules"] == "true" && len(r.Instances) == 0 {
		return true
	}
	skipR := strings.Split(cfg["rule-numbers-skip"], ",")
	for _, sR := range skipR {
		if sR == r.Num {
			return true
		}
	}
	skipRre := cfg["rule-numbers-skip-regex"]
	if  skipRre != "" {
		m, err := regexp.Match("^"+skipRre, []byte(r.Num))
		if m && err == nil {
			return true
		}
	}
	notSkip := false
	onlyR := []string{}
	for _, item := range strings.Split(cfg["rule-numbers-only"], ",") {
		if item == "" {
			continue
		}
		onlyR = append(onlyR, item)
	}
	for _, oR := range onlyR {
		if oR == r.Num {
			notSkip = true
		}
	}
	onlyRre := cfg["rule-numbers-only-regex"]
	if  onlyRre != "" {
		m, err := regexp.Match("^"+onlyRre, []byte(r.Num))
		if !m && err == nil {
			return true
		}
	}
	if len(onlyR) >= 1 && !notSkip {
		return true
	}
	return false
}
