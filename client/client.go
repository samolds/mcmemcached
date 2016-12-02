package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	// cold cache = not initialized
	mc := memcache.New("localhost:11211", "localhost:11212")

	// initialize random number generator with a zipfian distribution
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	zipf := rand.NewZipf(r, 3.12, 2.72, 5000)

	// keep track of the number of cache misses
	cache_misses := 0

	// fake database delay in milliseconds
	database_delay := 1

	// simulate n cache requests
	n := 100
	for i := 0; i < n; i++ {
		key := strconv.Itoa(int(zipf.Uint64())) // convert to string

		// try and get the key from the cache
		item, err := mc.Get(key) // returns item, err
		if err == memcache.ErrCacheMiss {
			// sleep here to simulate call to db to get value to add to cache
			time.Sleep(time.Duration(database_delay) * time.Millisecond)

			// cache miss, so add the key/value to the cache
			cache_misses++
			mc.Set(&memcache.Item{Key: key, Value: []byte("fake value")})
			log.Printf("Using key: '%s', cache miss! adding to cache", key)
		} else {
			log.Printf("\tUsing key: '%s', cache hit! value: '%#v'", key,
				string(item.Value))
		}
	}

	log.Printf("Got %d cache misses", cache_misses)
}