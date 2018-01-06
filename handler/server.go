package handler

import (
	"fmt"
	"net/http"
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

func slashHandler(w http.ResponseWriter, r *http.Request) {
	html := `<html>
	<head>
	<title>Password Hasher</title>
	</head>
	<body>
	<img src="https://azure.microsoft.com/svghandler/security-center?width=600&height=315">
	</body>
	</html>`
	fmt.Fprintf(w, html)
}

func createHandler(webServer *WebServer) http.Handler {
	h := http.NewServeMux()

	h.HandleFunc("/", slashHandler)
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
	webServer.Server.Handler = handler
	return webServer
}
