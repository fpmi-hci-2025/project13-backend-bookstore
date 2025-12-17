package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/domain"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/handler/middleware"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetOrders godoc
// @Summary Get user's orders
// @Description Get a paginated list of user's orders
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {array} domain.Order
// @Router /orders [get]
func (h *Handler) GetOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, err := h.orderService.GetUserOrders(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder godoc
// @Summary Get order by ID
// @Description Get detailed information about a specific order
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} domain.Order
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [get]
func (h *Handler) GetOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetByID(c.Request.Context(), orderID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order"})
		return
	}

	// Check if user owns the order
	if order.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create an order from cart items
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateOrderRequest true "Order details"
// @Success 201 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Router /orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyCart) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
			return
		}
		if errors.Is(err, domain.ErrInsufficientStock) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for some items"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel a pending order
// @Tags orders
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{id}/cancel [post]
func (h *Handler) CancelOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if err := h.orderService.Cancel(c.Request.Context(), orderID, userID); err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot cancel this order"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled"})
}

