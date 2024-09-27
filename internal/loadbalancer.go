package internal

import (
	"log"
	"net/http"
	"time"
)

type balancer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	HealthCheck()
}

type LoadBalancer struct {
	pool balancer
	algo string
}

func NewLoadBalancer(s []*Server, a string) *LoadBalancer {
	var p balancer
	switch a {
	case "round-robin":
		p = NewRRPool(s)
	}
	return &LoadBalancer{
		pool: p,
		algo: a,
	}
}

func (lb *LoadBalancer) Serve(addr string) {
	t := time.NewTicker(time.Second * 2)
	go func() {
		defer t.Stop()
		for {
			<-t.C
			log.Println("Health check underway")
			lb.pool.HealthCheck()
		}
	}()

	s := &http.Server{
		Addr:    addr,
		Handler: lb.pool,
	}

	log.Fatal(s.ListenAndServe())
}
