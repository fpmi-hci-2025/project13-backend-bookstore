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

type favoriteRepository struct {
	db *pgxpool.Pool
}

func NewFavoriteRepository(db *pgxpool.Pool) repository.FavoriteRepository {
	return &favoriteRepository{db: db}
}

func (r *favoriteRepository) Add(ctx context.Context, favorite *domain.Favorite) error {
	query := `
		INSERT INTO favorites (id, user_id, book_id, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, query,
		favorite.ID,
		favorite.UserID,
		favorite.BookID,
		favorite.CreatedAt,
	)
	return err
}

func (r *favoriteRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Favorite, error) {
	query := `
		SELECT f.id, f.user_id, f.book_id, f.created_at,
			   b.title, b.price, b.image_url
		FROM favorites f
		JOIN books b ON f.book_id = b.id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []domain.Favorite
	for rows.Next() {
		var fav domain.Favorite
		var book domain.Book
		if err := rows.Scan(
			&fav.ID,
			&fav.UserID,
			&fav.BookID,
			&fav.CreatedAt,
			&book.Title,
			&book.Price,
			&book.ImageURL,
		); err != nil {
			return nil, err
		}
		book.ID = fav.BookID
		fav.Book = &book
		favorites = append(favorites, fav)
	}
	return favorites, nil
}

func (r *favoriteRepository) Exists(ctx context.Context, userID, bookID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND book_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, bookID).Scan(&exists)
	return exists, err
}

func (r *favoriteRepository) Remove(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM favorites WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrFavoriteNotFound
	}
	return nil
}

func (r *favoriteRepository) RemoveByUserAndBook(ctx context.Context, userID, bookID uuid.UUID) error {
	query := `DELETE FROM favorites WHERE user_id = $1 AND book_id = $2`
	result, err := r.db.Exec(ctx, query, userID, bookID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrFavoriteNotFound
	}
	return nil
}

// Ensure interface compliance
var _ repository.FavoriteRepository = (*favoriteRepository)(nil)

// Helper function
func isFavoriteNotFoundError(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

