package service_test

import (
	"books/internal/gutendex"
	"books/internal/mock"
	"books/internal/model"
	"books/internal/repo"
	"books/internal/service"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubmitReview(t *testing.T) {
	tests := []struct {
		name          string
		req           model.ReviewCreateRequest
		repoInsertErr error
		publishErr    error
		bookLookupErr error
		wantErr       bool
		wantErrIs     error
	}{
		{
			name: "valid review",
			req: model.ReviewCreateRequest{
				BookID: 1,
				Review: "Nice book, very entertaining.",
				Score:  8,
			},
			wantErr: false,
		},
		{
			name: "score too high",
			req: model.ReviewCreateRequest{
				BookID: 1,
				Review: "Wonderful",
				Score:  11,
			},
			wantErr:   true,
			wantErrIs: service.ErrValidation,
		},
		{
			name: "score negative",
			req: model.ReviewCreateRequest{
				BookID: 1,
				Review: "It sucks",
				Score:  -1,
			},
			wantErr:   true,
			wantErrIs: service.ErrValidation,
		},
		{
			name: "book id missing",
			req: model.ReviewCreateRequest{
				Review: "Wonderful",
				Score:  8,
			},
			wantErr:   true,
			wantErrIs: service.ErrValidation,
		},
		{
			name: "review empty",
			req: model.ReviewCreateRequest{
				BookID: 1,
				Review: "   ",
				Score:  8,
			},
			wantErr:   true,
			wantErrIs: service.ErrValidation,
		},
		{
			name: "score zero is valid",
			req: model.ReviewCreateRequest{
				BookID: 1,
				Review: "Did not like it at all.",
				Score:  0,
			},
			wantErr: false,
		},
		{
			name: "score ten is valid",
			req: model.ReviewCreateRequest{
				BookID: 1,
				Review: "Masterpiece.",
				Score:  10,
			},
			wantErr: false,
		},
		{
			name: "review too long",
			req: model.ReviewCreateRequest{
				BookID: 1,
				Review: strings.Repeat("a", 2049),
				Score:  5,
			},
			wantErr:   true,
			wantErrIs: service.ErrValidation,
		},
		{
			name: "review at max length is valid",
			req: model.ReviewCreateRequest{
				BookID: 1,
				Review: strings.Repeat("a", 2048),
				Score:  5,
			},
			wantErr: false,
		},
		{
			name: "book id not found on gutendex",
			req: model.ReviewCreateRequest{
				BookID: 99999999,
				Review: "Nice cover.",
				Score:  5,
			},
			bookLookupErr: gutendex.ErrNoMatch,
			wantErr:       true,
			wantErrIs:     gutendex.ErrNoMatch,
		},
		{
			name: "repo insert error",
			req: model.ReviewCreateRequest{
				BookID: 1,
				Review: "Ok.",
				Score:  5,
			},
			repoInsertErr: errors.New("db down"),
			wantErr:       true,
		},
		{
			name: "publisher error",
			req: model.ReviewCreateRequest{
				BookID: 1,
				Review: "Ok.",
				Score:  5,
			},
			publishErr: errors.New("rabbitmq down"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockReviewRepo{
				InsertFn: func(ctx context.Context, r model.Review) error {
					return tt.repoInsertErr
				},
				UpdateFn: func(ctx context.Context, r model.Review) error {
					return nil
				},
			}
			mockPub := &mock.MockPublisher{
				PublishFn: func(ctx context.Context, reviewID string, bookID int) error {
					return tt.publishErr
				},
			}
			mockBooks := &mock.MockBookLookup{
				SearchByIDFn: func(id int) (gutendex.Book, error) {
					return gutendex.Book{}, tt.bookLookupErr
				},
			}

			svc := service.NewReviewService(mockRepo, mockPub, mockBooks)
			reviewID, err := svc.SubmitReview(context.Background(), tt.req)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
				assert.Empty(t, reviewID)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, reviewID)
			}
		})
	}
}

func TestGetReview(t *testing.T) {
	wantReview := model.Review{
		ID:         "review-id",
		BookID:     2701,
		ReviewText: "Good book.",
		Score:      8,
		Status:     model.StatusReady,
		BookTitle:  "Moby Dick",
		Authors:    "Herman Melville",
		CoverURL:   "https://example.com/cover.jpg",
	}

	mockRepo := &mock.MockReviewRepo{
		GetByIDFn: func(ctx context.Context, id string) (model.Review, error) {
			assert.Equal(t, "review-id", id)
			return wantReview, nil
		},
	}

	svc := service.NewReviewService(mockRepo, &mock.MockPublisher{}, &mock.MockBookLookup{})

	gotReview, err := svc.GetReview(context.Background(), "review-id")

	require.NoError(t, err)
	assert.Equal(t, wantReview, gotReview)
}

