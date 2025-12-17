package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/domain"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetAuthors godoc
// @Summary Get all authors
// @Description Get a paginated list of all authors
// @Tags authors
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {array} domain.Author
// @Router /authors [get]
func (h *Handler) GetAuthors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	authors, err := h.authorService.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get authors"})
		return
	}

	c.JSON(http.StatusOK, authors)
}

// GetAuthor godoc
// @Summary Get an author by ID
// @Description Get detailed information about an author including their books
// @Tags authors
// @Produce json
// @Param id path string true "Author ID"
// @Success 200 {object} domain.Author
// @Failure 404 {object} map[string]string
// @Router /authors/{id} [get]
func (h *Handler) GetAuthor(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	author, err := h.authorService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrAuthorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get author"})
		return
	}

	c.JSON(http.StatusOK, author)
}

// CreateAuthor godoc
// @Summary Create a new author
// @Description Create a new author (requires authentication)
// @Tags authors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateAuthorRequest true "Author details"
// @Success 201 {object} domain.Author
// @Failure 400 {object} map[string]string
// @Router /authors [post]
func (h *Handler) CreateAuthor(c *gin.Context) {
	var req service.CreateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	author, err := h.authorService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
		return
	}

	c.JSON(http.StatusCreated, author)
}

// UpdateAuthor godoc
// @Summary Update an author
// @Description Update an existing author (requires authentication)
// @Tags authors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Author ID"
// @Param request body service.UpdateAuthorRequest true "Author details"
// @Success 200 {object} domain.Author
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /authors/{id} [put]
func (h *Handler) UpdateAuthor(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	var req service.UpdateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	author, err := h.authorService.Update(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, domain.ErrAuthorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		return
	}

	c.JSON(http.StatusOK, author)
}

// DeleteAuthor godoc
// @Summary Delete an author
// @Description Delete an author (requires authentication)
// @Tags authors
// @Security BearerAuth
// @Param id path string true "Author ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /authors/{id} [delete]
func (h *Handler) DeleteAuthor(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	if err := h.authorService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete author"})
		return
	}

	c.Status(http.StatusNoContent)
}

