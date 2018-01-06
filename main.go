package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/jakepearson/jcpassword/handler"
)

func port() int {
	value, exists := os.LookupEnv("PORT") // Heroku port environment variable
	if !exists {
		return 8080 // Default port if environment variable not set
	}
	port, _ := strconv.Atoi(value)
	return port
}

func startServer() *handler.WebServer {
	webServer := handler.CreateServer(port(), true, 5)

	go func() {
		fmt.Printf("Listening at %s\n", webServer.Server.Addr)
		if err := webServer.Server.ListenAndServe(); err != nil {
			fmt.Printf("Httpserver: ListenAndServe() error: %s\n", err)
		}
	}()

	return webServer
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	startServer()
	wg.Wait()
}
