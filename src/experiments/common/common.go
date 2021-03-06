package common

import (
	"errors"
	"fmt"
	"log"
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

	FETCH_DELAY             float32 = 0.3
	DATABASE_DELAY          float32 = 8.0
	MEMCACHE_VALUE_FILENAME string  = "data/memcache_value.txt"
	ITERATION_COUNT         int     = 475000

	ZIPF_S    float64 = 1.01
	ZIPF_V    float64 = 20.1
	ZIPF_IMAX uint64  = 200000

	PRINT_CACHE_MISS_RATIO bool = false
	PRINT_TIME_STATS       bool = true

	STATS_WINDOW_SIZE float32 = 1000.0

	WARM_UP_RATIO        int = 2  // 1/WARM_UP_RATIO for time to allow initial warm up
	ADD_SERVER_THRESHOLD int = 10 // 10% cache miss before adding new server
)

// time stored in milliseconds
type TimeStats struct {
	RunTime     float32
	Means       []TimePair
	WindowSum   float32
	WindowCount float32
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
	stats.WindowSum = stats.WindowSum + delay
	stats.WindowCount = stats.WindowCount + 1
	stats.RunTime = stats.RunTime + delay

	if stats.WindowCount == STATS_WINDOW_SIZE {
		pair := TimePair{
			MeanValue: stats.WindowSum / stats.WindowCount,
			AtTime:    stats.RunTime,
		}
		stats.Means = append(stats.Means, pair)

		stats.WindowSum = 0
		stats.WindowCount = 0
	}
}

func WriteTimeStatsHeader() {
	if PRINT_TIME_STATS {
		fmt.Printf("mean_response_time,time\n")
	}
}

func WriteTimeStats(stats *TimeStats) {
	if PRINT_TIME_STATS {
		for _, m := range stats.Means {
			fmt.Printf("%v,%v\n", m.MeanValue, m.AtTime)
		}
	}
}

func WriteCacheMissRatioHeader() (int, error) {
	if PRINT_CACHE_MISS_RATIO {
		return fmt.Printf("iteration,cache_miss_ratio\n")
	}
	return 0, nil
}

func WriteCacheMissRatio(cache_misses int, i int) (int, error) {
	if PRINT_CACHE_MISS_RATIO {
		return fmt.Printf("%d,%0.3f\n", i+1, float64(cache_misses)/float64(i+1))
	}
	return 0, nil
}

func LogResults(mc *memcache.Client, key_distribution map[string]int,
	cache_misses int, n int, memcache_value []byte) {
	hot_key_servers, err := GetHotKeysPerServer(mc, key_distribution)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Generated %d unique keys and %0.2fMB of data",
		len(key_distribution),
		float32(len(key_distribution))*float32(len(memcache_value))/1000000.0)

	log.Printf("Got %d cache misses for %d requests", cache_misses, n)
	for _, hot_key_server := range hot_key_servers {
		log.Printf("%s\n", hot_key_server.String(5))
	}
}

type HotKeysPerServer struct {
	Server           string   // address of the server
	Keys             []string // list of all keys active on server
	KeyDistribution  PairList // ordered list of {key, hit_count} for all active keys
	TotalKeyHitCount int      // the total aggregate hit_count for all keys on server
}

func (hk *HotKeysPerServer) String(num_top_keys int) string {
	top_keys := make(map[string]float64)
	for _, pair := range hk.KeyDistribution[:num_top_keys] {
		top_keys[pair.Key] = float64(pair.Value) / float64(hk.TotalKeyHitCount)
	}

	return fmt.Sprintf("server: %s, num active keys: %d, %d hottest keys: %v",
		hk.Server, len(hk.Keys), num_top_keys, top_keys)
}

// using the map of server to active keys and the map of all keys to counts,
// builds up a list of HotKeys on each server
func GetHotKeysPerServer(mc *memcache.Client,
	key_distributions map[string]int) ([]HotKeysPerServer, error) {

	var hot_keys []HotKeysPerServer
	key_owners, err := RevealKeyOwners(mc, key_distributions)
	if err != nil {
		return hot_keys, err
	}

	// iterate through each of the servers that have active keys
	for server, _ := range key_owners {
		distribution := make(map[string]int)
		active_server_keys, ok := key_owners[server]
		if !ok {
			return nil, errors.New("key didn't exist in key_owners")
		}

		total_key_hit_count := 0

		// get the key distributions for just the keys active on each server
		// we want to filter out the keys on other servers to just get the
		// distribution of active keys on each individual server
		for _, key := range active_server_keys {
			count, ok := key_distributions[key]
			if !ok {
				return nil, errors.New("key didn't exist in key_distributions")
			}

			total_key_hit_count += count
			distribution[key] = count
		}

		// build new struct that the server name, the active keys it has, and the
		// distribution for each of those keys
		hot_keys = append(hot_keys, HotKeysPerServer{
			Server:           server,
			Keys:             active_server_keys,
			KeyDistribution:  OrderByValue(distribution),
			TotalKeyHitCount: total_key_hit_count,
		})
	}

	return hot_keys, nil
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
		// no error means the key exists on the server
		if err == nil {
			// figure out which server the key would be put on
			server, err := mc.Servers.PickServer(key)
			if err != nil {
				return key_owners, err
			}
			server_str := server.String()

			// adds the key into the map for the server it belonged to
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
