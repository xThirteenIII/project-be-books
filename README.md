# Book review service

Small Go HTTP service for searching books on Gutendex and managing async book reviews.

## Table of Contents

- [Requirements](#requirements)
- [Run](#run)
- [Unit tests](#unit-tests)
- [Endpoints](#endpoints)
  - [Search books](#search-books)
  - [Submit review](#submit-review)
  - [Get review status](#get-review-status)
  - [Update review](#update-review)
  - [Delete review](#delete-review)
- [Troubleshooting](#troubleshooting)

## Requirements

- Docker and Docker Compose
- Go 1.24, only if you want to run tests or the API outside Docker
- curl, only if you want to test the endpoints manually

## Run

```bash
docker compose up --build
```

The API listens on `http://localhost:8010`.

The compose stack starts:

- MariaDB on port `3306`
- RabbitMQ on port `5672`
- RabbitMQ management UI on `http://localhost:15672` with `user/password`
- the Go API on port `8010`

The database schema is loaded automatically from `migration/001_create_reviews_table.up.sql` when the MariaDB container initializes.

If you reuse an old database container created before this migration was mounted, recreate it with:

```bash
docker compose down -v
docker compose up --build
```

To verify that the `reviews` table has been created:

```bash
docker compose exec -T db mariadb -u user -ppassword books -e "SHOW TABLES;"
```

To later verify table contents:

```bash
docker compose exec -T db mariadb -u user -ppassword books -e "SELECT * FROM reviews;"
```

## Unit tests

The project includes unit tests for the service layer.

They cover:

- review input validation
- book existence checks against a mocked Gutendex dependency
- repository and publisher error propagation
- review retrieval, update, and deletion
- async enrichment failure handling

Run the tests with:

```bash
go test ./...
```

## Endpoints

### Search books

`GET /book/search?q={keywords}`

```bash
curl "http://localhost:8010/book/search?q=moby%20dick"
```

### Submit review

`POST /review`

```bash
curl -X POST http://localhost:8010/review \
  -H "Content-Type: application/json" \
  -d '{
    "id": 2701,
    "review": "Bel libro, molto coinvolgente.",
    "score": 8
  }'
```

The response is `202 Accepted` and returns a generated `review_id`.

Additional notes:

- `POST /review` validates the book id against Gutendex before storing the review.
- Reviews are first saved with `pending` status.
- A RabbitMQ consumer enriches reviews asynchronously with title, authors, and cover URL.
- If async enrichment fails, the review is marked as `failed` instead of remaining `pending` indefinitely.

To verify this behaviour end-to-end I put a time.Sleep(10 * time.Second) before enriching the review in `internal/rmq/consumer.go`, `line 57`.  
The line is now commented.  
This gives enough time to check DB content before the review goes from `pending` to `ready`.

### Get review status

`GET /review/{review_id}`

```bash
curl "http://localhost:8010/review/{review_id}"
```

The response is:

- `202 Accepted` while the review is still pending
- `200 OK` when the review has been enriched with book metadata

### Update review

`PUT /review/{review_id}`

```bash
curl -X PUT http://localhost:8010/review/{review_id} \
  -H "Content-Type: application/json" \
  -d '{
    "review": "Updated review text.",
    "score": 9
  }'
```

Returns `204 No Content`.

### Delete review

`DELETE /review/{review_id}`

```bash
curl -X DELETE http://localhost:8010/review/{review_id}
```

Returns `204 No Content`.

## Troubleshooting

On some work laptops, corporate VPNs or proxies may intercept HTTPS traffic and prevent the application from reaching the Gutendex API.

In this case, calls such as `GET /book/search?q={keywords}` or `POST /review` may fail with TLS or x509 certificate errors.

If this happens:

1. Disable the VPN and try again.
2. If the problem persists, run only MariaDB and RabbitMQ with Docker Compose.
3. Start the Go API separately from your local machine.

Start MariaDB and RabbitMQ:

```bash
docker compose up db rabbitmq
```

Then, in a separate terminal:

```bash
go run cmd/api/main.go
```
