package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
	"time"
)

type Payload struct {
	Game    string `json:"game"`
	GamerID string `json:"gamerID"`
	Points  int    `json:"points"`
}

func waitForLoadBalancer() error {
	client := &http.Client{Timeout: 10 * time.Second}
	healthCheckPayload := `{"game": "healthcheck", "gamerID": "healthcheck", "points": 0}`
	for i := 0; i < 30; i++ {
		resp, err := client.Post("http://localhost:8080", "application/json", bytes.NewBuffer([]byte(healthCheckPayload)))
		if err != nil {
			fmt.Printf("Attempt #%d: Error reaching load balancer: %v\n", i+1, err)
		} else {
			fmt.Printf("Attempt #%d: Load balancer response code: %d\n", i+1, resp.StatusCode)
		}

		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("load balancer not ready")
}

func TestLoadBalancerRouting(t *testing.T) {
	err := waitForLoadBalancer()
	require.NoError(t, err, "Load balancer is not ready")

	testPayloads := []Payload{
		{
			Game:    "COD",
			GamerID: "GYUTDTE",
			Points:  20,
		},
		{
			Game:    "Tetris",
			GamerID: "ABCD1234",
			Points:  15,
		},
		{
			Game:    "Snake",
			GamerID: "SNAKE001",
			Points:  50,
		},
		{
			Game:    "Pac-Man",
			GamerID: "PACMAN123",
			Points:  100,
		},
		{
			Game:    "Space Invaders",
			GamerID: "INVADER555",
			Points:  75,
		},
	}

	var visitedServers []string

	for _, testPayload := range testPayloads {
		t.Run(fmt.Sprintf("Testing payload: %s", testPayload.Game), func(t *testing.T) {
			payloadBytes, err := json.Marshal(testPayload)
			require.NoError(t, err)

			client := &http.Client{}
			resp, err := client.Post("http://localhost:8080", "application/json", bytes.NewBuffer(payloadBytes))
			require.NoError(t, err, "Error sending request to load balancer")

			assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK from load balancer")

			var responsePayload Payload
			err = json.NewDecoder(resp.Body).Decode(&responsePayload)
			require.NoError(t, err, "Error decoding response body")
			resp.Body.Close()

			assert.Equal(t, testPayload, responsePayload, "Response body does not match the sent payload")

			xForwardedServer := resp.Header.Get("X-Forwarded-Server")
			serverName := strings.Split(strings.Split(xForwardedServer, "://")[1], ":")[0]

			backendServers := []string{
				"server1",
				"server2",
				"server3",
				"server4",
				"server5",
			}

			assert.Contains(t, backendServers, serverName, "Routing to an unexpected server")

			visitedServers = append(visitedServers, serverName)

			if len(visitedServers) > 1 {
				for i := 1; i < len(visitedServers); i++ {
					if visitedServers[i] == visitedServers[i-1] {
						t.Errorf("Round-robin failed: Consecutive requests routed to the same server (%s)", visitedServers[i])
					}
				}
			}
		})
	}

	assert.Len(t, visitedServers, len(testPayloads), "Not all requests were routed to different servers")
}
