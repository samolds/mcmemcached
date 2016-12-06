package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"experiments/common"
	"github.com/consistent"
)

func main() {
	chash_config1 := consistent.New()
	chash_config1.Add(common.MCACHES[0])
	chash_config1.Add(common.MCACHES[1])
	chash_config1.Add(common.MCACHES[2])

	chash_config2 := consistent.New()
	chash_config2.Add(common.MCACHES[0])
	chash_config2.Add(common.MCACHES[1])
	chash_config2.Add(common.MCACHES[2])
	chash_config2.Add(common.MCACHES[3])

	// initialize random number generator with a zipfian distribution
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	zipf := rand.NewZipf(r, 1.1, 5.0, 300000)
	if zipf == nil {
		log.Printf("Error creating zipfian thing\n")
		return
	}

	server_key_count1 := make(map[string]int)
	server_key_count2 := make(map[string]int)

	n := 10

	for i := 0; i < n; i++ {
		key := strconv.Itoa(int(zipf.Uint64()))

		server1, err := chash_config1.Get(key)
		if err != nil {
			log.Fatal(err)
			return
		}
		server_key_count1[server1]++
		log.Printf("Key %s, mapped to server1: %s", key, server1)

		server2, err := chash_config2.Get(key)
		if err != nil {
			log.Fatal(err)
			return
		}
		server_key_count2[server2]++
		log.Printf("Key %s, mapped to server2: %s", key, server2)

		if server1 != server2 {
			log.Printf("\t\tKey %s, on different servers: %s, %s", key, server1, server2)
		}
	}

	log.Printf("Server 1 key counts: %#v", server_key_count1)
	log.Printf("Server 2 key counts: %#v", server_key_count2)
}
