package leastconnection

import (
	"fmt"
	"load-balancer/lb"
	randomgen "load-balancer/random_gen"
	"load-balancer/server"
	"net/http"
	"sync"
	"time"
)

type loadBalancerImpl struct {
	servers []server.Server
	mu sync.RWMutex
	done chan struct{}
	isServerHealthy map[string]bool
}

const serverHealthCheckURLPatten = "http://localhost:%s/health"

func NewLeastConnectionLoadBalancer(count int) lb.LoadBalancer {
	servers := []server.Server{}
	for count > 0 {
		port := randomgen.RandRange()
		server := server.NewServer(port, port)
		servers = append(servers, server)
		count--
	}
	return &loadBalancerImpl {
		servers: servers,
		done: make(chan struct{}),
		isServerHealthy: make(map[string]bool),
	}
}

func (lb *loadBalancerImpl) updateServerHealth (serverPort string) {
	healthCheckURL := fmt.Sprintf(serverHealthCheckURLPatten, serverPort)
	resp, err := http.Get(healthCheckURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("Server %s is unhealthy\n", serverPort)
		lb.isServerHealthy[serverPort] = false
		return
	}
	fmt.Printf("Server %s is HEALTHY\n", serverPort)
	lb.isServerHealthy[serverPort] = true
}


func (lb *loadBalancerImpl) FindServer() server.Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	var minServer server.Server
	for _, server := range lb.servers {
		if !lb.isServerHealthy[server.GetName()] {
			continue
		}
		if minServer == nil ||  server.GetActiveConnections() < minServer.GetActiveConnections(){
			minServer = server
		}
	}
	return minServer
}

func (lb *loadBalancerImpl) PerformHealthChecks() {
	for _, server := range lb.servers {
		lb.updateServerHealth(server.GetName())
	}
}

func (lb *loadBalancerImpl) handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Request received to least connection LB")

	server := lb.FindServer()
	server.HandleRequest(w, r)
}

func (lb *loadBalancerImpl) Start(port string) {
	for _,server := range lb.servers {
		server.Start()
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", lb.handleRequest)
	go http.ListenAndServe(":" + port, mux)
	go func() {
		for {
			select {
			case <-lb.done:
				fmt.Println("done called for least connection load balancer")
				return
			default:
				lb.PerformHealthChecks()
				time.Sleep(5 * time.Second)
			}
		}
		
	}()
}

func (lb *loadBalancerImpl) Stop() {
	for _,server := range lb.servers {
		server.Stop()
	}
	lb.done <- struct{}{}
}

func (lb *loadBalancerImpl) AddServer(backendServer server.Server) {
	lb.servers = append(lb.servers, backendServer)
}

func (lb *loadBalancerImpl) RemoveServer(serverName string) {
	servers := []server.Server{}
	for _, server := range lb.servers {
		if(server.GetName() == serverName) {
			continue
		}
		servers = append(servers, server)
	}
	lb.servers = servers
}
