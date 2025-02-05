package lb

import (
	"load-balancer/server"
)

type LoadBalancerAlgo interface {
	FindServer() server.Server
	PerformHealthChecks()
}

type LoadBalancer interface {
	Start(port string)
	Stop()
	AddServer(backendServer server.Server)
	RemoveServer(serverName string)
	LoadBalancerAlgo
}

