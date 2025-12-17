package service

import (
	"context"
	"time"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/domain"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/repository"
	"github.com/google/uuid"
)

type AuthorService struct {
	authorRepo repository.AuthorRepository
	bookRepo   repository.BookRepository
}

func NewAuthorService(authorRepo repository.AuthorRepository, bookRepo repository.BookRepository) *AuthorService {
	return &AuthorService{
		authorRepo: authorRepo,
		bookRepo:   bookRepo,
	}
}

type CreateAuthorRequest struct {
	Name      string  `json:"name" binding:"required"`
	Role      string  `json:"role"`
	ImageURL  string  `json:"image_url"`
	Biography string  `json:"biography"`
	Rating    float64 `json:"rating"`
}

type UpdateAuthorRequest struct {
	Name      string  `json:"name"`
	Role      string  `json:"role"`
	ImageURL  string  `json:"image_url"`
	Biography string  `json:"biography"`
	Rating    float64 `json:"rating"`
}

func (s *AuthorService) Create(ctx context.Context, req *CreateAuthorRequest) (*domain.Author, error) {
	author := &domain.Author{
		ID:        uuid.New(),
		Name:      req.Name,
		Role:      req.Role,
		ImageURL:  req.ImageURL,
		Biography: req.Biography,
		Rating:    req.Rating,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if author.Role == "" {
		author.Role = "Writer"
	}

	if err := s.authorRepo.Create(ctx, author); err != nil {
		return nil, err
	}

	return author, nil
}

func (s *AuthorService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Author, error) {
	author, err := s.authorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get author's books
	books, err := s.bookRepo.GetByAuthorID(ctx, id)
	if err != nil {
		return nil, err
	}
	author.Books = books

	return author, nil
}

func (s *AuthorService) GetAll(ctx context.Context, page, pageSize int) ([]domain.Author, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.authorRepo.GetAll(ctx, pageSize, offset)
}

func (s *AuthorService) Update(ctx context.Context, id uuid.UUID, req *UpdateAuthorRequest) (*domain.Author, error) {
	author, err := s.authorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		author.Name = req.Name
	}
	if req.Role != "" {
		author.Role = req.Role
	}
	if req.ImageURL != "" {
		author.ImageURL = req.ImageURL
	}
	if req.Biography != "" {
		author.Biography = req.Biography
	}
	if req.Rating > 0 {
		author.Rating = req.Rating
	}
	author.UpdatedAt = time.Now()

	if err := s.authorRepo.Update(ctx, author); err != nil {
		return nil, err
	}

	return author, nil
}

func (s *AuthorService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.authorRepo.Delete(ctx, id)
}

