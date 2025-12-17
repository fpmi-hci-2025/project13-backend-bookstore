package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(db *pgxpool.Pool) error {
	ctx := context.Background()

	migrations := []string{
		// Enable UUID extension
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,

		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// Authors table
		`CREATE TABLE IF NOT EXISTS authors (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			role VARCHAR(100) DEFAULT 'Writer',
			image_url TEXT,
			biography TEXT,
			rating DECIMAL(3,2) DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// Books table
		`CREATE TABLE IF NOT EXISTS books (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			title VARCHAR(255) NOT NULL,
			author_id UUID REFERENCES authors(id) ON DELETE SET NULL,
			genre VARCHAR(100),
			price DECIMAL(10,2) NOT NULL,
			stock INTEGER DEFAULT 0,
			image_url TEXT,
			description TEXT,
			rating DECIMAL(3,2) DEFAULT 0,
			review_count INTEGER DEFAULT 0,
			publisher VARCHAR(255),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// Orders table
		`CREATE TABLE IF NOT EXISTS orders (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			total_price DECIMAL(10,2) NOT NULL,
			delivery_address TEXT NOT NULL,
			status VARCHAR(50) DEFAULT 'pending',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,

		// Order items table
		`CREATE TABLE IF NOT EXISTS order_items (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
			book_id UUID REFERENCES books(id) ON DELETE SET NULL,
			quantity INTEGER NOT NULL,
			price DECIMAL(10,2) NOT NULL
		)`,

		// Reviews table
		`CREATE TABLE IF NOT EXISTS reviews (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			book_id UUID REFERENCES books(id) ON DELETE CASCADE,
			rating INTEGER CHECK (rating >= 1 AND rating <= 5),
			comment TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			UNIQUE(user_id, book_id)
		)`,

		// Cart items table
		`CREATE TABLE IF NOT EXISTS cart_items (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			book_id UUID REFERENCES books(id) ON DELETE CASCADE,
			quantity INTEGER DEFAULT 1,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			UNIQUE(user_id, book_id)
		)`,

		// Favorites table
		`CREATE TABLE IF NOT EXISTS favorites (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			book_id UUID REFERENCES books(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			UNIQUE(user_id, book_id)
		)`,

		// Indexes for better performance
		`CREATE INDEX IF NOT EXISTS idx_books_author_id ON books(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_books_genre ON books(genre)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_reviews_book_id ON reviews(book_id)`,
		`CREATE INDEX IF NOT EXISTS idx_cart_items_user_id ON cart_items(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(ctx, migration); err != nil {
			return err
		}
	}

	return nil
}

