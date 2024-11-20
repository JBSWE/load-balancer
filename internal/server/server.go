package server

import (
	"github.com/JBSWE/load-balancer/internal/config"
	"github.com/JBSWE/load-balancer/internal/loadbalancer"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func setupHealthCheck(s *loadbalancer.Server, healthCheckInterval time.Duration, logger *zap.Logger) {
	go func() {
		for {
			s.Latency = performHealthCheck(s)

			if !isHealthy(s.Latency, healthCheckInterval) {
				s.IsHealthy = false
				s.ExclusionTime = time.Now()
				logger.Warn("Server excluded due to failure or high latency",
					zap.String("url", s.URL),
					zap.Duration("latency", s.Latency))
			} else {
				s.IsHealthy = true
				logger.Info("Server health check passed",
					zap.String("url", s.URL),
					zap.Duration("latency", s.Latency))
			}

			s.Mu.Unlock()
		}
	}()
}

func performHealthCheck(s *loadbalancer.Server) time.Duration {
	start := time.Now()
	res, err := http.Head(s.URL)
	s.Mu.Lock()
	if err != nil || res.StatusCode != http.StatusOK {
		return 0
	}
	return time.Since(start)
}

func isHealthy(latency time.Duration, latencyThreshold time.Duration) bool {
	return latency > 0 && latency < latencyThreshold
}

func StartServer(cfg config.Config, lb loadbalancer.LoadBalancer, servers []*loadbalancer.Server, logger *zap.Logger) {
	healthCheckInterval, err := time.ParseDuration(cfg.HealthCheckInterval)
	if err != nil {
		logger.Fatal("Failed to parse health check interval", zap.String("interval", cfg.HealthCheckInterval), zap.Error(err))
	}

	for _, server := range servers {
		setupHealthCheck(server, healthCheckInterval, logger)
	}

	SetupRoutes(lb, servers, logger)

	logger.Info("Starting load balancer", zap.String("port", cfg.Port))
	err = http.ListenAndServe(cfg.Port, nil)
	if err != nil {
		logger.Fatal("Error starting load balancer", zap.Error(err))
	}
}
