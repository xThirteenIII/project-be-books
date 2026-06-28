package repo

import (
	"books/internal/model"
	"context"
	"database/sql"
	"errors"
	"time"
)

// Concrete implementation of ReviewRepo
type MariaDBReviewRepo struct {
	db *sql.DB
}

// Verify implementation at compile time.
// var _ ReviewRepo = &MariaDBReviewRepo{}

var ErrNotFound = errors.New("review not found")

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

func (r *MariaDBReviewRepo) GetByID(ctx context.Context, id string) (model.Review, error) {
	query := `SELECT id, book_id, score, review_text, status, book_title, authors, cover_url, created_at, updated_at
			  FROM reviews
			  WHERE id = ?`

	var review model.Review

	// Use of NullString resolves the problematic scan of NULL fields into strings.
	// These fields can be NULL while review status is pending.
	var bookTitle sql.NullString
	var authors sql.NullString
	var coverURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&review.ID,
		&review.BookID,
		&review.Score,
		&review.ReviewText,
		&review.Status,
		&bookTitle,
		&authors,
		&coverURL,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	// If there's no reviews corresponding to given ID, return ErrNotFound.
	if errors.Is(err, sql.ErrNoRows) {
		return model.Review{}, ErrNotFound
	}
	if err != nil {
		return model.Review{}, err
	}

	// NULL fields become empty strings.
	review.BookTitle = bookTitle.String
	review.Authors = authors.String
	review.CoverURL = coverURL.String

	return review, nil
}

func (r *MariaDBReviewRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM reviews WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were updated, review id does not exist.
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *MariaDBReviewRepo) Update(ctx context.Context, review model.Review) error {
	query := `UPDATE reviews
			  SET book_title = ?, authors = ?, cover_url = ?, status = ?, updated_at = ?
			  WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query,
		review.BookTitle,
		review.Authors,
		review.CoverURL,
		review.Status,
		review.UpdatedAt,
		review.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateReviewContent updates score and review_text only.
// We avoid using Update() method for both PUT and POST endpoints.
func (r *MariaDBReviewRepo) UpdateReviewContent(ctx context.Context, id string, reviewText string, score int) error {
	query := `UPDATE reviews
			  SET review_text = ?, score = ?, updated_at = ?
			  WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query, reviewText, score, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
