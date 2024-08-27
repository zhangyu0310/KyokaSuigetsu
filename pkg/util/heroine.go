package util

import (
	"errors"
	"sync"
)

type Heroine struct {
	Name           string
	Characteristic string
	LuckyNumber    int
}

var (
	heroineMap = map[string]Heroine{
		"": {
			Name:           "",
			Characteristic: "",
			LuckyNumber:    0,
		},
	}

	// MakeHeroine not 'make' in English, it is「負け」in Japanese.
	makeHeroine []string

	heroineLock sync.Mutex
)

var (
	ErrorNoHeroine = errors.New("no heroine")
)

func init() {
	for _, heroine := range heroineMap {
		makeHeroine = append(makeHeroine, heroine.Name)
	}
}

func MaKeRuNa() (Heroine, error) {
	heroineLock.Lock()
	defer heroineLock.Unlock()
	if len(makeHeroine) == 0 {
		return Heroine{}, ErrorNoHeroine
	}
	heroine := heroineMap[makeHeroine[0]]
	makeHeroine = makeHeroine[1:]
	return heroine, nil
}

func MaKeRu(name string) {
	heroineLock.Lock()
	defer heroineLock.Unlock()
	makeHeroine = append(makeHeroine, name)
}
