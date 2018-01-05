package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jakepearson/jcpassword/encoder"
)

func hashHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second) //Slow down request to meet requirement

	password := r.URL.Query().Get("password")
	fmt.Fprintf(w, encoder.HashAndEncode(password))
}

// Create a ServeMux handler will all routes and return
func Create() http.Handler {
	h := http.NewServeMux()

	h.HandleFunc("/hash", hashHandler)

	return h
}
