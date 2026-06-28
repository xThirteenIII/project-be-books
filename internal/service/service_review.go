package service

import (
	"books/internal/gutendex"
	"books/internal/model"
	"books/internal/repo"
	"books/internal/rmq"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	maxChars = 2048
	maxScore = 10
	minScore = 0
)

var ErrValidation = errors.New("validation error")

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
	if req.BookID <= 0 {
		return "", fmt.Errorf("%w: book id must be positive", ErrValidation)
	}

	if strings.TrimSpace(req.Review) == "" {
		return "", fmt.Errorf("%w: review must not be empty", ErrValidation)
	}

	if len(req.Review) > maxChars {
		return "", fmt.Errorf("%w: review exceeds the max length (2048): %d", ErrValidation, len(req.Review))
	}

	if req.Score > maxScore || req.Score < minScore {
		return "", fmt.Errorf("%w: score must be within 0 and 10", ErrValidation)
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
		if updateErr := s.repo.Update(ctx, model.Review{
			ID:        reviewID,
			BookID:    bookID,
			Status:    model.StatusFailed,
			UpdatedAt: time.Now(),
		}); updateErr != nil {
			return fmt.Errorf("lookup book: %w; mark failed: %v", err, updateErr)
		}
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

func (s *ReviewService) GetReview(ctx context.Context, reviewID string) (model.Review, error) {
	return s.repo.GetByID(ctx, reviewID)
}

func (s *ReviewService) UpdateReview(ctx context.Context, reviewID string, req model.ReviewUpdateRequest) error {
	if strings.TrimSpace(req.Review) == "" {
		return fmt.Errorf("%w: review must not be empty", ErrValidation)
	}

	if len(req.Review) > maxChars {
		return fmt.Errorf("%w: review exceeds the max length (2048): %d", ErrValidation, len(req.Review))
	}

	if req.Score > maxScore || req.Score < minScore {
		return fmt.Errorf("%w: score must be within 0 and 10", ErrValidation)
	}

	return s.repo.UpdateReviewContent(ctx, reviewID, req.Review, req.Score)
}

func (s *ReviewService) DeleteReview(ctx context.Context, reviewID string) error {
	return s.repo.Delete(ctx, reviewID)
}
