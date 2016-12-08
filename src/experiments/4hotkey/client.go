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
	//var old_mc *memcache.Client
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

	// simulate n cache requests
	n := common.ITERATION_COUNT
	warm_up_its := n / 8 // give the servers an 8th of the iterations to "warm up"
	for i := 0; i < n; i++ {
		key := strconv.Itoa(int(zipf.Uint64())) // convert to string
		key_distribution[key]++

		// try and get the key from the cache
		_, err := active_mc.Get(key) // returns item, err
		if err == memcache.ErrCacheMiss {
			common.AddDelayPoint(&stats, database_delay)

			// cache miss, so add the key/value to the cache
			cache_misses++
			active_mc.Set(&memcache.Item{Key: key, Value: memcache_value})
		} else {
			common.AddDelayPoint(&stats, fetch_delay)
		}

		// after a fraction of cache requests, to give servers time to "warm up",
		// if there have been more than 35% cache misses for the requests thus far,
		// "spin up new server" (switch to configuration two with 4 cache servers)
		if i > warm_up_its && (cache_misses*100)/i >= 25 && active_mc == config1 {
			log.Printf("\tAdded new server!! cache misses: %d, requests sent: %d\n",
				cache_misses, i)
			active_mc = config2
			err = catchUpColdServer(config1, config2, key_distribution)
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		// for printing out csv data of cache misses for graphing ratio
		common.WriteCacheMissRatio(cache_misses, i)
	}

	common.WriteTimeStats(&stats)

	hot_key_servers, err := common.GetHotKeysPerServer(active_mc,
		key_distribution)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Got %d cache misses for %d requests", cache_misses, n)
	for _, hot_key_server := range hot_key_servers {
		log.Printf("%s\n", hot_key_server.String(5))
	}
}

// TODO(sam): add optimizations here. additional memcache server will be cold
// and will need to be caught up with existing key/value pairs all at once
// right here
//
// once we hit the threshold and moved from using 3 servers to 4, we'd need to
// figure out which keys have been remapped to the new server. and once we've
// figured out the hottest ~20 keys or something that have been remapped, we
// just set those on the new server
func catchUpColdServer(old_mc *memcache.Client, new_mc *memcache.Client,
	key_distribution map[string]int) error {

	// get all of the server objects from the memcache client
	servers, err := new_mc.Servers.Servers()
	if err != nil {
		return err
	}

	// we know the cold server is going to be the 4th on in the new server list
	new_server_addr := servers[3].String()

	// keys that are now on the new cold server
	// using map[string]struct{} because it's easier to check if key is in a map
	// rather than a slice, and the struct{} provides no additional memory
	// requirements
	remapped_keys := make(map[string]struct{})

	// find all of the keys that have been remapped to the new server
	for key, _ := range key_distribution {
		server, err := new_mc.Servers.PickServer(key)
		if err != nil {
			return err
		}

		if server.String() == new_server_addr {
			var empty struct{}
			remapped_keys[key] = empty
		}
	}

	// get the hottest keys on each of the original server
	hot_key_servers, err := common.GetHotKeysPerServer(old_mc,
		key_distribution)
	if err != nil {
		return err
	}

	// read in a file that contains the value that will be set for every memcache
	// key/value pair
	memcache_value, err := ioutil.ReadFile(common.MEMCACHE_VALUE_FILENAME)
	if err != nil {
		return err
	}

	// function to move a key from one server to another if it exists and return
	// "how hot" that key was relative to the other keys it knows about
	remap_key_from_old_to_new := func(server_index int, index int) (float64, bool) {
		if len(hot_key_servers[server_index].KeyDistribution) <= index {
			return 0.0, true
		}

		key_dist := hot_key_servers[server_index].KeyDistribution[index]
		_, key_needs_remapping := remapped_keys[key_dist.Key]
		if key_needs_remapping {
			// TODO(sam): add sleep or something to indicate that server needs to go
			// get the value from the old server
			new_mc.Set(&memcache.Item{Key: key_dist.Key, Value: memcache_value})
			percentage := float64(key_dist.Value) / float64(
				hot_key_servers[server_index].TotalKeyHitCount)
			return percentage, false
		}
		return 0.0, false
	}

	total_key_hit_percentage := 0.0
	index := 0

	// iterate through all of the keys on each existing server from hottest to
	// coldest until the new server has 20% of the hottest keys from the old
	// servers
	for total_key_hit_percentage < 20.0 {
		key_hit_percentage, too_much := remap_key_from_old_to_new(0, index)
		if too_much {
			break
		}
		total_key_hit_percentage += key_hit_percentage

		key_hit_percentage, too_much = remap_key_from_old_to_new(1, index)
		if too_much {
			break
		}
		total_key_hit_percentage += key_hit_percentage

		key_hit_percentage, too_much = remap_key_from_old_to_new(2, index)
		if too_much {
			break
		}
		total_key_hit_percentage += key_hit_percentage

		index++
	}

	return nil
}
