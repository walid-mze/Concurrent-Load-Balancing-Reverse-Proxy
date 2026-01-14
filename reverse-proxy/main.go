package main

import (
	"fmt"
	"log"

	"ReverseProxy/config"
)

func main() {
	pool, err := config.LoadConfig("reverse-proxy/config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	fmt.Printf("Loaded %d backends:\n", len(pool.Backends))
	for i, backend := range pool.Backends {
		fmt.Printf("  %d: %s (Alive: %v)\n", i+1, backend.URL, backend.Alive)
	}
}
