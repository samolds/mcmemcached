package main

import (
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
	zipf := rand.NewZipf(r, 1.1, 5.0, 50000)

	// print out the generated key distribution at the end
	key_distribution := make(map[string]int)

	// keep track of the number of cache misses
	cache_misses := 0

	// fake database delay in milliseconds
	database_delay := 1 // want 8, but i don't like waiting so i set it to 1 for now

	// simulate n cache requests
	n := 100000
	for i := 0; i < n; i++ {
		key := strconv.Itoa(int(zipf.Uint64())) // convert to string
		key_distribution[key]++

		// try and get the key from the cache
		//item, err := mc.Get(key) // returns item, err
		_, err := mc.Get(key) // returns item, err
		if err == memcache.ErrCacheMiss {
			time.Sleep(time.Duration(database_delay) * time.Millisecond)

			// cache miss, so add the key/value to the cache
			cache_misses++
			mc.Set(&memcache.Item{Key: key, Value: []byte("fake value")})

			//log.Printf("Using key: '%s', cache miss! adding to cache", key)
		} else {
			//log.Printf("\tUsing key: '%s', cache hit! value: '%#v'", key,
			//	string(item.Value))
		}
	}

	log.Printf("Key access distribtuion {key access_count}: %v",
		common.OrderByValue(key_distribution))
	log.Printf("Got %d cache misses for %d requests", cache_misses, n)
}
