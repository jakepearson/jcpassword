package handler

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

//Statistics holds metrics information about the server
type Statistics struct {
	TotalRequests         uint64 `json:"total"`
	totalResponseTime     time.Duration
	AverageResponseTimeMS float32 `json:"average"`
}

//Hash contains a hash string and whether it is complete
type Hash struct {
	ID       int
	Complete bool
	Value    string
}

//WebServer contains all fields of the server to allow easier testing and multiple instances of the server if needed
type WebServer struct {
	Closed       bool
	Port         int
	ShutdownChan chan int
	SleepSeconds int
	nextID       int
	idLock       *sync.Mutex
	Server       *http.Server
	Handler      *http.Handler
	Hashes       *sync.Map
	Statistics   *Statistics
}

//createHandler sets up the routes for the web server to serve
func createHandler(webServer *WebServer) http.Handler {
	h := http.NewServeMux()

	h.HandleFunc("/hash", hashHandler(webServer))
	h.HandleFunc("/hash/", hashHandler(webServer))
	h.HandleFunc("/shutdown", shutdownHandler(webServer))
	h.HandleFunc("/stats", statsHandler(webServer))

	return h
}

// CreateServer will return an instance of a webserver (does not open a port)
func CreateServer(port int, sleepSeconds int) *WebServer {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port)}
	webServer := &WebServer{
		Closed:       false,
		SleepSeconds: sleepSeconds,
		Port:         port,
		Server:       server,
		idLock:       &sync.Mutex{},
		Hashes:       new(sync.Map),
		Statistics:   &Statistics{},
		ShutdownChan: make(chan int),
	}
	handler := createHandler(webServer)
	webServer.Handler = &handler
	webServer.Server.Handler = handler
	return webServer
}
