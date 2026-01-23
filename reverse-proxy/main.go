package main

import (
	"ReverseProxy/config"
	"ReverseProxy/healthCheker"
	"fmt"
	"log"
	"net/http"
	"time"
	"ReverseProxy/admin"
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


	//start the proxy server in a separate go routine
	go func() {
	fmt.Printf("starting proxy server on %s\n", addr)
	http.ListenAndServe(addr, proxyHandler)
	time.Sleep(5 * time.Second)

	}()


	// test the adminAPI
	admin:=admin.AdminAPI{
		ServerPool: pool,
	}
	http.HandleFunc("/status",admin.StatusHandler)

	http.HandleFunc("/backends",func (w http.ResponseWriter, r *http.Request){
		switch r.Method{
		case http.MethodPost:
					admin.AddBackendHandler(w,r)
		case http.MethodDelete:
					admin.DeleteBackendHandler(w,r)
		default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	})
	
	//start the admin server in a separate go routine
	go func (){
		fmt.Printf("starting proxy server on :8081\n")

		http.ListenAndServe(":8081",nil)
	}()

	select {}


	
/*issues in the code : 
	1. when i add a backend from the admin api the round robin does not work correctly : 
		2026/01/21 22:02:38 Forwarding request to http://localhost:8082
		2026/01/21 22:02:39 Forwarding request to http://localhost:8083
		2026/01/21 22:02:48 Forwarding request to http://localhost:8085
		2026/01/21 22:02:52 Forwarding request to http://localhost:8085

	==> solved 

	2.the http://localhost:8081/status have always the curConn=0 for all backend 

	3.the delete func is not working : 
	==> solved 
*/




}
