package main

import (
	"ReverseProxy/admin"
	"ReverseProxy/config"
	cheker "ReverseProxy/healthCheker"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	// Create proxy server with http.Server for graceful shutdown
	proxyServer := &http.Server{
		Addr:    addr,
		Handler: proxyHandler,
	}

	//start the proxy server in a separate go routine
	go func() {
		fmt.Printf("starting proxy server on %s\n", addr)
		if err := proxyServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("proxy server error: %v", err)
		}
	}()

	// test the adminAPI
	adminAPI := admin.AdminAPI{
		ServerPool: pool,
	}

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/status", adminAPI.StatusHandler)

	adminMux.HandleFunc("/backends", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			adminAPI.AddBackendHandler(w, r)
		case http.MethodDelete:
			adminAPI.DeleteBackendHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Create admin server with http.Server for graceful shutdown
	adminServer := &http.Server{
		Addr:    ":8081",
		Handler: adminMux,
	}

	//start the admin server in a separate go routine
	go func() {
		fmt.Printf("starting admin server on :8081\n")
		if err := adminServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("admin server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\nShutting down servers...")

	// create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := proxyServer.Shutdown(ctx); err != nil {
		log.Printf("proxy server shutdown error: %v", err)
	} else {
		fmt.Println("proxy server stopped gracefully")
	}

	if err := adminServer.Shutdown(ctx); err != nil {
		log.Printf("admin server shutdown error: %v", err)
	} else {
		fmt.Println("admin server stopped gracefully")
	}
	fmt.Println("All servers stopped!")

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
