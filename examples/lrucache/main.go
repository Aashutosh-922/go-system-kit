package main

import (
	"fmt"

	"github.com/Aashutosh-922/go-system-kit.git/lrucache"
)

func main() {
	cache := lrucache.New[string, int](2)

	cache.Put("session:a", 100)
	cache.Put("session:b", 200)

	// Touch "session:a" so it becomes most recently used.
	_, _ = cache.Get("session:a")

	// This insert evicts "session:b" because capacity is 2.
	cache.Put("session:c", 300)

	fmt.Println("cache length:", cache.Len())

	if v, ok := cache.Get("session:a"); ok {
		fmt.Println("session:a =", v)
	}
	if _, ok := cache.Get("session:b"); !ok {
		fmt.Println("session:b evicted (LRU)")
	}
	if v, ok := cache.Get("session:c"); ok {
		fmt.Println("session:c =", v)
	}
}
