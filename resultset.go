package secbench

import (
	"sync"
	"fmt"
	"sort"
)

type ResultSets struct {
	mu sync.Mutex
	data map[string]Rule

}


func NewResultSet() *ResultSets {
	rs := &ResultSets{
		data: map[string]Rule{},
	}
	return rs
}

func (rs *ResultSets) AddRule(r Rule) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.data[r.Num] = r
}

func (rs *ResultSets) AddInstance(num, mode, msg string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	inst := NewInstance(mode, msg)
	r := rs.data[num]
	r.AddInstance(inst)
	fmt.Println(inst.String())
	rs.data[num] = r
}

func (rs *ResultSets) Eval(cfg map[string]string) {
	nums := []string{}
	for num, _ := range rs.data {
		nums = append(nums, num)
	}
	sort.Strings(nums)
	for _, num := range nums {
		r := rs.data[num]
		if r.Skip(cfg) {
			continue
		}
		fmt.Println(r.String())
	}
}
