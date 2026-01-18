package models

import (
	"net/url"
	"sync"
)

type Backend struct {
	URL                *url.URL `json:"url"`
	Alive              bool     `json:"alive"`
	CurrentConnections int64     `json:"current_connections"`
	Mux                sync.RWMutex
}
func (b *Backend) IncrementConns(){
	b.Mux.Lock()
	b.CurrentConnections++
	b.Mux.Unlock()
}

func (b *Backend) DecrementConns(){
	b.Mux.Lock()
	b.CurrentConnections--
	b.Mux.Unlock()

}

