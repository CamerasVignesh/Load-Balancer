package server

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Server interface {
	Start()
	GetName() string
	HandleRequest(w http.ResponseWriter, r *http.Request)
	GetActiveConnections() int
	Stop()
}

type serverImpl struct {
	name string
	port string
	connections int
	mu sync.RWMutex
	wg sync.WaitGroup
}

func NewServer(serverID string, port string) Server {
	s := &serverImpl{
		name: serverID,
		port: port,
	}
	return s
}

func (s *serverImpl) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HandleRequest)
	mux.HandleFunc("/health", s.healthCheck)
	go func() {
		http.ListenAndServe(":" + s.port, mux)
	}()
	fmt.Println("server started")
}

func (s *serverImpl) healthCheck(w http.ResponseWriter, r *http.Request) {
	if val, _ := strconv.Atoi(s.port); val % 2 == 0 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Server on port: " + s.port + " is healthy")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Server on port: " + s.port + " is unhealthy")
	}
}

func (s *serverImpl) Stop() {
	s.wg.Wait()
	fmt.Println("server: " + s.name + " stopped")
}


func (s *serverImpl) HandleRequest(w http.ResponseWriter, r *http.Request) {
	/*
	fmt.Println("./be")
    fmt.Printf("Received request from %s\n", r.RemoteAddr)
    fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
    fmt.Printf("Host: %s\n", r.Host)
    fmt.Printf("User-Agent: %s\n", r.Header.Get("User-Agent"))
    fmt.Printf("Accept: %s\n", r.Header.Get("Accept"))
	*/
	fmt.Fprintln(w, "Request received to BE: " + s.port)
	s.increaseActiveConnections()
    defer s.decreaseActiveConnections()
	processRequest()
}

func (s *serverImpl) GetActiveConnections() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	connections := s.connections
	return connections
}

func (s *serverImpl) GetName() string {
	return s.name
}

func (s *serverImpl) increaseActiveConnections() {
	s.mu.Lock()
	s.wg.Add(1)
	defer s.mu.Unlock()
	s.connections++
}

func (s *serverImpl) decreaseActiveConnections() {
	s.mu.Lock()
	s.wg.Done()
	defer s.mu.Unlock()
	s.connections--
}

func processRequest() {
	time.Sleep(3*time.Second)
}
