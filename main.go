package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jakepearson/jcpassword/handler"
)

//port will get the port to use from an environment variable.
//if the environent is not set the port will be `8080``
func port() int {
	value, exists := os.LookupEnv("PORT") // Heroku port environment variable
	if !exists {
		return 8080 // Default port if environment variable not set
	}
	port, _ := strconv.Atoi(value)
	return port
}

//startServer will create an instance of the `webServer` object then
//open a listener in the background
func startServer() *handler.WebServer {
	webServer := handler.CreateServer(port(), 5)

	go func() {
		fmt.Printf("Listening at %s\n", webServer.Server.Addr)
		if err := webServer.Server.ListenAndServe(); err != nil {
			fmt.Printf("Httpserver: ListenAndServe() error: %s\n", err)
		}
	}()

	return webServer
}

//main will create a wait group, start the server
//then wait for the waitgroup to complete (which will never happen)
func main() {
	webServer := startServer()
	<-webServer.ShutdownChan // Exit process when a message shows on the channel
}
