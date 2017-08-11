package secbench

import "fmt"

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
	for _, rule := range s.Rules {
		if rule.Skip(cfg) {
			continue
		}
		fmt.Println(rule.String())
	}
}
