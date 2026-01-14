package models

import (
	"net/url"
	"sync"
)

type Backend struct {
	URL                *url.URL `json:"url"`
	Alive              bool     `json:"alive"`
	CurrentConnections int      `json:"current_connections"`
	Mux                sync.RWMutex
}
