package service

import (
	"books/internal/gutendex"
	"books/internal/model"
	"books/internal/repo"
	"books/internal/rmq"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	max_chars = 2048
	max_score = 10
	min_score = 0
)

// ReviewService depends on the concrete ReviewRepo and ReviewPublisher interface.
type ReviewService struct {
	repo      repo.ReviewRepo
	publisher rmq.ReviewPublisher
	books     BookLookup
}

type BookLookup interface {
	// We accept the redundancy of the method returning a book both in SubmitReview and EnrichReview.
	// Creating two different methods seems an overkill.
	SearchByID(id int) (gutendex.Book, error)
}

func NewReviewService(r repo.ReviewRepo, p rmq.ReviewPublisher, b BookLookup) *ReviewService {
	return &ReviewService{
		repo:      r,
		publisher: p,
		books:     b,
	}
}

// SubmitReview service handles business logic.
// Validates the request parameters.
// It returns a unique review id and error.
func (s *ReviewService) SubmitReview(ctx context.Context, req model.ReviewCreateRequest) (string, error) {

	// Validate review text length and score range.
	if len(req.Review) > max_chars {
		return "", fmt.Errorf("review exceeds the max length (2048): %d\n", len(req.Review))
	}

	if req.Score > max_score || req.Score < min_score {
		return "", fmt.Errorf("score must be within 0 and 10\n")
	}

	// Check if book id exists on the Gutendex API.
	_, err := s.books.SearchByID(req.BookID)
	if err != nil {
		return "", err
	}

	// Create new unique ID for review.
	reviewID := uuid.NewString()
	now := time.Now()

	// Review will persist in MariaDB.
	// We save a first, not complete data, with status PENDING.
	// It will be enriched with metadata asyncly.
	review := model.Review{
		ID:         reviewID,
		BookID:     req.BookID,
		ReviewText: req.Review,
		Score:      req.Score,
		UpdatedAt:  now,
		CreatedAt:  now,
		Status:     model.StatusPending,
	}

	if err := s.repo.Insert(ctx, review); err != nil {
		return "", err
	}

	if err := s.publisher.PublishReviewJob(ctx, review.ID, review.BookID); err != nil {

		// If the publish fails, update review status to Failed.
		// This avoids zombie reviews with status "pending" in DB.
		review.Status = model.StatusFailed
		if updateErr := s.repo.Update(ctx, review); updateErr != nil {
			return "", fmt.Errorf("publish review Job: %w; mark failed: %v", err, updateErr)
		}
		return "", err
	}

	return reviewID, nil
}

func (s *ReviewService) EnrichReview(ctx context.Context, reviewID string, bookID int) error {
	book, err := s.books.SearchByID(bookID)
	if err != nil {
		// opzionale: update status failed
		return err
	}

	// Collect authors name in a slice
	authorSlice := []string{}
	for _, author := range book.Authors {
		authorSlice = append(authorSlice, author.Name)

	}
	// Separate them by , and join them into a single string
	authors := strings.Join(authorSlice, ", ")

	review := model.Review{
		ID:        reviewID,
		BookID:    bookID,
		BookTitle: book.Title,
		Authors:   authors,
		CoverURL:  book.Formats.ImageJpeg,
		Status:    model.StatusReady,
		UpdatedAt: time.Now(),
	}

	return s.repo.Update(ctx, review)
}
