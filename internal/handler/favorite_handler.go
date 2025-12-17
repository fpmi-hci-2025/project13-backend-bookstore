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

// GetFavorites godoc
// @Summary Get user's favorites
// @Description Get all books in user's favorites
// @Tags favorites
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Favorite
// @Router /favorites [get]
func (h *Handler) GetFavorites(c *gin.Context) {
	userID := middleware.GetUserID(c)

	favorites, err := h.favoriteService.GetUserFavorites(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorites"})
		return
	}

	c.JSON(http.StatusOK, favorites)
}

// AddFavorite godoc
// @Summary Add book to favorites
// @Description Add a book to user's favorites
// @Tags favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.AddFavoriteRequest true "Book to add"
// @Success 201 {object} domain.Favorite
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /favorites [post]
func (h *Handler) AddFavorite(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	favorite, err := h.favoriteService.Add(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrFavoriteAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Book already in favorites"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to favorites"})
		return
	}

	c.JSON(http.StatusCreated, favorite)
}

// RemoveFavorite godoc
// @Summary Remove book from favorites
// @Description Remove a book from user's favorites
// @Tags favorites
// @Security BearerAuth
// @Param book_id path string true "Book ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /favorites/{book_id} [delete]
func (h *Handler) RemoveFavorite(c *gin.Context) {
	userID := middleware.GetUserID(c)

	bookID, err := uuid.Parse(c.Param("book_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := h.favoriteService.Remove(c.Request.Context(), userID, bookID); err != nil {
		if errors.Is(err, domain.ErrFavoriteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not in favorites"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from favorites"})
		return
	}

	c.Status(http.StatusNoContent)
}

