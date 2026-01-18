package main

import (
	"ReverseProxy/config"
	"ReverseProxy/healthCheker"
	"fmt"
	"log"
	"net/http"
	"time"
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

	//test the health cheker function 
	cheker.StartHealthCheck(pool,5*time.Second)

	fmt.Println("starting proxy server on 8080: ")
	http.ListenAndServe(":8080", proxyHandler)


}
