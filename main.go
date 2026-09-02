package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	m := initRoutes()

	startServer(m)
}

func initRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", wsHandler)

	return mux
}

func startServer(mux *http.ServeMux) {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Server listening on", server.Addr)
	log.Fatal(server.ListenAndServe())
}
