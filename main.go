package main

import (
	"fmt"
	"net/http"

	"github.com/jakepearson/jcpassword/encoder"
)

func hashHandler(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("password")
	fmt.Fprintf(w, encoder.HashAndEncode(password))
}

func createHandler() http.Handler {
	h := http.NewServeMux()

	h.HandleFunc("/hash", hashHandler)

	return h
}

func main() {
	http.ListenAndServe(":8080", createHandler())
}
