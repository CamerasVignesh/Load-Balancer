package leastconnection

import (
	"fmt"
	"load-balancer/lb"
	randomgen "load-balancer/random_gen"
	"load-balancer/server"
	"net/http"
)

type loadBalancerImpl struct {
	servers []server.Server
}

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
	}
}

func (lb *loadBalancerImpl) FindServer() server.Server {
	minServer := lb.servers[0]
	for _, server := range lb.servers {
		if(server.GetActiveConnections() < minServer.GetActiveConnections()) {
			minServer = server
		}
	}
	return minServer
}

func (lb *loadBalancerImpl) handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Request received to round robin LB")

	server := lb.FindServer()
	server.HandleRequest(w, r)
}

func (lb *loadBalancerImpl) Start(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", lb.handleRequest)
	go func() {
		http.ListenAndServe(":" + port, mux)
	}()

	for _,server := range lb.servers {
		server.Start()
	}
}

func (lb *loadBalancerImpl) Stop() {
	for _,server := range lb.servers {
		server.Stop()
	}
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



