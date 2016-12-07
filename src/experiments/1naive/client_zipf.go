package main

import (
	//"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"

	"experiments/common"

	"github.com/rainycape/memcache"
)

// uses a zipfian distribution to generate workload (requests)

func main() {
	// just use 3 memcache servers total
	mc, err := memcache.New(common.MCACHES[0], common.MCACHES[1],
		common.MCACHES[2])
	if err != nil {
		log.Fatal(err)
		return
	}

	// initialize random number generator with a zipfian distribution
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	zipf := rand.NewZipf(r, 1.1, 5.0, 300000)

	// print out the generated key distribution at the end
	key_distribution := make(map[string]int)

	// keep track of the number of cache misses
	cache_misses := 0

	// fake memcached fetch delay
	var fetch_delay float32 = 0.3

	// fake database delay in milliseconds
	var database_delay float32 = 8.0

	var stats common.TimeStats

	// read in a file that contains the value that will be set for every memcache
	// key/value pair
	memcache_value, err := ioutil.ReadFile("data/memcache_value.txt")
	if err != nil {
		log.Fatal(err)
		return
	}

	//fmt.Printf("iteration,cache_miss_ratio\n")
	// simulate n cache requests
	n := 1000000
	for i := 0; i < n; i++ {
		key := strconv.Itoa(int(zipf.Uint64())) // convert to string
		key_distribution[key]++

		// try and get the key from the cache
		//item, err := mc.Get(key) // returns item, err
		_, err := mc.Get(key) // returns item, err
		if err == memcache.ErrCacheMiss {
			common.AddDelayPoint(&stats, database_delay)

			// cache miss, so add the key/value to the cache
			cache_misses++
			mc.Set(&memcache.Item{Key: key, Value: memcache_value})

			//log.Printf("Using key: '%s', cache miss! adding to cache", key)
		} else {
			common.AddDelayPoint(&stats, fetch_delay)
			//log.Printf("\tUsing key: '%s', cache hit! value: '%#v'", key,
			//	string(item.Value))
		}

		// ratio:
		//fmt.Printf("%d,%0.3f\n", i+1, float64(cache_misses)/float64(i+1))

		// non-ratio:
		//fmt.Printf("%d,%d\n", i+1, cache_misses)
	}

	//log.Printf("Key access distribtuion {key access_count}: %v",
	//	common.OrderByValue(key_distribution))
	log.Printf("Got %d cache misses for %d requests", cache_misses, n)
	//common.WriteTimeStats(&stats)

	hot_key_servers, err := common.GetHotKeysPerServer(mc,
		key_distribution)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, hot_key_server := range hot_key_servers {
		log.Printf("%s\n", hot_key_server.String(5))
	}
}
