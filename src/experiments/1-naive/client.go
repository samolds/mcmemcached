package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"experiments/common"

	"github.com/rainycape/memcache"
)

func main() {
	// cold cache = not initialized
	config1, err := memcache.New(common.MCACHES[0], common.MCACHES[1],
		common.MCACHES[2])
	if err != nil {
		log.Fatal(err)
		return
	}

	config2, err := memcache.New(common.MCACHES[0], common.MCACHES[1],
		common.MCACHES[2], common.MCACHES[3])
	if err != nil {
		log.Fatal(err)
		return
	}

	// the configurations being used
	//var old_mc *memcache.Client
	var active_mc *memcache.Client
	//old_mc = config1
	active_mc = config1

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
	warm_up_its := n / 8 // give the servers an 8th of the iterations to "warm up"
	for i := 0; i < n; i++ {
		key := strconv.Itoa(int(zipf.Uint64())) // convert to string
		key_distribution[key]++

		// try and get the key from the cache
		//item, err := active_mc.Get(key) // returns item, err
		_, err := active_mc.Get(key) // returns item, err
		if err == memcache.ErrCacheMiss {
			// if active_mc == config2 {
			// TODO(sam): add optimizations here. when there is a cache miss after a
			// new server has been added and it is within a threshold of being new,
			// refer back to the old memcache server instead of "going to the db"
			// } else {
			// sleep here to simulate call to db to get value to add to cache
			time.Sleep(time.Duration(database_delay) * time.Millisecond)
			//}

			// cache miss, so add the key/value to the cache
			cache_misses++
			active_mc.Set(&memcache.Item{Key: key, Value: []byte("fake value")})

			//log.Printf("Using key: '%s', cache miss! adding to cache", key)
		} else {
			//log.Printf("\tUsing key: '%s', cache hit! value: '%#v'", key,
			//	string(item.Value))
		}

		// after a fraction of cache requests, to give servers time to "warm up",
		// if there have been more than 35% cache misses for the requests thus far,
		// "spin up new server" (switch to configuration two with 4 cache servers)
		if i > warm_up_its && (cache_misses*100)/i >= 35 && active_mc == config1 {
			log.Printf("\tAdded new server!! cache misses: %d, requests sent: %d\n",
				cache_misses, i)
			active_mc = config2
			// TODO(sam): add optimizations here. additional memcache server will be
			// cold and will need to be caught up with existing key/value pairs
		}
	}

	log.Printf("Key access distribtuion {key access_count}: %v",
		common.OrderByValue(key_distribution))
	log.Printf("Got %d cache misses for %d requests", cache_misses, n)
}
