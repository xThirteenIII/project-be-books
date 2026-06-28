package mock

import (
	"books/internal/gutendex"
	"books/internal/model"
	"context"
)

type MockReviewRepo struct {
	InsertFn              func(ctx context.Context, r model.Review) error
	UpdateFn              func(ctx context.Context, r model.Review) error
	UpdateReviewContentFn func(ctx context.Context, id string, reviewText string, score int) error
	GetByIDFn             func(ctx context.Context, id string) (model.Review, error)
	DeleteFn              func(ctx context.Context, id string) error
}

func (m *MockReviewRepo) Insert(ctx context.Context, r model.Review) error {
	return m.InsertFn(ctx, r)
}
func (m *MockReviewRepo) Update(ctx context.Context, r model.Review) error {
	return m.UpdateFn(ctx, r)
}
func (m *MockReviewRepo) UpdateReviewContent(ctx context.Context, id string, reviewText string, score int) error {
	if m.UpdateReviewContentFn == nil {
		return nil
	}
	return m.UpdateReviewContentFn(ctx, id, reviewText, score)
}
func (m *MockReviewRepo) GetByID(ctx context.Context, id string) (model.Review, error) {
	if m.GetByIDFn == nil {
		return model.Review{}, nil
	}
	return m.GetByIDFn(ctx, id)
}
func (m *MockReviewRepo) Delete(ctx context.Context, id string) error {
	if m.DeleteFn == nil {
		return nil
	}
	return m.DeleteFn(ctx, id)
}

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
