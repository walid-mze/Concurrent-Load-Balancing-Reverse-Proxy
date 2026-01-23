package models

import (
	"net/url"
	"sync"
	"sync/atomic"

)

type Backend struct {
	URL                *url.URL `json:"url"`
	Alive              bool     `json:"alive"`
	CurrentConnections int64     `json:"current_connections"`
	Mux                sync.RWMutex
}

func (b *Backend) IncrementConns() {
    atomic.AddInt64(&b.CurrentConnections, 1)
}

func (b *Backend) DecrementConns() {
    atomic.AddInt64(&b.CurrentConnections, -1)
}



