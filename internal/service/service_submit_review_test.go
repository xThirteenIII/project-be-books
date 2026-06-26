package service_test

import (
	"books/internal/gutendex"
	"books/internal/mock"
	"books/internal/model"
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
			wantErr: true,
		},
		{
			name: "score negative",
			req: model.ReviewCreateRequest{
				BookID: 1,
				Review: "It sucks",
				Score:  -1,
			},
			wantErr: true,
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
			wantErr: true,
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
