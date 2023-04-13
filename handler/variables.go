package handler

import "sync"

var (
	GlobalVariables *Variables
	gvMutex         sync.RWMutex
)

func init() {
	GlobalVariables = NewVariables()
}

func GetGlobalVariable(k string) interface{} {
	gvMutex.RLock()
	defer gvMutex.RUnlock()
	return GlobalVariables.GetVariable(k)
}

func SetGlobalVariable(k string, v interface{}) {
	gvMutex.Lock()
	defer gvMutex.Unlock()
	GlobalVariables.SetVariable(k, v)
}

type Variables struct {
	s map[string]interface{}
}

func NewVariables() *Variables {
	return &Variables{
		s: make(map[string]interface{}),
	}
}

func (variables *Variables) SetVariable(k string, v interface{}) {
	variables.s[k] = v
}

func (variables *Variables) GetVariable(k string) interface{} {
	return variables.s[k]
}
