package roundrobin

import (
	"fmt"
	"net/http"
	"load-balancer/server"
	"load-balancer/lb"
	"load-balancer/random_gen"
)

type loadBalancerImpl struct {
	servers []server.Server
	serverIndex int
}


func NewRoundRobinLoadBalancer(count int) lb.LoadBalancer {
	servers := []server.Server{}
	for count > 0 {
		newServerPort := randomgen.RandRange()
		newServer := server.NewServer(newServerPort, newServerPort)
		servers = append(servers, newServer)
		count--
	}
	return &loadBalancerImpl{
		servers: servers,
	}
}

func (lb *loadBalancerImpl) FindServer() server.Server {
	lb.serverIndex++
	if lb.serverIndex >= len(lb.servers) {
		lb.serverIndex %= len(lb.servers)
	}
	return lb.servers[lb.serverIndex]
}

func (lb *loadBalancerImpl) handleRequest (w http.ResponseWriter, r *http.Request) {
	/*
	fmt.Println("./lb")
    fmt.Printf("Received request from %s\n", r.RemoteAddr)
    fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
    fmt.Printf("Host: %s\n", r.Host)
    fmt.Printf("User-Agent: %s\n", r.Header.Get("User-Agent"))
    fmt.Printf("Accept: %s\n", r.Header.Get("Accept"))
	*/
    
    fmt.Fprintln(w, "Request received to round robin LB")

	server := lb.FindServer()
	server.HandleRequest(w, r)
}

func (lb *loadBalancerImpl) Start(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", lb.handleRequest)
	go func() {
		http.ListenAndServe(":" + port, mux)
	} ()

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
