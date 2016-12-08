package main

import (
	"io/ioutil"
	"log"

	"experiments/common"
)

func main() {
	memcache_value, err := ioutil.ReadFile(common.MEMCACHE_VALUE_FILENAME)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("the file is %d bytes", len(memcache_value))
}
