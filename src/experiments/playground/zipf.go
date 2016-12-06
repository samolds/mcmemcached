package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"experiments/common"
)

func main() {
	//s := 3.12
	//v := 2.72
	s := 1.1
	v := 5.0
	imax := uint64(50000)

	//its := 1000
	its := 1
	var avg_smallest int
	var avg_biggest int
	for o := 0; o < its; o++ {
		// initialize random number generator with a zipfian distribution
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		zipf := rand.NewZipf(r, s, v, imax)
		if zipf == nil {
			log.Printf("Error creating zipfian thing\n")
			return
		}

		// print out the generated key distribution at the end
		key_distribution := make(map[string]int)

		for i := 0; i < 1000; i++ {
			key_distribution[strconv.Itoa(int(zipf.Uint64()))]++
		}

		smallest := int(imax)
		biggest := -1
		for k, _ := range key_distribution {
			key, err := strconv.Atoi(k)
			if err != nil {
				log.Fatal(err)
			}

			if key > biggest {
				biggest = key
			}

			if key < smallest {
				smallest = key
			}
		}

		avg_smallest += smallest
		avg_biggest += biggest

		log.Printf("Smallest Key: %d, Biggest Key: %d\n", smallest, biggest)
		log.Printf("Key access distribtuion {key access_count}: %v",
			common.OrderByValue(key_distribution))
	}

	asmall := float64(avg_smallest) / float64(its)
	abig := float64(avg_biggest) / float64(its)
	log.Printf("rand.NewZipf(r, %0.2f, %0.2f, %d)", s, v, imax)
	log.Printf("Avg Smallest Key: %0.3f Avg Biggest Key: %0.3f\n\n", asmall, abig)
}
