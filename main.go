package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

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

func port() int {
	value, exists := os.LookupEnv("PORT")
	if !exists {
		return 8080
	}
	port, _ := strconv.Atoi(value)
	return port
}

func main() {
	http.ListenAndServe(fmt.Sprintf(":%d", port()), createHandler())
}
