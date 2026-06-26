package mock

import (
	"books/internal/gutendex"
	"books/internal/model"
	"context"
)

type MockReviewRepo struct {
	InsertFn func(ctx context.Context, r model.Review) error
	UpdateFn func(ctx context.Context, r model.Review) error
}

func (m *MockReviewRepo) Insert(ctx context.Context, r model.Review) error {
	return m.InsertFn(ctx, r)
}
func (m *MockReviewRepo) Update(ctx context.Context, r model.Review) error {
	return m.UpdateFn(ctx, r)
}
func (m *MockReviewRepo) GetByID(ctx context.Context, id string) error { return nil }
func (m *MockReviewRepo) Delete(ctx context.Context, id string) error  { return nil }

type MockPublisher struct {
	PublishFn func(ctx context.Context, reviewID string, bookID int) error
}

func (m *MockPublisher) PublishReviewJob(ctx context.Context, reviewID string, bookID int) error {
	return m.PublishFn(ctx, reviewID, bookID)
}

type MockBookLookup struct {
	SearchByIDFn func(id int) (gutendex.Book, error)
}

func (m *MockBookLookup) SearchByID(id int) (gutendex.Book, error) {
	return m.SearchByIDFn(id)
}
