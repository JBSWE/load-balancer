package algorithms

import (
	"github.com/JBSWE/load-balancer/internal/loadbalancer"
	"sync"
)

type RoundRobin struct {
	Current int
	Mutex   sync.Mutex
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{Current: 0}
}

func (lb *RoundRobin) GetServer(servers []*loadbalancer.Server) *loadbalancer.Server {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()

	for i := 0; i < len(servers); i++ {
		idx := lb.Current % len(servers)
		nextServer := servers[idx]
		lb.Current++

		if nextServer.IsHealthy && !nextServer.IsExcludable() {
			return nextServer
		}
	}

	return nil
}
