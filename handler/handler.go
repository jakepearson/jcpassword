package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jakepearson/jcpassword/encoder"
)

//WebServer contains all fields of the server to allow easier testing and multiple instances of the server if needed
type WebServer struct {
	Closed                bool
	Port                  int
	KillProcessOnShutdown bool
	SleepSeconds          int
	Server                *http.Server
	Handler               *http.Handler
}

func hashHandler(webServer *WebServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		password := r.URL.Query().Get("password")
		if password == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Password missing (use password query parameter)")
			return
		}

		time.Sleep(time.Duration(webServer.SleepSeconds) * time.Second) //Slow down request to meet requirement

		fmt.Fprintf(w, encoder.HashAndEncode(password))
	}
}

func shutdownProcess(server *WebServer) {
	for i := 0; i < server.SleepSeconds; i++ {
		fmt.Printf("Shutting down service in %d seconds\n", server.SleepSeconds-i)
		time.Sleep(1 * time.Second)
	}

	server.Closed = true
	if server.KillProcessOnShutdown {
		os.Exit(0)
	}
}

func shutdownHandler(webServer *WebServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		go shutdownProcess(webServer)
		webServer.Server.Shutdown(context.Background())
	}
}

func createHandler(webServer *WebServer) http.Handler {
	h := http.NewServeMux()

	h.HandleFunc("/hash", hashHandler(webServer))
	h.HandleFunc("/shutdown", shutdownHandler(webServer))

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
	}
	handler := createHandler(webServer)
	webServer.Handler = &handler
	return webServer
}
