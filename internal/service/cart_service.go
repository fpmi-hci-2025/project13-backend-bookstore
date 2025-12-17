package service

import (
	"context"
	"time"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/domain"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/repository"
	"github.com/google/uuid"
)

type CartService struct {
	cartRepo repository.CartRepository
	bookRepo repository.BookRepository
}

func NewCartService(cartRepo repository.CartRepository, bookRepo repository.BookRepository) *CartService {
	return &CartService{
		cartRepo: cartRepo,
		bookRepo: bookRepo,
	}
}

type AddToCartRequest struct {
	BookID   uuid.UUID `json:"book_id" binding:"required"`
	Quantity int       `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type CartResponse struct {
	Items      []domain.CartItem `json:"items"`
	TotalPrice float64           `json:"total_price"`
	TotalItems int               `json:"total_items"`
}

func (s *CartService) AddItem(ctx context.Context, userID uuid.UUID, req *AddToCartRequest) (*domain.CartItem, error) {
	// Check if book exists and has enough stock
	book, err := s.bookRepo.GetByID(ctx, req.BookID)
	if err != nil {
		return nil, err
	}

	if book.Stock < req.Quantity {
		return nil, domain.ErrInsufficientStock
	}

	item := &domain.CartItem{
		ID:        uuid.New(),
		UserID:    userID,
		BookID:    req.BookID,
		Quantity:  req.Quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.cartRepo.AddItem(ctx, item); err != nil {
		return nil, err
	}

	item.Book = book
	return item, nil
}

func (s *CartService) GetCart(ctx context.Context, userID uuid.UUID) (*CartResponse, error) {
	items, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var totalPrice float64
	var totalItems int

	for _, item := range items {
		if item.Book != nil {
			totalPrice += item.Book.Price * float64(item.Quantity)
		}
		totalItems += item.Quantity
	}

	return &CartResponse{
		Items:      items,
		TotalPrice: totalPrice,
		TotalItems: totalItems,
	}, nil
}

func (s *CartService) UpdateQuantity(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req *UpdateCartItemRequest) error {
	// Get cart item to verify ownership
	items, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	var found bool
	var bookID uuid.UUID
	for _, item := range items {
		if item.ID == itemID {
			found = true
			bookID = item.BookID
			break
		}
	}

	if !found {
		return domain.ErrCartItemNotFound
	}

	// Check stock
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return err
	}

	if book.Stock < req.Quantity {
		return domain.ErrInsufficientStock
	}

	return s.cartRepo.UpdateQuantity(ctx, itemID, req.Quantity)
}

func (s *CartService) RemoveItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	// Verify ownership
	items, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	var found bool
	for _, item := range items {
		if item.ID == itemID {
			found = true
			break
		}
	}

	if !found {
		return domain.ErrCartItemNotFound
	}

	return s.cartRepo.RemoveItem(ctx, itemID)
}

func (s *CartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	return s.cartRepo.ClearCart(ctx, userID)
}

