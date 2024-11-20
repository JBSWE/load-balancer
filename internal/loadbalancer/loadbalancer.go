package loadbalancer

import (
	"sync"
	"time"
)

type Server struct {
	URL           string
	IsHealthy     bool
	ExclusionTime time.Time
	Latency       time.Duration
	Mu            sync.Mutex
}

func NewServer(urlString string, isHealthy bool, exclusionTime time.Time) (*Server, error) {

	return &Server{
		URL:           urlString,
		IsHealthy:     isHealthy,
		ExclusionTime: exclusionTime,
		Mu:            sync.Mutex{},
	}, nil
}

func (s *Server) IsExcludable() bool {
	return !s.ExclusionTime.IsZero() && time.Now().Before(s.ExclusionTime)
}

type LoadBalancer interface {
	GetServer(servers []*Server) *Server
}
