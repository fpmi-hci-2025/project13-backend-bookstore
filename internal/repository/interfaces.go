package repository

import (
	"context"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/domain"
	"github.com/google/uuid"
)

// UserRepository defines methods for user data access
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// BookRepository defines methods for book data access
type BookRepository interface {
	Create(ctx context.Context, book *domain.Book) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Book, error)
	GetByAuthorID(ctx context.Context, authorID uuid.UUID) ([]domain.Book, error)
	Search(ctx context.Context, query string, limit, offset int) ([]domain.Book, error)
	GetByGenre(ctx context.Context, genre string, limit, offset int) ([]domain.Book, error)
	Update(ctx context.Context, book *domain.Book) error
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
	UpdateRating(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int, error)
}

// AuthorRepository defines methods for author data access
type AuthorRepository interface {
	Create(ctx context.Context, author *domain.Author) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Author, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Author, error)
	Update(ctx context.Context, author *domain.Author) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// OrderRepository defines methods for order data access
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ReviewRepository defines methods for review data access
type ReviewRepository interface {
	Create(ctx context.Context, review *domain.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error)
	GetByBookID(ctx context.Context, bookID uuid.UUID, limit, offset int) ([]domain.Review, error)
	GetByUserAndBook(ctx context.Context, userID, bookID uuid.UUID) (*domain.Review, error)
	Update(ctx context.Context, review *domain.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// CartRepository defines methods for cart data access
type CartRepository interface {
	AddItem(ctx context.Context, item *domain.CartItem) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.CartItem, error)
	GetItem(ctx context.Context, userID, bookID uuid.UUID) (*domain.CartItem, error)
	UpdateQuantity(ctx context.Context, id uuid.UUID, quantity int) error
	RemoveItem(ctx context.Context, id uuid.UUID) error
	ClearCart(ctx context.Context, userID uuid.UUID) error
}

// FavoriteRepository defines methods for favorites data access
type FavoriteRepository interface {
	Add(ctx context.Context, favorite *domain.Favorite) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Favorite, error)
	Exists(ctx context.Context, userID, bookID uuid.UUID) (bool, error)
	Remove(ctx context.Context, id uuid.UUID) error
	RemoveByUserAndBook(ctx context.Context, userID, bookID uuid.UUID) error
}

