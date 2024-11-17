package server

import (
	"bytes"
	"encoding/json"
	"github.com/JBSWE/load-balancer/internal/loadbalancer"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func SetupRoutes(lb loadbalancer.LoadBalancer, servers []*loadbalancer.Server, logger *zap.Logger) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server := lb.GetServer(servers)
		if server == nil {
			logger.Error("No healthy server available")
			http.Error(w, "No healthy server available", http.StatusServiceUnavailable)
			return
		}

		w.Header().Add("X-Forwarded-Server", server.URL)

		proxy := setupReverseProxy(server, logger)

		logger.Info("Proxying request", zap.String("target", server.URL))

		proxy.ServeHTTP(w, r)
	})
}

func setupReverseProxy(server *loadbalancer.Server, logger *zap.Logger) http.Handler {
	proxy := http.NewServeMux()

	proxy.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("Error reading request body", zap.Error(err))
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		logger.Info("Request Body:", zap.ByteString("body", bodyBytes))

		client := &http.Client{}

		req, err := http.NewRequest(r.Method, server.URL, bytes.NewReader(bodyBytes))
		if err != nil {
			logger.Error("Error creating new request", zap.Error(err))
			http.Error(w, "Failed to create request to backend", http.StatusInternalServerError)
			return
		}

		req.Header = r.Header

		resp, err := client.Do(req)
		if err != nil {
			logger.Error("Error forwarding request to backend", zap.Error(err))
			http.Error(w, "Failed to forward request to backend", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("Error reading response body", zap.Error(err))
			http.Error(w, "Failed to read response body", http.StatusInternalServerError)
			return
		}

		var jsonResponse map[string]interface{}
		if err := json.Unmarshal(respBody, &jsonResponse); err != nil {
			logger.Error("Error unmarshaling JSON response", zap.Error(err))
			http.Error(w, "Failed to parse response body", http.StatusInternalServerError)
			return
		}

		if jsonContent, ok := jsonResponse["json"]; ok {
			jsonResponseBytes, err := json.Marshal(jsonContent)
			if err != nil {
				logger.Error("Error marshaling json content", zap.Error(err))
				http.Error(w, "Failed to marshal JSON content", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResponseBytes)
			return
		} else {
			logger.Error("No json field found in the response body")
			http.Error(w, "No JSON content in the response", http.StatusInternalServerError)
			return
		}
	})

	return proxy
}
