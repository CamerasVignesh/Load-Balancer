package main

import (
	"fmt"
	roundrobin "load-balancer/lb/round_robin"
	"os"
	"os/signal"
	"syscall"
)

// https://codingchallenges.fyi/challenges/challenge-load-balancer
func main() {
	lb := roundrobin.NewRoundRobinLoadBalancer(3)
	lb.Start("8080")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal (like SIGINT, Ctrl+C) to stop the load balancer
	select {
	case <-stopChan:
		// Graceful shutdown triggered, stop the load balancer
		fmt.Println("Received shutdown signal, stopping load balancer.")
		lb.Stop() // Gracefully stop the load balancer
	}
}
