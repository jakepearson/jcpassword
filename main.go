package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

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

func main() {
	http.ListenAndServe(fmt.Sprintf(":%d", port()), handler.Create())
}
