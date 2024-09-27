package internal

import "net/http"

type loadBalancer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	HealthCheck()
}
