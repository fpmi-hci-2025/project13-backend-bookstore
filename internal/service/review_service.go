package service

import (
	"context"
	"time"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/domain"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/repository"
	"github.com/google/uuid"
)

type ReviewService struct {
	reviewRepo repository.ReviewRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo}
}

type CreateReviewRequest struct {
	BookID  uuid.UUID `json:"book_id" binding:"required"`
	Rating  int       `json:"rating" binding:"required,min=1,max=5"`
	Comment string    `json:"comment"`
}

type UpdateReviewRequest struct {
	Rating  int    `json:"rating" binding:"min=1,max=5"`
	Comment string `json:"comment"`
}

func (s *ReviewService) Create(ctx context.Context, userID uuid.UUID, req *CreateReviewRequest) (*domain.Review, error) {
	// Check if user already reviewed this book
	if _, err := s.reviewRepo.GetByUserAndBook(ctx, userID, req.BookID); err == nil {
		return nil, domain.ErrReviewAlreadyExists
	}

	review := &domain.Review{
		ID:        uuid.New(),
		UserID:    userID,
		BookID:    req.BookID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) GetByBookID(ctx context.Context, bookID uuid.UUID, page, pageSize int) ([]domain.Review, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.reviewRepo.GetByBookID(ctx, bookID, pageSize, offset)
}

func (s *ReviewService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *UpdateReviewRequest) (*domain.Review, error) {
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if user owns the review
	if review.UserID != userID {
		return nil, domain.ErrForbidden
	}

	if req.Rating > 0 {
		review.Rating = req.Rating
	}
	if req.Comment != "" {
		review.Comment = req.Comment
	}
	review.UpdatedAt = time.Now()

	if err := s.reviewRepo.Update(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if user owns the review
	if review.UserID != userID {
		return domain.ErrForbidden
	}

	return s.reviewRepo.Delete(ctx, id)
}

