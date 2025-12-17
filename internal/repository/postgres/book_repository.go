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

type bookRepository struct {
	db *pgxpool.Pool
}

func NewBookRepository(db *pgxpool.Pool) repository.BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(ctx context.Context, book *domain.Book) error {
	query := `
		INSERT INTO books (id, title, author_id, genre, price, stock, image_url, description, rating, review_count, publisher, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(ctx, query,
		book.ID,
		book.Title,
		book.AuthorID,
		book.Genre,
		book.Price,
		book.Stock,
		book.ImageURL,
		book.Description,
		book.Rating,
		book.ReviewCount,
		book.Publisher,
		book.CreatedAt,
		book.UpdatedAt,
	)
	return err
}

func (r *bookRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	query := `
		SELECT b.id, b.title, b.author_id, b.genre, b.price, b.stock, b.image_url, b.description, b.rating, b.review_count, b.publisher, b.created_at, b.updated_at,
			   a.id, a.name, a.role, a.image_url, a.biography, a.rating
		FROM books b
		LEFT JOIN authors a ON b.author_id = a.id
		WHERE b.id = $1
	`
	var book domain.Book
	var author domain.Author
	var authorID, authorName, authorRole, authorImageURL, authorBiography *string
	var authorRating *float64

	err := r.db.QueryRow(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.AuthorID,
		&book.Genre,
		&book.Price,
		&book.Stock,
		&book.ImageURL,
		&book.Description,
		&book.Rating,
		&book.ReviewCount,
		&book.Publisher,
		&book.CreatedAt,
		&book.UpdatedAt,
		&authorID,
		&authorName,
		&authorRole,
		&authorImageURL,
		&authorBiography,
		&authorRating,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBookNotFound
	}
	if err != nil {
		return nil, err
	}

	if authorID != nil {
		author.ID = book.AuthorID
		if authorName != nil {
			author.Name = *authorName
		}
		if authorRole != nil {
			author.Role = *authorRole
		}
		if authorImageURL != nil {
			author.ImageURL = *authorImageURL
		}
		if authorBiography != nil {
			author.Biography = *authorBiography
		}
		if authorRating != nil {
			author.Rating = *authorRating
		}
		book.Author = &author
	}

	return &book, nil
}

func (r *bookRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Book, error) {
	query := `
		SELECT id, title, author_id, genre, price, stock, image_url, description, rating, review_count, publisher, created_at, updated_at
		FROM books
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var book domain.Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.AuthorID,
			&book.Genre,
			&book.Price,
			&book.Stock,
			&book.ImageURL,
			&book.Description,
			&book.Rating,
			&book.ReviewCount,
			&book.Publisher,
			&book.CreatedAt,
			&book.UpdatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *bookRepository) GetByAuthorID(ctx context.Context, authorID uuid.UUID) ([]domain.Book, error) {
	query := `
		SELECT id, title, author_id, genre, price, stock, image_url, description, rating, review_count, publisher, created_at, updated_at
		FROM books
		WHERE author_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var book domain.Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.AuthorID,
			&book.Genre,
			&book.Price,
			&book.Stock,
			&book.ImageURL,
			&book.Description,
			&book.Rating,
			&book.ReviewCount,
			&book.Publisher,
			&book.CreatedAt,
			&book.UpdatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *bookRepository) Search(ctx context.Context, query string, limit, offset int) ([]domain.Book, error) {
	sqlQuery := `
		SELECT id, title, author_id, genre, price, stock, image_url, description, rating, review_count, publisher, created_at, updated_at
		FROM books
		WHERE title ILIKE $1 OR description ILIKE $1
		ORDER BY rating DESC
		LIMIT $2 OFFSET $3
	`
	searchTerm := "%" + query + "%"
	rows, err := r.db.Query(ctx, sqlQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var book domain.Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.AuthorID,
			&book.Genre,
			&book.Price,
			&book.Stock,
			&book.ImageURL,
			&book.Description,
			&book.Rating,
			&book.ReviewCount,
			&book.Publisher,
			&book.CreatedAt,
			&book.UpdatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *bookRepository) GetByGenre(ctx context.Context, genre string, limit, offset int) ([]domain.Book, error) {
	query := `
		SELECT id, title, author_id, genre, price, stock, image_url, description, rating, review_count, publisher, created_at, updated_at
		FROM books
		WHERE genre = $1
		ORDER BY rating DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, genre, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var book domain.Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.AuthorID,
			&book.Genre,
			&book.Price,
			&book.Stock,
			&book.ImageURL,
			&book.Description,
			&book.Rating,
			&book.ReviewCount,
			&book.Publisher,
			&book.CreatedAt,
			&book.UpdatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *bookRepository) Update(ctx context.Context, book *domain.Book) error {
	query := `
		UPDATE books SET title = $2, author_id = $3, genre = $4, price = $5, stock = $6, image_url = $7, description = $8, publisher = $9, updated_at = $10
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		book.ID,
		book.Title,
		book.AuthorID,
		book.Genre,
		book.Price,
		book.Stock,
		book.ImageURL,
		book.Description,
		book.Publisher,
		book.UpdatedAt,
	)
	return err
}

func (r *bookRepository) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	query := `UPDATE books SET stock = stock + $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, quantity)
	return err
}

func (r *bookRepository) UpdateRating(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE books 
		SET rating = (SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE book_id = $1),
			review_count = (SELECT COUNT(*) FROM reviews WHERE book_id = $1),
			updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *bookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *bookRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM books`).Scan(&count)
	return count, err
}

