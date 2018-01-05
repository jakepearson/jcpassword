package main

import (
	"fmt"
	"net/http"
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

func startServer() *http.Server {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port()),
	}

	server.Handler = handler.Create(server)

	go func() {
		fmt.Printf("Listening at %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Httpserver: ListenAndServe() error: %s\n", err)
		}
	}()

	return server
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	startServer()
	wg.Wait()
	fmt.Println("Falling off the end")
}
