package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jakepearson/jcpassword/encoder"
)

func hashHandler(webServer *WebServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		time.Sleep(time.Duration(webServer.SleepSeconds) * time.Second) //Slow down request to meet requirement

		fmt.Fprintf(w, encoder.HashAndEncode(password))
	}
}
