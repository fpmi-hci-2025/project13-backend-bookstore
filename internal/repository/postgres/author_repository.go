package postgres

import (
	"context"
	"errors"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/domain"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authorRepository struct {
	db *pgxpool.Pool
}

func NewAuthorRepository(db *pgxpool.Pool) repository.AuthorRepository {
	return &authorRepository{db: db}
}

func (r *authorRepository) Create(ctx context.Context, author *domain.Author) error {
	query := `
		INSERT INTO authors (id, name, role, image_url, biography, rating, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		author.ID,
		author.Name,
		author.Role,
		author.ImageURL,
		author.Biography,
		author.Rating,
		author.CreatedAt,
		author.UpdatedAt,
	)
	return err
}

func (r *authorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Author, error) {
	query := `
		SELECT id, name, role, image_url, biography, rating, created_at, updated_at
		FROM authors WHERE id = $1
	`
	var author domain.Author
	err := r.db.QueryRow(ctx, query, id).Scan(
		&author.ID,
		&author.Name,
		&author.Role,
		&author.ImageURL,
		&author.Biography,
		&author.Rating,
		&author.CreatedAt,
		&author.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAuthorNotFound
	}
	return &author, err
}

func (r *authorRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Author, error) {
	query := `
		SELECT id, name, role, image_url, biography, rating, created_at, updated_at
		FROM authors
		ORDER BY name
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []domain.Author
	for rows.Next() {
		var author domain.Author
		if err := rows.Scan(
			&author.ID,
			&author.Name,
			&author.Role,
			&author.ImageURL,
			&author.Biography,
			&author.Rating,
			&author.CreatedAt,
			&author.UpdatedAt,
		); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}

func (r *authorRepository) Update(ctx context.Context, author *domain.Author) error {
	query := `
		UPDATE authors SET name = $2, role = $3, image_url = $4, biography = $5, rating = $6, updated_at = $7
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		author.ID,
		author.Name,
		author.Role,
		author.ImageURL,
		author.Biography,
		author.Rating,
		author.UpdatedAt,
	)
	return err
}

func (r *authorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM authors WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

