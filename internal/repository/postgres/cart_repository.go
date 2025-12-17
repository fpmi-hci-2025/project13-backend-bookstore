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

type cartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) repository.CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddItem(ctx context.Context, item *domain.CartItem) error {
	query := `
		INSERT INTO cart_items (id, user_id, book_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, book_id) DO UPDATE SET quantity = cart_items.quantity + $4, updated_at = $6
	`
	_, err := r.db.Exec(ctx, query,
		item.ID,
		item.UserID,
		item.BookID,
		item.Quantity,
		item.CreatedAt,
		item.UpdatedAt,
	)
	return err
}

func (r *cartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.CartItem, error) {
	query := `
		SELECT ci.id, ci.user_id, ci.book_id, ci.quantity, ci.created_at, ci.updated_at,
			   b.title, b.price, b.image_url, b.stock
		FROM cart_items ci
		JOIN books b ON ci.book_id = b.id
		WHERE ci.user_id = $1
		ORDER BY ci.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.CartItem
	for rows.Next() {
		var item domain.CartItem
		var book domain.Book
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.BookID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
			&book.Title,
			&book.Price,
			&book.ImageURL,
			&book.Stock,
		); err != nil {
			return nil, err
		}
		book.ID = item.BookID
		item.Book = &book
		items = append(items, item)
	}
	return items, nil
}

func (r *cartRepository) GetItem(ctx context.Context, userID, bookID uuid.UUID) (*domain.CartItem, error) {
	query := `
		SELECT id, user_id, book_id, quantity, created_at, updated_at
		FROM cart_items WHERE user_id = $1 AND book_id = $2
	`
	var item domain.CartItem
	err := r.db.QueryRow(ctx, query, userID, bookID).Scan(
		&item.ID,
		&item.UserID,
		&item.BookID,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCartItemNotFound
	}
	return &item, err
}

func (r *cartRepository) UpdateQuantity(ctx context.Context, id uuid.UUID, quantity int) error {
	query := `UPDATE cart_items SET quantity = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, quantity)
	return err
}

func (r *cartRepository) RemoveItem(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *cartRepository) ClearCart(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

