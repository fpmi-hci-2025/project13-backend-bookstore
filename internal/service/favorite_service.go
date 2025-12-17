package service

import (
	"context"
	"time"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/domain"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/repository"
	"github.com/google/uuid"
)

type FavoriteService struct {
	favoriteRepo repository.FavoriteRepository
}

func NewFavoriteService(favoriteRepo repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{favoriteRepo: favoriteRepo}
}

type AddFavoriteRequest struct {
	BookID uuid.UUID `json:"book_id" binding:"required"`
}

func (s *FavoriteService) Add(ctx context.Context, userID uuid.UUID, req *AddFavoriteRequest) (*domain.Favorite, error) {
	// Check if already in favorites
	exists, err := s.favoriteRepo.Exists(ctx, userID, req.BookID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrFavoriteAlreadyExists
	}

	favorite := &domain.Favorite{
		ID:        uuid.New(),
		UserID:    userID,
		BookID:    req.BookID,
		CreatedAt: time.Now(),
	}

	if err := s.favoriteRepo.Add(ctx, favorite); err != nil {
		return nil, err
	}

	return favorite, nil
}

func (s *FavoriteService) GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]domain.Favorite, error) {
	return s.favoriteRepo.GetByUserID(ctx, userID)
}

func (s *FavoriteService) Remove(ctx context.Context, userID uuid.UUID, bookID uuid.UUID) error {
	return s.favoriteRepo.RemoveByUserAndBook(ctx, userID, bookID)
}

func (s *FavoriteService) IsFavorite(ctx context.Context, userID uuid.UUID, bookID uuid.UUID) (bool, error) {
	return s.favoriteRepo.Exists(ctx, userID, bookID)
}

