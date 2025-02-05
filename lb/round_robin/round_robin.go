package roundrobin

import (
	"fmt"
	"net/http"
	"load-balancer/server"
	"load-balancer/lb"
	"load-balancer/random_gen"
	"time"
)

type loadBalancerImpl struct {
	servers []server.Server
	serverIndex int
	done chan struct{}
	isServerHealthy map[string]bool
}

const serverHealthCheckURLPatten = "http://localhost:%s/health"


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
		done: make(chan struct{}),
		isServerHealthy: make(map[string]bool),
	}
}

func (lb *loadBalancerImpl) FindServer() server.Server {
	counter := 0
	for counter < len(lb.servers){
		lb.serverIndex++
		if lb.serverIndex >= len(lb.servers) {
			lb.serverIndex %= len(lb.servers)
		}
		if lb.isServerHealthy[lb.servers[lb.serverIndex].GetName()] {
			return lb.servers[lb.serverIndex]
		}
		counter++
	}
	return nil
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

func (lb *loadBalancerImpl) PerformHealthChecks() {
	for _, server := range lb.servers {
		lb.updateServerHealth(server.GetName())
	}
}

func (lb *loadBalancerImpl) Start(port string) {
	
	for _,server := range lb.servers {
		server.Start()
	}
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", lb.handleRequest)
	go func() {
		http.ListenAndServe(":" + port, mux)
	} ()


	go func() {
		for {
			select {
			case <-lb.done:
				fmt.Println("done called for round robin load balancer")
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
