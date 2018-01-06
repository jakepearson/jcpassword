package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/jakepearson/jcpassword/handler"
)

func port() int {
	value, exists := os.LookupEnv("PORT")
	if !exists {
		return 8080
	}
	port, _ := strconv.Atoi(value)
	return port
}

func startServer() {
	webServer := handler.CreateServer(port(), true, 5)

	go func() {
		fmt.Printf("Listening at %s\n", webServer.Server.Addr)
		if err := webServer.Server.ListenAndServe(); err != nil {
			fmt.Printf("Httpserver: ListenAndServe() error: %s\n", err)
		}
	}()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	startServer()
	wg.Wait()
}
