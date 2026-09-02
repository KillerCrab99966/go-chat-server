package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	h := newHub()

	m := initRoutes(h)

	startServer(m)
}

func initRoutes(h *hub) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", h.wsHandler)

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
