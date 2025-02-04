package server

import (
	"fmt"
	"net/http"
	"time"
)

type Server interface {
	Start()
	HandleRequest(w http.ResponseWriter, r *http.Request)
	GetName() string
	GetActiveConnections() int
	Stop()
}

type serverImpl struct {
	name string
	port string
	connections int
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
	go func() {
		http.ListenAndServe(":" + s.port, mux)
	}()
	fmt.Println("server started")
}

func (s *serverImpl) Stop() {
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
	s.connections++
    fmt.Fprintln(w, "Request received to BE: " + s.port)
	time.Sleep(3*time.Second)
	s.connections--
}

func (s *serverImpl) GetActiveConnections() int {
	return s.connections
}

func (s *serverImpl) GetName() string {
	return s.name
}

