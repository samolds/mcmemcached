package main

import (
	"log"
	"strconv"

	"experiments/common"

	"github.com/rainycape/memcache"
)

// simply iterates from 0 to n to generate workload (requests)
// results in a miss every time

func main() {
	// just use 3 memcache servers total
	mc, err := memcache.New(common.MCACHES[0], common.MCACHES[1],
		common.MCACHES[2])
	if err != nil {
		log.Fatal(err)
		return
	}

	// print out the generated key distribution at the end
	key_distribution := make(map[string]int)

	// keep track of the number of cache misses
	cache_misses := 0

	// fake memcached fetch delay
	var fetch_delay float32 = 0.3

	// fake database delay in milliseconds
	var database_delay float32 = 8.0

	var stats common.TimeStats

	// simulate n cache requests
	n := 100000
	for i := 0; i < n; i++ {
		key := strconv.Itoa(i) // convert to string
		key_distribution[key]++

		// try and get the key from the cache
		//item, err := mc.Get(key) // returns item, err
		_, err := mc.Get(key) // returns item, err
		if err == memcache.ErrCacheMiss {
			common.AddDelayPoint(&stats, database_delay)

			// cache miss, so add the key/value to the cache
			cache_misses++
			mc.Set(&memcache.Item{Key: key, Value: []byte("fake value")})

			//log.Printf("Using key: '%s', cache miss! adding to cache", key)
		} else {
			common.AddDelayPoint(&stats, fetch_delay)
			//log.Printf("\tUsing key: '%s', cache hit! value: '%#v'", key,
			//	string(item.Value))
		}
	}

	log.Printf("Key access distribtuion {key access_count}: %v",
		common.OrderByValue(key_distribution))
	log.Printf("Got %d cache misses for %d requests", cache_misses, n)
	common.WriteTimeStats(&stats)
}
