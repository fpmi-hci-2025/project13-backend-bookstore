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

// GetBookReviews godoc
// @Summary Get reviews for a book
// @Description Get a paginated list of reviews for a specific book
// @Tags reviews
// @Produce json
// @Param id path string true "Book ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {array} domain.Review
// @Router /books/{id}/reviews [get]
func (h *Handler) GetBookReviews(c *gin.Context) {
	bookID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	reviews, err := h.reviewService.GetByBookID(c.Request.Context(), bookID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// CreateReview godoc
// @Summary Create a new review
// @Description Create a review for a book
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateReviewRequest true "Review details"
// @Success 201 {object} domain.Review
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /reviews [post]
func (h *Handler) CreateReview(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.reviewService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrReviewAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "You already reviewed this book"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// UpdateReview godoc
// @Summary Update a review
// @Description Update an existing review
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Review ID"
// @Param request body service.UpdateReviewRequest true "Review details"
// @Success 200 {object} domain.Review
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /reviews/{id} [put]
func (h *Handler) UpdateReview(c *gin.Context) {
	userID := middleware.GetUserID(c)

	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var req service.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.reviewService.Update(c.Request.Context(), reviewID, userID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrReviewNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own reviews"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	c.JSON(http.StatusOK, review)
}

// DeleteReview godoc
// @Summary Delete a review
// @Description Delete an existing review
// @Tags reviews
// @Security BearerAuth
// @Param id path string true "Review ID"
// @Success 204 "No Content"
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /reviews/{id} [delete]
func (h *Handler) DeleteReview(c *gin.Context) {
	userID := middleware.GetUserID(c)

	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := h.reviewService.Delete(c.Request.Context(), reviewID, userID); err != nil {
		if errors.Is(err, domain.ErrReviewNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own reviews"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	c.Status(http.StatusNoContent)
}

