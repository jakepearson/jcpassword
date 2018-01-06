package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

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
