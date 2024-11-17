package loadbalancer

import (
	"time"
)

type Server struct {
	URL           string
	IsHealthy     bool
	ExclusionTime time.Time
	Latency       time.Duration
}

func NewServer(urlString string, isHealthy bool, exclusionTime time.Time) (*Server, error) {

	return &Server{
		URL:           urlString,
		IsHealthy:     isHealthy,
		ExclusionTime: exclusionTime,
	}, nil
}

func (s *Server) IsExcludable() bool {
	return !s.ExclusionTime.IsZero() && time.Now().Before(s.ExclusionTime)
}

type LoadBalancer interface {
	GetServer(servers []*Server) *Server
}
