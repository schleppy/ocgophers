package main

import (
	"fmt"
	"time"

	"github.com/hashicorp/golang-lru"
)

func main() {
	cache, _ := lru.New(5)
	for i, key := 0, 0; i < 100; i, key = i+1, i%5 {
		if res, ok := cache.Get(key); ok {
			fmt.Printf("Got item %d from cache\n", res)
			continue
		}
		item := getSlowThing(key)
		fmt.Printf("Adding %d to cache\n", item)
		cache.Add(key, item)

	}
	fmt.Printf("Cache size: %d\n", cache.Len())
	time.Sleep(1 * time.Second)
}

func getSlowThing(i int) int {
	fmt.Printf("\tRetrieving %d\n", i)
	time.Sleep(time.Duration(i*250) * time.Millisecond)
	return i * 100
}
