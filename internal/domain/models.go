package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Book represents a book in the store
type Book struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	AuthorID    uuid.UUID `json:"author_id"`
	Author      *Author   `json:"author,omitempty"`
	Genre       string    `json:"genre"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	Rating      float64   `json:"rating"`
	ReviewCount int       `json:"review_count"`
	Publisher   string    `json:"publisher"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Author represents a book author
type Author struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	ImageURL  string    `json:"image_url"`
	Biography string    `json:"biography"`
	Rating    float64   `json:"rating"`
	Books     []Book    `json:"books,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents a customer order
type Order struct {
	ID              uuid.UUID   `json:"id"`
	UserID          uuid.UUID   `json:"user_id"`
	User            *User       `json:"user,omitempty"`
	Items           []OrderItem `json:"items"`
	TotalPrice      float64     `json:"total_price"`
	DeliveryAddress string      `json:"delivery_address"`
	Status          OrderStatus `json:"status"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// OrderItem represents a single item in an order
type OrderItem struct {
	ID       uuid.UUID `json:"id"`
	OrderID  uuid.UUID `json:"order_id"`
	BookID   uuid.UUID `json:"book_id"`
	Book     *Book     `json:"book,omitempty"`
	Quantity int       `json:"quantity"`
	Price    float64   `json:"price"`
}

// Review represents a book review
type Review struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	User      *User     `json:"user,omitempty"`
	BookID    uuid.UUID `json:"book_id"`
	Book      *Book     `json:"book,omitempty"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	BookID    uuid.UUID `json:"book_id"`
	Book      *Book     `json:"book,omitempty"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Favorite represents a user's favorite book
type Favorite struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	BookID    uuid.UUID `json:"book_id"`
	Book      *Book     `json:"book,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

