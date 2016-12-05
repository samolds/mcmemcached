package common

import (
	"fmt"
	"sort"

	"github.com/rainycape/memcache"
)

var (
	MCACHES = []string{
		"localhost:11211",
		"localhost:11212",
		"localhost:11213",
		"localhost:11214",
	}
)

// time stored in milliseconds
type TimeStats struct {
	RunTime float32
	Means   []TimePair
}

type TimePair struct {
	MeanValue float32
	AtTime    float32
}

// function to keep timing statistics
// below are the time values from the HotCloud paper:
// key-value fetch, 0.3 ms
// goes to DB, 8 ms
func AddDelayPoint(stats *TimeStats, delay float32) {
	stats.RunTime = stats.RunTime + delay
	pair := TimePair{
		MeanValue: stats.RunTime / float32(len(stats.Means)+1),
		AtTime:    stats.RunTime,
	}
	stats.Means = append(stats.Means, pair)
}

func WriteTimeStats(stats *TimeStats) {
	for i, m := range stats.Means {
		fmt.Printf("[%v, %v]", m.MeanValue, m.AtTime)
		if i < len(stats.Means)-1 {
			fmt.Printf(",")
		}
	}
	fmt.Printf("\n")
}

// performs a single get for each key in the key_distribution map to
// figure out which key is active on each server
// returns a map of server address to keys active on each server
func RevealKeyOwners(mc *memcache.Client, key_distributions map[string]int) (
	map[string][]string, error) {

	key_owners := make(map[string][]string)

	// loop through all keys from the map of key distributions
	for key, _ := range key_distributions {

		// perform a get to see if the key even exists on the memcache servers
		_, err := mc.Get(key)
		if err == nil {
			// figure out which server the key would be put on
			server, err := mc.Servers.PickServer(key)
			if err != nil {
				return key_owners, err
			}
			server_str := server.String()

			key_owners[server_str] = append(key_owners[server_str], key)
		}
	}

	return key_owners, nil
}

// used for ordering a map[string]int by the values
// found here: http://stackoverflow.com/a/18695740
func OrderByValue(stringIntMap map[string]int) PairList {
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
