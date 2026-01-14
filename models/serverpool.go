package models

type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current  uint64     `json:"current"`
}

func (s *ServerPool) AddBackend(b *Backend) {
	s.Backends = append(s.Backends, b)
}
