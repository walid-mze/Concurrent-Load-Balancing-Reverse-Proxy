package main

import (
	"ReverseProxy/models"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
	"context"
)

type ProxyHandler struct {
	ServerPool *models.ServerPool
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := p.ServerPool.GetNextValidPeer()

	if backend == nil {
		http.Error(w, "Service Unavailable No backends available", http.StatusServiceUnavailable)
		log.Println("No backends available")
		return
	}
	backend.IncrementConns()
	defer backend.DecrementConns()
	time.Sleep(5 * time.Second)



	//creating context with timeout
	ctx,cancel:=context.WithTimeout(r.Context(),30*time.Second)
	defer cancel()
	r=r.WithContext(ctx)


	log.Printf("Forwarding request to %s", backend.URL)
	reverseProxy := httputil.NewSingleHostReverseProxy(backend.URL)


	reverseProxy.ErrorHandler=func(w http.ResponseWriter,r *http.Request, err error){
		log.Printf("Backend %s error:%v",backend.URL,err)
		p.ServerPool.SetBackendStatus(backend.URL,false)
		http.Error(w,"Bad Gateway",http.StatusBadGateway)
	}
	reverseProxy.ServeHTTP(w, r)
}
