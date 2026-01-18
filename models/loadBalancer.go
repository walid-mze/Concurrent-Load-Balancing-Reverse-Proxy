package models 

import  ("net/url")


type LoadBalancer interface {
	GetNextValidPeer() *Backend
	AddBackend(Backend *Backend) 
	SetBackendStatus(uri *url.URL, alive bool) 
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

func (s *ServerPool)SetBackendStatus(uri *url.URL, alive bool){
	for i:=0;i<len(s.Backends);i++{
		if s.Backends[i].URL.String()==uri.String(){
			s.Backends[i].Mux.Lock()
			s.Backends[i].Alive=alive
			s.Backends[i].Mux.Unlock()
			return 
		}
	}

}
