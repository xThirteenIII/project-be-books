package app

import (
	"books/internal/db"
	"books/internal/gutendex"
	"books/internal/handler"
	"books/internal/repo"
	"books/internal/rmq"
	"books/internal/service"
	"log"
	"net/http"
)

func Run() {
	db, err := db.Connect()
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

	queueName := "review_job"
	q, err := publishCh.QueueDeclare(
		queueName,
		true,  // durable
		false, // autodelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("can't declare publisher queue: %v\n", err)
	}

	reviewRepo := repo.NewMariaDBReviewRepository(db)
	publisher := rmq.NewPublisher(publishCh, queueName)
	var gutendexClient gutendex.GutendexClient
	reviewService := service.NewReviewService(reviewRepo, publisher, &gutendexClient)
	reviewHandler := handler.NewReviewHandler(reviewService)

	consumer := rmq.NewConsumer(consumerCh, queueName, reviewService)
	err = consumer.StartReviewConsumer()
	if err != nil {
		log.Printf("can't Consume on %s queue: %v\n", q.Name, err)
	}

	mux := http.NewServeMux()
	const port = "8010"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Use of singular instead of plural as written in the assignment.
	mux.HandleFunc("GET /book/search", handler.SearchBooksByKeywords)
	mux.HandleFunc("POST /review", reviewHandler.SubmitReview)
	mux.HandleFunc("GET /review/{id}", reviewHandler.GetReview)
	mux.HandleFunc("PUT /review/{id}", reviewHandler.UpdateReview)
	mux.HandleFunc("DELETE /review/{id}", reviewHandler.DeleteReview)

	// this blocks forever, until the server
	// has an unrecoverable error
	log.Println("server started on", port)
	err = server.ListenAndServe()
	log.Fatal(err)
}
