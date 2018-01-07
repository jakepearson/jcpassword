package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

//shutdownProcess will print out to the console how long until the process is going down
//after the time has elapsed the process will end if `KillProcessOnShutdown` is enabled
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

//shutdownHandler will call the "graceful" shutdown method of the listener then run a background job
//that will kill the process after 5 seconds have elapsed
func shutdownHandler(webServer *WebServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		go shutdownProcess(webServer)
		webServer.Server.Shutdown(context.Background())
	}
}