func TestGetReviewReturnsRepoError(t *testing.T) {
	mockRepo := &mock.MockReviewRepo{
		GetByIDFn: func(ctx context.Context, id string) (model.Review, error) {
			return model.Review{}, repo.ErrNotFound
		},
	}

	svc := service.NewReviewService(mockRepo, &mock.MockPublisher{}, &mock.MockBookLookup{})

	gotReview, err := svc.GetReview(context.Background(), "missing-id")

	require.ErrorIs(t, err, repo.ErrNotFound)
	assert.Empty(t, gotReview)
}

func TestUpdateReview(t *testing.T) {
	tests := []struct {
		name        string
		req         model.ReviewUpdateRequest
		repoErr     error
		wantErr     bool
		wantRepoHit bool
	}{
		{
			name: "valid update",
			req: model.ReviewUpdateRequest{
				Review: "Updated review.",
				Score:  7,
			},
			wantRepoHit: true,
		},
		{
			name: "score too high",
			req: model.ReviewUpdateRequest{
				Review: "Updated review.",
				Score:  11,
			},
			wantErr: true,
		},
		{
			name: "score negative",
			req: model.ReviewUpdateRequest{
				Review: "Updated review.",
				Score:  -1,
			},
			wantErr: true,
		},
		{
			name: "review too long",
			req: model.ReviewUpdateRequest{
				Review: strings.Repeat("a", 2049),
				Score:  5,
			},
			wantErr: true,
		},
		{
			name: "review empty",
			req: model.ReviewUpdateRequest{
				Review: "   ",
				Score:  5,
			},
			wantErr: true,
		},
		{
			name: "repo error",
			req: model.ReviewUpdateRequest{
				Review: "Updated review.",
				Score:  7,
			},
			repoErr:     repo.ErrNotFound,
			wantErr:     true,
			wantRepoHit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoHit := false
			mockRepo := &mock.MockReviewRepo{
				UpdateReviewContentFn: func(ctx context.Context, id string, reviewText string, score int) error {
					repoHit = true
					assert.Equal(t, "review-id", id)
					assert.Equal(t, tt.req.Review, reviewText)
					assert.Equal(t, tt.req.Score, score)
					return tt.repoErr
				},
			}

			svc := service.NewReviewService(mockRepo, &mock.MockPublisher{}, &mock.MockBookLookup{})

			err := svc.UpdateReview(context.Background(), "review-id", tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.wantRepoHit, repoHit)
		})
	}
}

func TestEnrichReviewMarksReviewFailedWhenBookLookupFails(t *testing.T) {
	mockRepo := &mock.MockReviewRepo{
		UpdateFn: func(ctx context.Context, r model.Review) error {
			assert.Equal(t, "review-id", r.ID)
			assert.Equal(t, 2701, r.BookID)
			assert.Equal(t, model.StatusFailed, r.Status)
			return nil
		},
	}
	mockBooks := &mock.MockBookLookup{
		SearchByIDFn: func(id int) (gutendex.Book, error) {
			return gutendex.Book{}, gutendex.ErrNoMatch
		},
	}

	svc := service.NewReviewService(mockRepo, &mock.MockPublisher{}, mockBooks)

	err := svc.EnrichReview(context.Background(), "review-id", 2701)

	require.ErrorIs(t, err, gutendex.ErrNoMatch)
}

func TestDeleteReview(t *testing.T) {
	mockRepo := &mock.MockReviewRepo{
		DeleteFn: func(ctx context.Context, id string) error {
			assert.Equal(t, "review-id", id)
			return nil
		},
	}

	svc := service.NewReviewService(mockRepo, &mock.MockPublisher{}, &mock.MockBookLookup{})

	err := svc.DeleteReview(context.Background(), "review-id")

	require.NoError(t, err)
}

func TestDeleteReviewReturnsRepoError(t *testing.T) {
	mockRepo := &mock.MockReviewRepo{
		DeleteFn: func(ctx context.Context, id string) error {
			return repo.ErrNotFound
		},
	}

	svc := service.NewReviewService(mockRepo, &mock.MockPublisher{}, &mock.MockBookLookup{})

	err := svc.DeleteReview(context.Background(), "missing-id")

	require.ErrorIs(t, err, repo.ErrNotFound)
}
