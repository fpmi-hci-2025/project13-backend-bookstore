package service

import (
	"context"
	"time"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/domain"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/repository"
	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo repository.OrderRepository
	cartRepo  repository.CartRepository
	bookRepo  repository.BookRepository
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, bookRepo repository.BookRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
		bookRepo:  bookRepo,
	}
}

type CreateOrderRequest struct {
	DeliveryAddress string `json:"delivery_address" binding:"required"`
}

func (s *OrderService) Create(ctx context.Context, userID uuid.UUID, req *CreateOrderRequest) (*domain.Order, error) {
	// Get cart items
	cartItems, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, domain.ErrEmptyCart
	}

	// Calculate total price and create order items
	var totalPrice float64
	var orderItems []domain.OrderItem

	for _, item := range cartItems {
		// Check stock
		book, err := s.bookRepo.GetByID(ctx, item.BookID)
		if err != nil {
			return nil, err
		}
		if book.Stock < item.Quantity {
			return nil, domain.ErrInsufficientStock
		}

		itemPrice := book.Price * float64(item.Quantity)
		totalPrice += itemPrice

		orderItems = append(orderItems, domain.OrderItem{
			ID:       uuid.New(),
			BookID:   item.BookID,
			Quantity: item.Quantity,
			Price:    itemPrice,
		})
	}

	// Create order
	order := &domain.Order{
		ID:              uuid.New(),
		UserID:          userID,
		Items:           orderItems,
		TotalPrice:      totalPrice,
		DeliveryAddress: req.DeliveryAddress,
		Status:          domain.OrderStatusPending,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Update stock for each book
	for _, item := range cartItems {
		if err := s.bookRepo.UpdateStock(ctx, item.BookID, -item.Quantity); err != nil {
			return nil, err
		}
	}

	// Clear cart
	if err := s.cartRepo.ClearCart(ctx, userID); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]domain.Order, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.orderRepo.GetByUserID(ctx, userID, pageSize, offset)
}

func (s *OrderService) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	return s.orderRepo.UpdateStatus(ctx, id, status)
}

func (s *OrderService) Cancel(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if user owns the order
	if order.UserID != userID {
		return domain.ErrForbidden
	}

	// Only pending orders can be cancelled
	if order.Status != domain.OrderStatusPending {
		return domain.ErrForbidden
	}

	// Restore stock
	for _, item := range order.Items {
		if err := s.bookRepo.UpdateStock(ctx, item.BookID, item.Quantity); err != nil {
			return err
		}
	}

	return s.orderRepo.UpdateStatus(ctx, id, domain.OrderStatusCancelled)
}

