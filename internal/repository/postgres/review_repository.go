package postgres

import (
	"context"
	"errors"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/domain"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type reviewRepository struct {
	db *pgxpool.Pool
}

func NewReviewRepository(db *pgxpool.Pool) repository.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *domain.Review) error {
	query := `
		INSERT INTO reviews (id, user_id, book_id, rating, comment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		review.ID,
		review.UserID,
		review.BookID,
		review.Rating,
		review.Comment,
		review.CreatedAt,
		review.UpdatedAt,
	)
	return err
}

func (r *reviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error) {
	query := `
		SELECT r.id, r.user_id, r.book_id, r.rating, r.comment, r.created_at, r.updated_at,
			   u.username
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.id = $1
	`
	var review domain.Review
	var username *string
	err := r.db.QueryRow(ctx, query, id).Scan(
		&review.ID,
		&review.UserID,
		&review.BookID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
		&review.UpdatedAt,
		&username,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrReviewNotFound
	}
	if err != nil {
		return nil, err
	}
	if username != nil {
		review.User = &domain.User{ID: review.UserID, Username: *username}
	}
	return &review, nil
}

func (r *reviewRepository) GetByBookID(ctx context.Context, bookID uuid.UUID, limit, offset int) ([]domain.Review, error) {
	query := `
		SELECT r.id, r.user_id, r.book_id, r.rating, r.comment, r.created_at, r.updated_at,
			   u.username
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.book_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, bookID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []domain.Review
	for rows.Next() {
		var review domain.Review
		var username *string
		if err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.BookID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
			&review.UpdatedAt,
			&username,
		); err != nil {
			return nil, err
		}
		if username != nil {
			review.User = &domain.User{ID: review.UserID, Username: *username}
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (r *reviewRepository) GetByUserAndBook(ctx context.Context, userID, bookID uuid.UUID) (*domain.Review, error) {
	query := `
		SELECT id, user_id, book_id, rating, comment, created_at, updated_at
		FROM reviews WHERE user_id = $1 AND book_id = $2
	`
	var review domain.Review
	err := r.db.QueryRow(ctx, query, userID, bookID).Scan(
		&review.ID,
		&review.UserID,
		&review.BookID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrReviewNotFound
	}
	return &review, err
}

func (r *reviewRepository) Update(ctx context.Context, review *domain.Review) error {
	query := `
		UPDATE reviews SET rating = $2, comment = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		review.ID,
		review.Rating,
		review.Comment,
		review.UpdatedAt,
	)
	return err
}

func (r *reviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reviews WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

