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
- [Notes](#notes)


## Requirements

- Docker and Docker Compose
- Go 1.24, only if you want to run tests or the API outside Docker
- curl, only if you want to run tests on the endpoints

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

Check if the 'reviews' table has been created  

```bash
docker compose exec -T db mariadb -u user -ppassword books -e "SHOW TABLES;"
```

## Unit tests

The project includes unit tests for the service layer.

They validate:
- review input validation
- book existence checks against the mocked Gutendex dependency
- repository and publisher error propagation
- review retrieval, update and deletion
- async enrichment failure handling

Run the tests with:

```bash
go test ./...
```

## Endpoints

### Search books

```bash
curl "http://localhost:8010/book/search?q=moby%20dick"
```

### Submit review

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

### Get review status

```bash
curl "http://localhost:8010/review/{review_id}"
```

The response is:

- `202 Accepted` while the review is still pending
- `200 OK` when the review has been enriched with book metadata

### Update review

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

```bash
curl -X DELETE http://localhost:8010/review/{review_id}
```

Returns `204 No Content`.



- `POST /review` validates the book id against Gutendex before storing the review.
- Reviews are first saved with `pending` status.
- A RabbitMQ consumer enriches reviews asynchronously with title, authors and cover URL.
- If async enrichment fails, the review is marked as `failed` instead of staying pending forever.
