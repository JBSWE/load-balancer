package main

import (
	"github.com/JBSWE/load-balancer/internal/config"
	"github.com/JBSWE/load-balancer/internal/loadbalancer"
	"github.com/JBSWE/load-balancer/internal/loadbalancer/algorithms"
	"github.com/JBSWE/load-balancer/internal/server"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	lb := algorithms.NewRoundRobin()

	var servers []*loadbalancer.Server
	for _, serverURL := range cfg.Servers {
		server := &loadbalancer.Server{URL: serverURL, IsHealthy: true}
		servers = append(servers, server)
	}

	server.StartServer(*cfg, lb, servers, logger)
}
