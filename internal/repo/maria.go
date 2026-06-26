package repo

import (
	"books/internal/model"
	"context"
	"database/sql"
)

// Concrete implementation of ReviewRepo
type MariaDBReviewRepo struct {
	db *sql.DB
}

// Verify implementation at compile time.
// var _ ReviewRepo = &MariaDBReviewRepo{}

func NewMariaDBReviewRepository(db *sql.DB) *MariaDBReviewRepo {
	return &MariaDBReviewRepo{db: db}
}

func (r *MariaDBReviewRepo) Insert(ctx context.Context, review model.Review) error {
	query := `INSERT INTO reviews (id, book_id, score, review_text, status, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		review.ID,
		review.BookID,
		review.Score,
		review.ReviewText,
		review.Status,
		review.CreatedAt,
		review.UpdatedAt,
	)
	return err
}

func (r *MariaDBReviewRepo) GetByID(ctx context.Context, id string) error {
	query := `SELECT FROM reviews id
              VALUES ?`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *MariaDBReviewRepo) Delete(ctx context.Context, id string) error {
	query := `SELECT FROM reviews id
              VALUES ?`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *MariaDBReviewRepo) Update(ctx context.Context, review model.Review) error {
	query := `UPDATE reviews
			  SET book_title = ?, authors = ?, cover_url = ?, status = ?, updated_at = ?
			  WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		review.BookTitle,
		review.Authors,
		review.CoverURL,
		review.Status,
		review.UpdatedAt,
		review.ID,
	)
	return err
}
