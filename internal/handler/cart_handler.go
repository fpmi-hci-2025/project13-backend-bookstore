package handler

import (
	"errors"
	"net/http"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/domain"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/handler/middleware"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetCart godoc
// @Summary Get user's cart
// @Description Get all items in the user's shopping cart
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} service.CartResponse
// @Router /cart [get]
func (h *Handler) GetCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	cart, err := h.cartService.GetCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddToCart godoc
// @Summary Add item to cart
// @Description Add a book to the shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.AddToCartRequest true "Cart item details"
// @Success 201 {object} domain.CartItem
// @Failure 400 {object} map[string]string
// @Router /cart [post]
func (h *Handler) AddToCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.cartService.AddItem(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		if errors.Is(err, domain.ErrInsufficientStock) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// UpdateCartItem godoc
// @Summary Update cart item quantity
// @Description Update the quantity of an item in the cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Cart item ID"
// @Param request body service.UpdateCartItemRequest true "Updated quantity"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /cart/{id} [put]
func (h *Handler) UpdateCartItem(c *gin.Context) {
	userID := middleware.GetUserID(c)

	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req service.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.UpdateQuantity(c.Request.Context(), userID, itemID, &req); err != nil {
		if errors.Is(err, domain.ErrCartItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
			return
		}
		if errors.Is(err, domain.ErrInsufficientStock) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item updated"})
}

// RemoveFromCart godoc
// @Summary Remove item from cart
// @Description Remove a specific item from the cart
// @Tags cart
// @Security BearerAuth
// @Param id path string true "Cart item ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /cart/{id} [delete]
func (h *Handler) RemoveFromCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	if err := h.cartService.RemoveItem(c.Request.Context(), userID, itemID); err != nil {
		if errors.Is(err, domain.ErrCartItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove cart item"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ClearCart godoc
// @Summary Clear the cart
// @Description Remove all items from the cart
// @Tags cart
// @Security BearerAuth
// @Success 204 "No Content"
// @Router /cart [delete]
func (h *Handler) ClearCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if err := h.cartService.ClearCart(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.Status(http.StatusNoContent)
}

