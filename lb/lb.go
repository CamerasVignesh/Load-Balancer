package lb

import (
	"load-balancer/server"
)

type LoadBalancerAlgo interface {
	FindServer() server.Server
}

type LoadBalancer interface {
	Start(port string)
	Stop()
	AddServer(backendServer server.Server)
	RemoveServer(serverName string)
	LoadBalancerAlgo
}
