package entity

import "fmt"

type GroupQueryFrequency map[string]int64

func (g GroupQueryFrequency) String() string {
	var s string

	for query, freq := range g {
		s += fmt.Sprintf("%s %d\n", query, freq)
	}

	return s
}

func (g GroupQueryFrequency) BytesEncode() []byte {
	return []byte(g.String())
}

func (g GroupQueryFrequency) Clear() {
	for key := range g {
		delete(g, key)
	}
}
func (g GroupQueryFrequency) Add(addFreq GroupQueryFrequency, max int) {
	for k, v := range addFreq {
		if val, ok := g[k]; ok {
			g[k] = val + v
			delete(addFreq, k)
		}
	}

	if len(g) == max {
		return
	}

	for k, v := range addFreq {
		if len(g) < max {
			g[k] = v
			delete(addFreq, k)
		}
	}
}
