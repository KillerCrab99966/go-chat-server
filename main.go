package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	h := newHub()

func main() {
	rm := newRoom()

	m := initRoutes(rm)

	startServer(m)
}

func initRoutes(rm *room) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", rm.wsHandler)

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
