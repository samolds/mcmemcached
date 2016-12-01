package main

import (
    "github.com/bradfitz/gomemcache/memcache"
)
import "log"

func main() {
    mc := memcache.New("localhost:11211", "localhost:11212")

    log.Printf("Sending foo=bar")

    mc.Set(&memcache.Item{Key: "foo", Value: []byte("bar")})

    it, _ := mc.Get("foo")

    log.Printf("Got %s for key foo", it.foo)
}