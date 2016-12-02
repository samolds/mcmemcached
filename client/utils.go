package main

import (
	"sort"
)

// used for ordering a map[string]int by the values
// found here: http://stackoverflow.com/a/18695740
func orderByValue(stringIntMap map[string]int) PairList {
	pl := make(PairList, len(stringIntMap))
	i := 0
	for k, v := range stringIntMap {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
