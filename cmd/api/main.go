package main

import (
	"books/internal/db"
	rabbitmq "books/internal/queue"
	"log"
	"net/http"
)

func main() {
	err := db.Connect()
	if err != nil {
		log.Printf("Error while connecting to MariaDB: %v\nClosing app\n", err)
		return
	}

	err = rabbitmq.AttemptConnect()
	if err != nil {
		log.Printf("Error while connecting to RabbitMQ: %v\nClosing app\n", err)
		return
	}

	mux := http.NewServeMux()
	const port = "8010"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// this blocks forever, until the server
	// has an unrecoverable error
	log.Println("server started on", port)
	err = server.ListenAndServe()
	log.Fatal(err)
}
