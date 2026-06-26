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
		GetByID(ctx context.Context, id string) error
		Delete(ctx context.Context, id string) error
	}
)
