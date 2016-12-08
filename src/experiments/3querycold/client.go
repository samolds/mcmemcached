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
	var old_mc *memcache.Client
	old_valid := false
	var active_mc *memcache.Client
	//old_mc = config1
	active_mc = config1

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

	// threshold of cache misses before adding 4th server
	thresh := common.ADD_SERVER_THRESHOLD

	// simulate n cache requests
	n := common.ITERATION_COUNT

	// give the servers an 8th of the iterations to "warm up"
	warm_up_its := n / common.WARM_UP_RATIO
	for i := 0; i < n; i++ {
		key := strconv.Itoa(int(zipf.Uint64())) // convert to string
		key_distribution[key]++

		// try and get the key from the cache
		_, err := active_mc.Get(key) // returns item, err
		if err == memcache.ErrCacheMiss {
			miss := true
			if old_valid {
				// Delay for consulting second memcached server
				common.AddDelayPoint(&stats, fetch_delay)

				_, er := old_mc.Get(key)
				if er != memcache.ErrCacheMiss {
					miss = false
					active_mc.Set(&memcache.Item{Key: key, Value: memcache_value})

					// Delay for the first memcache server
					// since it won't be captured in a miss
					common.AddDelayPoint(&stats, fetch_delay)
				}
			}

			if miss {
				// cache miss, so add the key/value to the cache
				common.AddDelayPoint(&stats, database_delay)

				cache_misses++
				active_mc.Set(&memcache.Item{Key: key, Value: memcache_value})
			}
		} else {
			common.AddDelayPoint(&stats, fetch_delay)
		}

		// after a fraction of cache requests, to give servers time to "warm up",
		// if there have been more than 35% cache misses for the requests thus far,
		// "spin up new server" (switch to configuration two with 4 cache servers)
		if i > warm_up_its && (cache_misses*100)/i >= thresh && active_mc == config1 {
			log.Printf("\tAdded new server!! cache misses: %d, requests sent: %d\n",
				cache_misses, i)
			active_mc = config2
		}

		// for printing out csv data of cache misses for graphing ratio
		common.WriteCacheMissRatio(cache_misses, i)
	}

	common.WriteTimeStats(&stats)

	common.LogResults(active_mc, key_distribution, cache_misses, n, memcache_value)
}
