package secbench

import (
	"fmt"
	"strings"
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
	Instances	[]Result
	Pass		bool
}

func NewRule(num, desc, mode string) Rule {
	return Rule{
		Num: num,
		Description: desc,
		CurrentMode: ModeToInt(mode),
		DesiredMode: NIL,
		Instances: []Result{},
	}
}

func (r *Rule) String() string {
	return fmt.Sprintf("%-4s | Cur:%-5s || %s", r.Num, ModeToStr(r.CurrentMode), r.Description)

}

func (r *Rule) AddCurrentMode(current string) {
	r.CurrentMode = ModeToInt(current)
}

func (r *Rule) AddDescription(desc string) {
	r.Description = desc
}

func (r Rule) AddResult(res Result) {
	r.Instances = append(r.Instances, res)
}

func (r *Rule) Skip(cfg map[string]string) bool {
	modes, ok := cfg["modes-show"]
	if !ok {
		return false
	}
	for _, mode := range strings.Split(modes, ",") {
		if mode == ModeToStr(r.CurrentMode) {
			return false
		}
	}
	return true
}
