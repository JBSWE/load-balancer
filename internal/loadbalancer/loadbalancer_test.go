package loadbalancer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	url := "http://localhost:8081"
	isHealthy := true
	exclusionTime := time.Now().Add(1 * time.Hour)

	server, err := NewServer(url, isHealthy, exclusionTime)

	assert.NoError(t, err)
	assert.NotNil(t, server)
	assert.Equal(t, url, server.URL)
	assert.Equal(t, isHealthy, server.IsHealthy)
	assert.Equal(t, exclusionTime, server.ExclusionTime)
}

func TestIsExcludable_NoExclusionTime(t *testing.T) {

	server := &Server{
		URL:           "http://localhost:8081",
		IsHealthy:     true,
		ExclusionTime: time.Time{},
	}

	excludable := server.IsExcludable()

	assert.False(t, excludable)
}

func TestIsExcludable_ExclusionTimeInThePast(t *testing.T) {
	server := &Server{
		URL:           "http://localhost:8081",
		IsHealthy:     true,
		ExclusionTime: time.Now().Add(-1 * time.Hour),
	}

	excludable := server.IsExcludable()

	assert.False(t, excludable)
}

func TestIsExcludable_ExclusionTimeInTheFuture(t *testing.T) {
	server := &Server{
		URL:           "http://localhost:8081",
		IsHealthy:     true,
		ExclusionTime: time.Now().Add(1 * time.Hour),
	}

	excludable := server.IsExcludable()

	assert.True(t, excludable)
}
