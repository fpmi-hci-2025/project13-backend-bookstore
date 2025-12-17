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

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) repository.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert order
	orderQuery := `
		INSERT INTO orders (id, user_id, total_price, delivery_address, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, orderQuery,
		order.ID,
		order.UserID,
		order.TotalPrice,
		order.DeliveryAddress,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Insert order items
	itemQuery := `
		INSERT INTO order_items (id, order_id, book_id, quantity, price)
		VALUES ($1, $2, $3, $4, $5)
	`
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, itemQuery,
			item.ID,
			order.ID,
			item.BookID,
			item.Quantity,
			item.Price,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	orderQuery := `
		SELECT id, user_id, total_price, delivery_address, status, created_at, updated_at
		FROM orders WHERE id = $1
	`
	var order domain.Order
	err := r.db.QueryRow(ctx, orderQuery, id).Scan(
		&order.ID,
		&order.UserID,
		&order.TotalPrice,
		&order.DeliveryAddress,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}

	// Get order items
	itemsQuery := `
		SELECT oi.id, oi.order_id, oi.book_id, oi.quantity, oi.price,
			   b.title, b.image_url
		FROM order_items oi
		LEFT JOIN books b ON oi.book_id = b.id
		WHERE oi.order_id = $1
	`
	rows, err := r.db.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.OrderItem
		var bookTitle, bookImageURL *string
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.BookID,
			&item.Quantity,
			&item.Price,
			&bookTitle,
			&bookImageURL,
		); err != nil {
			return nil, err
		}
		if bookTitle != nil {
			item.Book = &domain.Book{
				ID:       item.BookID,
				Title:    *bookTitle,
				ImageURL: *bookImageURL,
			}
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Order, error) {
	query := `
		SELECT id, user_id, total_price, delivery_address, status, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalPrice,
			&order.DeliveryAddress,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	// Load items for each order
	for i := range orders {
		itemsQuery := `
			SELECT oi.id, oi.order_id, oi.book_id, oi.quantity, oi.price,
				   b.title, b.image_url
			FROM order_items oi
			LEFT JOIN books b ON oi.book_id = b.id
			WHERE oi.order_id = $1
		`
		itemRows, err := r.db.Query(ctx, itemsQuery, orders[i].ID)
		if err != nil {
			return nil, err
		}

		for itemRows.Next() {
			var item domain.OrderItem
			var bookTitle, bookImageURL *string
			if err := itemRows.Scan(
				&item.ID,
				&item.OrderID,
				&item.BookID,
				&item.Quantity,
				&item.Price,
				&bookTitle,
				&bookImageURL,
			); err != nil {
				itemRows.Close()
				return nil, err
			}
			if bookTitle != nil {
				item.Book = &domain.Book{
					ID:       item.BookID,
					Title:    *bookTitle,
					ImageURL: *bookImageURL,
				}
			}
			orders[i].Items = append(orders[i].Items, item)
		}
		itemRows.Close()
	}

	return orders, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	query := `UPDATE orders SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, status)
	return err
}

func (r *orderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

