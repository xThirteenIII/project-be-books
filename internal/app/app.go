package app

import (
	"books/internal/db"
	"books/internal/handler"
	"books/internal/rmq"
	"log"
	"net/http"
)

func Run() {
	_, err := db.Connect()
	if err != nil {
		log.Printf("Error while connecting to MariaDB: %v\nClosing app\n", err)
		return
	}

	rabbitConn, err := rmq.AttemptConnect()
	if err != nil {
		log.Printf("Error while connecting to RabbitMQ: %v\nClosing app\n", err)
		return
	}
	defer rabbitConn.Close()

	publishCh, err := rabbitConn.Channel()
	if err != nil {
		// Let the program shut down if we can't create a channel.
		log.Fatalf("could not create publisher channel: %v\n", err)
	}
	consumerCh, err := rabbitConn.Channel()
	if err != nil {
		// Let the program shut down if we can't create a channel.
		log.Fatalf("could not create consumer channel: %v\n", err)
	}

	q, err := publishCh.QueueDeclare(
		"review_enrichment",
		true,  // durable
		false, // autodelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		log.Printf("can't declare publisher queue: %v\n", err)
	}

	err = rmq.PublishReviewJob(publishCh, q.Name, rmq.ReviewJob{
		ReviewID: "100",
		BookID:   10,
	})
	if err != nil {
		log.Printf("can't Publish on %s queue: %v\n", q.Name, err)
	}

	err = rmq.StartReviewConsumer(consumerCh, q.Name)
	if err != nil {
		log.Printf("can't Consume on %s queue: %v\n", q.Name, err)
	}

	mux := http.NewServeMux()
	const port = "8010"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("GET /books/search", handler.SearchBooksByKeywords)
	mux.HandleFunc("POST /review", handler.SubmitReview)

	// this blocks forever, until the server
	// has an unrecoverable error
	log.Println("server started on", port)
	err = server.ListenAndServe()
	log.Fatal(err)
}
