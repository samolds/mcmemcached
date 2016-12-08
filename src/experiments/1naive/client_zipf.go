package main

import (
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
	zipf := rand.NewZipf(r, common.ZIPF_S, common.ZIPF_V, common.ZIPF_IMAX)

	// print out the generated key distribution at the end
	key_distribution := make(map[string]int)

	// keep track of the number of cache misses
	cache_misses := 0

	// fake memcached fetch delay
	var fetch_delay float32 = common.FETCH_DELAY

	// fake database delay in milliseconds
	var database_delay float32 = common.DATABASE_DELAY

	var stats common.TimeStats

	// read in a file that contains the value that will be set for every memcache
	// key/value pair
	memcache_value, err := ioutil.ReadFile(common.MEMCACHE_VALUE_FILENAME)
	if err != nil {
		log.Fatal(err)
		return
	}

	// for printing out csv data of cache misses for graphing
	common.WriteTimeStatsHeader()
	common.WriteCacheMissRatioHeader()

	// simulate n cache requests
	n := common.ITERATION_COUNT
	for i := 0; i < n; i++ {
		key := strconv.Itoa(int(zipf.Uint64())) // convert to string
		key_distribution[key]++

		// try and get the key from the cache
		_, err := mc.Get(key) // returns item, err
		if err == memcache.ErrCacheMiss {
			common.AddDelayPoint(&stats, database_delay)

			// cache miss, so add the key/value to the cache
			cache_misses++
			mc.Set(&memcache.Item{Key: key, Value: memcache_value})
		} else {
			common.AddDelayPoint(&stats, fetch_delay)
		}

		// for printing out csv data of cache misses for graphing ratio
		common.WriteCacheMissRatio(cache_misses, i)
	}

	common.WriteTimeStats(&stats)

	common.LogResults(mc, key_distribution, cache_misses, n, memcache_value)
}
