package main

import "log"
import "math/rand"
import "time"

func main() {
     n := 100
     i := 0
     
     r := rand.New(rand.NewSource(time.Now().UnixNano()))
     zipf := rand.NewZipf(r, 3.12, 2.72, 5000)

     for i < n {
     	 log.Println(zipf.Uint64())
	 i += 1
     }
}