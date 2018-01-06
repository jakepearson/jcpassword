package handler

import (
	"fmt"
	"net/http"
)

//Statistics holds metrics information about the server
type Statistics struct {
	TotalRequests       uint64
	totalResponseMS     uint64
	AverageResponseTime uint64
}

//Hash contains the a hash string and whether it is complete
type Hash struct {
	ID       int
	Complete bool
	Value    string
}

//WebServer contains all fields of the server to allow easier testing and multiple instances of the server if needed
type WebServer struct {
	Closed                bool
	Port                  int
	KillProcessOnShutdown bool
	SleepSeconds          int
	Server                *http.Server
	Handler               *http.Handler
	Hashes                map[int]*Hash
	Statistics            *Statistics
}

func slashHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func createHandler(webServer *WebServer) http.Handler {
	h := http.NewServeMux()

	h.HandleFunc("/hash", hashHandler(webServer))
	h.HandleFunc("/hash/", hashHandler(webServer))
	h.HandleFunc("/shutdown", shutdownHandler(webServer))
	h.HandleFunc("/stats", statsHandler(webServer))

	return h
}

// CreateServer will return an instance of a webserver (does not open a port)
func CreateServer(port int, killProcessOnShutdown bool, sleepSeconds int) *WebServer {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port)}
	webServer := &WebServer{
		Closed:                false,
		KillProcessOnShutdown: killProcessOnShutdown,
		SleepSeconds:          sleepSeconds,
		Port:                  port,
		Server:                server,
		Hashes:                make(map[int]*Hash),
		Statistics:            &Statistics{},
	}
	handler := createHandler(webServer)
	webServer.Handler = &handler
	webServer.Server.Handler = handler
	return webServer
}
