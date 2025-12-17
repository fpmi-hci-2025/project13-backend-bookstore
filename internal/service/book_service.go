package service

import (
	"context"
	"time"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/domain"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/repository"
	"github.com/google/uuid"
)

type BookService struct {
	bookRepo repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}
}

type CreateBookRequest struct {
	Title       string    `json:"title" binding:"required"`
	AuthorID    uuid.UUID `json:"author_id" binding:"required"`
	Genre       string    `json:"genre"`
	Price       float64   `json:"price" binding:"required,gt=0"`
	Stock       int       `json:"stock" binding:"min=0"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	Publisher   string    `json:"publisher"`
}

type UpdateBookRequest struct {
	Title       string    `json:"title"`
	AuthorID    uuid.UUID `json:"author_id"`
	Genre       string    `json:"genre"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	Publisher   string    `json:"publisher"`
}

type BookListResponse struct {
	Books      []domain.Book `json:"books"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

func (s *BookService) Create(ctx context.Context, req *CreateBookRequest) (*domain.Book, error) {
	book := &domain.Book{
		ID:          uuid.New(),
		Title:       req.Title,
		AuthorID:    req.AuthorID,
		Genre:       req.Genre,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		Description: req.Description,
		Publisher:   req.Publisher,
		Rating:      0,
		ReviewCount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.bookRepo.Create(ctx, book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	return s.bookRepo.GetByID(ctx, id)
}

func (s *BookService) GetAll(ctx context.Context, page, pageSize int) (*BookListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	books, err := s.bookRepo.GetAll(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.bookRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &BookListResponse{
		Books:      books,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *BookService) Search(ctx context.Context, query string, page, pageSize int) (*BookListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	books, err := s.bookRepo.Search(ctx, query, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &BookListResponse{
		Books:    books,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *BookService) GetByGenre(ctx context.Context, genre string, page, pageSize int) (*BookListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	books, err := s.bookRepo.GetByGenre(ctx, genre, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &BookListResponse{
		Books:    books,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *BookService) Update(ctx context.Context, id uuid.UUID, req *UpdateBookRequest) (*domain.Book, error) {
	book, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		book.Title = req.Title
	}
	if req.AuthorID != uuid.Nil {
		book.AuthorID = req.AuthorID
	}
	if req.Genre != "" {
		book.Genre = req.Genre
	}
	if req.Price > 0 {
		book.Price = req.Price
	}
	if req.Stock >= 0 {
		book.Stock = req.Stock
	}
	if req.ImageURL != "" {
		book.ImageURL = req.ImageURL
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if req.Publisher != "" {
		book.Publisher = req.Publisher
	}
	book.UpdatedAt = time.Now()

	if err := s.bookRepo.Update(ctx, book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.bookRepo.Delete(ctx, id)
}

