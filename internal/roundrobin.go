package internal

import (
	"lb/queue"
	"log"
	"net/http"
	"sync"
)

type RRPool struct {
	queue     *queue.CircularQueue[*Server]
	serverMap map[*Server]bool
	mu        sync.Mutex
}

func NewRRPool(servers []*Server) *RRPool {
	n := len(servers)
	q := queue.NewCircularQueue[*Server](n)
	m := make(map[*Server]bool)
	for _, server := range servers {
		err := q.Push(server)
		if err != nil {
			log.Println(err)
		}
		m[server] = true
	}
	return &RRPool{
		queue:     q,
		serverMap: m,
	}
}

func (rp *RRPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	var server *Server
	var err error
	for {
		server, err = rp.queue.Pop()
		if err != nil {
			log.Println("All servers are down")
			w.WriteHeader(503)
			return
		}
		v, err := isAlive(server)
		if v {
			log.Printf("Request sent to %s\n", server._url.String())
			server.proxy.ServeHTTP(w, r)
			defer func() {
				err := rp.queue.Push(server)
				if err != nil {
					log.Println(err)
				}
			}()
			break
		} else {
			log.Printf("Cannot send request to server %s: %s\n", server._url.String(), err.Error())
			rp.serverMap[server] = false
		}
	}
}

func (rp *RRPool) HealthCheck() {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	log.Println("Health check underway")
	for server, alive := range rp.serverMap {
		v, err := isAlive(server)
		if v {
			log.Printf("Server %s up\n", server._url.String())
			if !alive {
				rp.serverMap[server] = true
				e := rp.queue.Push(server)
				if e != nil {
					log.Println(e)
				}
			}
		} else {
			log.Printf("Server %s down: %s\n", server._url.String(), err.Error())
		}
	}
}
