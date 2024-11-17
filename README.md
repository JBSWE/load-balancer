# Load Balancer API

This project implements a **Round Robin Load Balancer** that routes HTTP POST requests to a set of application APIs, distributing requests in a round-robin manner. The system handles server health checks, slow server exclusion, and provides a basic example of a load balancing mechanism with failover support.

## 1. Overview

The project is composed of two main components:
- **Application API**: A simple API that accepts POST requests and returns the received JSON payload as the response.
- **Round Robin API**: The load balancer that accepts POST requests, determines which application API instance to route the request to, and forwards the request accordingly.

### How it Works
The Round Robin API uses Go’s `http.ReverseProxy` to direct load based on a round-robin implementation. If the server is healthy and not excluded due to slow responses, the server is chosen in a round-robin approach (e.g., 1, 2, 3, 4). If the server is unhealthy or excluded, the system skips it and selects the next healthy server.

## 2. Run Docker Compose

To start the Round Robin API along with multiple instances of the Application API, use the following command:

```bash
docker-compose up
```

This will start:

The Round Robin API on port 8080
Multiple instances of the Application API (configured in docker-compose.yml)

## 3. Test Manually
You can test the setup by sending a POST request to the Round Robin API. For example:

```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"game": "COD", "gamerID": "GYUTDTE", "points": 10}'
```

This will send a request to the Round Robin API, which will then route it to one of the available Application API instances based on the round-robin logic.

### How the Round Robin API Handles Failures and Slow Servers
#### 1. Handling Application API Failures
If an Application API instance goes down, the Round Robin API will detect this using health checks. If the instance is not responsive, the Round Robin API will mark it as unhealthy and exclude it from the round-robin rotation. The health checks occur periodically, and once the server is back online, it will be included again in the rotation.

#### Implementation Details:
- The Round Robin API pings each Application API instance at regular intervals using HTTP HEAD requests.
- If the response status code is not 200 OK or if an error occurs, the server is marked as unhealthy.

#### 2. Handling Slow Application API Servers

If an Application API instance starts responding slowly (i.e., if the response latency exceeds a predefined threshold), the Round Robin API will temporarily exclude the server using an `ExclusionTime`. The instance will not be included in the round-robin rotation until it recovers, based on a certain recovery time or health check status.

#### Implementation Details:
The ExclusionTime field in the server object is used to mark a server that should not be selected for a period of time.
If the server's latency exceeds a certain threshold (e.g., 2 seconds), it is marked as "slow" and excluded from the round-robin rotation.

## 4. Testing
### Unit Tests
Unit tests for the Round Robin API and Application API can be run using the following command:
```bash
make test
```
This will run all unit tests and validate the behavior of the application components, including the health check and round-robin routing logic.

### Integration Tests

Integration tests ensure that the Round Robin API correctly routes requests to the available Application API instances and handles failures and slow responses appropriately. You can run the integration tests using:
```bash
make run-integration-tests
```

This will start up the system (using Docker Compose if configured) and perform end-to-end testing to verify the correctness of the round-robin routing and failure handling.

## 5. Project Structure
 - `/cmd/server`: The entry point for the Round Robin API server.
 - `/internal/server`: Contains the core logic for handling requests and managing round-robin routing, health checks, and server exclusions.
- `/internal/loadbalancer`: Contains the logic for the load balancing (round robin) and server management (health checks and exclusion).
- `/internal/config`: Configuration for server parameters and health check settings.
- `docker-compose.yml`: Docker Compose configuration to run multiple instances of the Application API and the Round Robin API.

## 6. Future Improvements
- Dynamic Scaling: Allow the Round Robin API to automatically discover new Application API instances based on service discovery or configuration updates.
- Advanced Health Checks: Implement more advanced health checks, such as checking server response time, load, or other metrics.
- Metrics and Monitoring: Integrate with a monitoring system to collect metrics about server health, request distribution, and performance.
- Retry Logic: Implement retry logic for failed requests to slow or temporarily failed servers.
- Configuration Management: Support for different configurations of the round-robin algorithm (e.g., weighted round robin) or different health check intervals.

## 7. Conclusion
This project demonstrates a simple, robust, and scalable implementation of a Round Robin Load Balancer API with features such as health checks and exclusion for slow servers. It can be easily extended with additional features such as metrics, retries, or dynamic scaling.