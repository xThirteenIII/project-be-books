// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"books/internal/model"
	"context"
)

type (
	ReviewRepo interface {
		Insert(ctx context.Context, r model.Review) error
		Update(ctx context.Context, r model.Review) error
		UpdateReviewContent(ctx context.Context, id string, reviewText string, score int) error
		GetByID(ctx context.Context, id string) (model.Review, error)
		Delete(ctx context.Context, id string) error
	}
)
