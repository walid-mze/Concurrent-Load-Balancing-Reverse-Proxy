package models

import  "sync"

type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current  uint64     `json:"current"`
	Mux 	sync.RWMutex
}

func (s *ServerPool) AddBackend(b *Backend) {
	s.Backends = append(s.Backends, b)
}

func (s *ServerPool) GetNextValidPeer() *Backend{
	l:=len(s.Backends)

	s.Mux.Lock()
	next:=s.Current
	s.Current++
	s.Mux.Unlock()

	for i:=0;i<l;i++{
		idx:=int(next+uint64(i))%l
		s.Backends[idx].Mux.RLock()
		alive:=s.Backends[idx].Alive
		s.Backends[idx].Mux.RUnlock()

		if alive{
			return s.Backends[idx]
		}
	}
	return nil

}
