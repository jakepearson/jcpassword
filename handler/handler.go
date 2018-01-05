package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jakepearson/jcpassword/encoder"
)

func hashHandler(w http.ResponseWriter, r *http.Request) {
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

	time.Sleep(5 * time.Second) //Slow down request to meet requirement

	fmt.Fprintf(w, encoder.HashAndEncode(password))
}

func shutdownProcess(server *http.Server) {
	for i := 0; i < 5; i++ {
		fmt.Printf("Shutting down service in %d seconds\n", 5-i)
		time.Sleep(1 * time.Second)
	}
	os.Exit(0)
}

func shutdownHandler(server *http.Server) func(w http.ResponseWriter, r *http.Request) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		go shutdownProcess(server)
		server.Shutdown(context.Background())
	}
	return handler
}

// Create a ServeMux handler will all routes and return
func Create(server *http.Server) http.Handler {
	h := http.NewServeMux()

	h.HandleFunc("/hash", hashHandler)
	h.HandleFunc("/shutdown", shutdownHandler(server))

	return h
}
