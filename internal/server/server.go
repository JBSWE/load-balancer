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
		for range time.Tick(healthCheckInterval) {
			start := time.Now()
			res, err := http.Head(s.URL)
			s.Latency = time.Since(start)
			if err != nil || res.StatusCode != http.StatusOK {
				logger.Warn("Server is down", zap.String("url", s.URL), zap.Error(err))
				s.IsHealthy = false
			} else {
				if s.Latency < 2*time.Second {
					s.IsHealthy = true
				} else {
					s.IsHealthy = false
					s.ExclusionTime = time.Now()
				}
			}
		}
	}()
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
