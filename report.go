package main

import (
	"fmt"
	"sort"
)

func printReport(pages map[string]int, baseURL string) {
	fmt.Println("=============================")
	fmt.Printf("REPORT for %s\n", baseURL)
	fmt.Println("=============================")

	vec := mapToSlice(pages)

	sort.Slice(vec, func(i, j int) bool {
		// 1. value is different - sort by value (in reverse order)
		if vec[i].value != vec[j].value {
			return vec[i].value > vec[j].value
		}
		// 2. only when value is the same - sort by key
		return vec[i].key < vec[j].key
	})

	for _, v := range vec {
		fmt.Printf("Found %v internal links to %s\n", v.value, v.key)
	}
}

func mapToSlice(in map[string]int) []KV {
	vec := make([]KV, len(in))
	i := 0
	for k, v := range in {
		vec[i].key = k
		vec[i].value = v
		i++
	}
	return vec
}

type KV struct {
	key   string
	value int
}
