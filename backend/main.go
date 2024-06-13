package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request from %s\n", r.RemoteAddr)
	fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	port := "127.0.0.1:8080"
	http.HandleFunc("/", requestHandler)
	http.HandleFunc("/health", healthHandler)
	for {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			port = findAvailablePort(port)
		} else {
			break
		}
	}
}

func findAvailablePort(port string) string {
	portNumber, err := strconv.Atoi(strings.TrimPrefix(port, "127.0.0.1:"))
	if err != nil {
		panic("invalid port number!")
	}
	portNumber++
	return "127.0.0.1:" + strconv.Itoa(portNumber)
}
