package spider

import "sync"

func NewStringSet() *StringSet {
	s := &StringSet{}
	s.set = map[string]bool{}
	return s
}

type StringSet struct {
	set map[string]bool
	mu  sync.RWMutex
}

func (v *StringSet) Add(val string) {
	v.mu.Lock()
	v.set[val] = true
	v.mu.Unlock()
}

func (v *StringSet) IsExists(val string) bool {
	v.mu.RLock()
	_, ok := v.set[val]
	v.mu.RUnlock()
	return ok
}

func (v StringSet) All() []string {
	v.mu.RLock()
	vals := []string{}
	for val, _ := range v.set {
		vals = append(vals, val)
	}
	v.mu.RUnlock()
	return vals
}
