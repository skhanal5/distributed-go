package main

import (
	"log"

	"github.com/skhanal5/distributed-go/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}