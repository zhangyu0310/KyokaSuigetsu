package fakemaster

import "sync"

type Variables struct {
	sync.Mutex
	s map[string]interface{}
}

func NewVariables() *Variables {
	return &Variables{
		s: make(map[string]interface{}),
	}
}

func (v *Variables) SetVariable(key string, val interface{}) {
	v.Lock()
	defer v.Unlock()
	v.s[key] = val
}

func (v *Variables) GetVariable(key string) interface{} {
	v.Lock()
	defer v.Unlock()
	return v.s[key]
}
