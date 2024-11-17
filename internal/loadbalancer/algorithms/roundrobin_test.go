package algorithms

import (
	"github.com/JBSWE/load-balancer/internal/loadbalancer"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func createTestServer(urlStr string, healthy bool, exclusionTime time.Time) *loadbalancer.Server {
	return &loadbalancer.Server{
		URL:           urlStr,
		IsHealthy:     healthy,
		ExclusionTime: exclusionTime,
	}
}

func TestRoundRobin_BasicRotation(t *testing.T) {
	servers := []*loadbalancer.Server{
		createTestServer("http://server1", true, time.Time{}),
		createTestServer("http://server2", true, time.Time{}),
		createTestServer("http://server3", true, time.Time{}),
	}

	lb := NewRoundRobin()

	assert.Equal(t, "http://server1", lb.GetServer(servers).URL)
	assert.Equal(t, "http://server2", lb.GetServer(servers).URL)
	assert.Equal(t, "http://server3", lb.GetServer(servers).URL)
	assert.Equal(t, "http://server1", lb.GetServer(servers).URL)
}

func TestRoundRobin_SkipUnhealthyServers(t *testing.T) {
	servers := []*loadbalancer.Server{
		createTestServer("http://server1", true, time.Time{}),
		createTestServer("http://server2", false, time.Time{}),
		createTestServer("http://server3", true, time.Time{}),
	}

	lb := NewRoundRobin()

	assert.Equal(t, "http://server1", lb.GetServer(servers).URL)
	assert.Equal(t, "http://server3", lb.GetServer(servers).URL)
	assert.Equal(t, "http://server1", lb.GetServer(servers).URL)
}

func TestRoundRobin_HealthyServersOnly(t *testing.T) {
	servers := []*loadbalancer.Server{
		createTestServer("http://server1", true, time.Time{}),
		createTestServer("http://server2", true, time.Time{}),
	}

	lb := NewRoundRobin()

	assert.Equal(t, "http://server1", lb.GetServer(servers).URL)
	assert.Equal(t, "http://server2", lb.GetServer(servers).URL)
	assert.Equal(t, "http://server1", lb.GetServer(servers).URL)
}

func TestRoundRobin_SingleServer(t *testing.T) {
	servers := []*loadbalancer.Server{
		createTestServer("http://server1", true, time.Time{}),
	}

	lb := NewRoundRobin()

	assert.Equal(t, "http://server1", lb.GetServer(servers).URL)
	assert.Equal(t, "http://server1", lb.GetServer(servers).URL)
}

func TestRoundRobin_Concurrency(t *testing.T) {
	servers := []*loadbalancer.Server{
		createTestServer("http://server1", true, time.Time{}),
		createTestServer("http://server2", true, time.Time{}),
	}

	lb := NewRoundRobin()

	var result []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				mu.Lock()
				result = append(result, lb.GetServer(servers).URL)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	assert.Contains(t, result, "http://server1")
	assert.Contains(t, result, "http://server2")
	assert.Equal(t, 20, len(result))
}

func TestRoundRobin_ServerListExhaustion(t *testing.T) {
	servers := []*loadbalancer.Server{
		createTestServer("http://server1", true, time.Time{}),
		createTestServer("http://server2", true, time.Time{}),
		createTestServer("http://server3", true, time.Time{}),
	}

	lb := NewRoundRobin()

	for i := 0; i < 10; i++ {
		server := lb.GetServer(servers)
		assert.Contains(t, []string{"http://server1", "http://server2", "http://server3"}, server.URL)
	}
}

func TestRoundRobin_ExclusionTime(t *testing.T) {
	currentTime := time.Now()

	server1, err := loadbalancer.NewServer("http://server1", true, currentTime.Add(time.Second))
	if err != nil {
		t.Fatalf("Failed to create server1: %v", err)
	}

	server2, err := loadbalancer.NewServer("http://server2", true, time.Time{})
	if err != nil {
		t.Fatalf("Failed to create server2: %v", err)
	}

	servers := []*loadbalancer.Server{server1, server2}

	lb := NewRoundRobin()

	selectedServer := lb.GetServer(servers)
	t.Logf("Selected server initially: %v", selectedServer.URL)
	assert.Equal(t, "http://server2", selectedServer.URL)

	time.Sleep(3 * time.Second)

	selectedServer = lb.GetServer(servers)
	t.Logf("Selected server after exclusion time: %v", selectedServer.URL)

	assert.Equal(t, "http://server1", selectedServer.URL)
}
