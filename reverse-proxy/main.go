package main

import (
	"ReverseProxy/config"
	"fmt"
	"log"
	"net/http"
)

func main() {
	//load backends
	pool, err := config.LoadConfig("reverse-proxy/config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	fmt.Printf("Loaded %d backends:\n", len(pool.Backends))

	proxyHandler := &ProxyHandler{
		ServerPool: pool,
	}

	fmt.Println("starting proxy server on 8080: ")
	http.ListenAndServe(":8080", proxyHandler)

}
