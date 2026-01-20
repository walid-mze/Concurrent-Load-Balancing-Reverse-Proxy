package main

import (
	"ReverseProxy/config"
	cheker "ReverseProxy/healthCheker"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	//load backends
	pool, err := config.LoadBackends("reverse-proxy/config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	fmt.Printf("Loaded %d backends:\n", len(pool.Backends))

	proxyHandler := &ProxyHandler{
		ServerPool: pool,
	}

	//test the health cheker function

	proxyConfig, err := config.LoadProxyConfig("reverse-proxy/proxyConfig.json")
	if err != nil {
		log.Fatalf("Error loading proxy config: %v", err)
	}

	healthCheckFreq, err := time.ParseDuration(proxyConfig.HealthChekerFreq)
	if err != nil {
		log.Fatalf("Error parsing health check frequency: %v", err)
	}

	cheker.StartHealthCheck(pool, healthCheckFreq)
	addr := fmt.Sprintf(":%d", proxyConfig.Port)

	fmt.Printf("starting proxy server on %s\n", addr)
	http.ListenAndServe(addr, proxyHandler)

}
