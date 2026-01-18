package models

import  ("sync" )

type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current  uint64     `json:"current"`
	Mux 	sync.RWMutex
}


