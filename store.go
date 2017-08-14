package secbench

import (
	"fmt"
	"os"
	"sort"
)

type Store struct {
	Chan <- chan interface{}
	Rules map[string]Rule
}

func NewStore(c chan interface{}) Store {
	return Store{
		Chan: c,
		Rules: map[string]Rule{},
	}
}

func (s *Store) Loop(done chan error) {
	var num string
	var err error
	stop := false
	for {
		if stop {
			break
		}
		val := <- s.Chan
		switch val.(type) {
		case Rule:
			r := val.(Rule)
			//fmt.Println(r.String())
			num = r.Num
			s.Rules[num] = r
		case Instance:
			i := val.(Instance)
			//fmt.Println(i.String())
			s.Rules[num] = s.Rules[num].AddInstance(i)
		default:
			stop = true
		}
	}
	done <- err
}

func (s *Store) Eval(cfg map[string]string) {
	warn_cnt := 0
	nums := []string{}
	for n, _ := range s.Rules {
		nums = append(nums, n)
	}
	sort.Strings(nums)
	for _, num := range nums {
		rule := s.Rules[num]
		if rule.Skip(cfg) {
			continue
		}
		if rule.CurrentMode == WARN {
			warn_cnt++
		}
		fmt.Println(rule.String())
	}
	if cfg["fail-on-warn"] == "true" {
		os.Exit(warn_cnt)
	}
}
