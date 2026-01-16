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

	//test Round Robin logic 
	for i:=0;i<10;i++{
		backend:=pool.GetNextValidPeer()
		if backend!=nil{
			fmt.Printf("Request %d -> %s\n",i+1,backend.URL)
		}else{
            fmt.Printf("Request %d -> No backends available\n", i+1)

		}
	}	

}
