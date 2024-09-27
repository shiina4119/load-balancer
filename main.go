package main

import (
	"flag"
	"fmt"
	"lb/internal"
	"log"
	"net/url"
)

func main() {
	port := flag.Int("port", 8090, "port number")
	algo := flag.String("algo", "round-robin", "algorithm to use")
	flag.Parse()
	urls := flag.Args()

	servers := make([]*internal.Server, 0)
	for _, addr := range urls {
		url, err := url.Parse(addr)
		if err != nil {
			log.Fatal(err)
		}
		servers = append(servers, internal.CreateServer(url))
	}

	lb := internal.NewLoadBalancer(servers, *algo)
	fmt.Printf("Load Balancer started on port %d", *port)
	lb.Serve(fmt.Sprintf(":%d", *port))
}
