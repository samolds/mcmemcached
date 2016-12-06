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
	chash := consistent.New()

	for _, server := range common.MCACHES {
		chash.Add(server)
	}

	// initialize random number generator with a zipfian distribution
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	zipf := rand.NewZipf(r, 1.1, 5.0, 300000)
	if zipf == nil {
		log.Printf("Error creating zipfian thing\n")
		return
	}

	server_key_count := make(map[string]int)

	for i := 0; i < 100; i++ {
		key := strconv.Itoa(int(zipf.Uint64()))
		server, err := chash.Get(key)
		if err != nil {
			log.Fatal(err)
			return
		}

		server_key_count[server]++

		log.Printf("Key %s, mapped to server: %s", key, server)
	}

	log.Printf("Server key counts: %#v", server_key_count)
}
