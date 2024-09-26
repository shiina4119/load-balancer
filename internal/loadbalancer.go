package internal

import "net/http"

type LoadBalancer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
