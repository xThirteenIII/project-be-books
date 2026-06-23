package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	const port = "8010"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// this blocks forever, until the server
	// has an unrecoverable error
	fmt.Println("server started on", port)
	err := server.ListenAndServe()
	log.Fatal(err)
}
