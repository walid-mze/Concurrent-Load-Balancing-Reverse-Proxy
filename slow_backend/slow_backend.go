package main

import (
	"fmt"
	"net/http"
	"time"
)

// this is a slow backend that sleeps for 60 seconds before responding
// I used for testing the cancellation of requests due to timeout
func main() {
	// Health check endpoint - responds quickly
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Slow endpoint for testing timeout
	http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Backend received slow request, sleeping for 60 seconds...")
		time.Sleep(60 * time.Second)
		w.Write([]byte("Done"))
	})

	// Root endpoint - also slow for testing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Backend received request, sleeping for 60 seconds...")
		time.Sleep(60 * time.Second)
		w.Write([]byte("Done"))
	})

	fmt.Println("slow backend running on :8085")
	http.ListenAndServe(":8085", nil)
}
