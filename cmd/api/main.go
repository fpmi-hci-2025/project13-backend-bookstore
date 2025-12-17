package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/config"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/handler"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/repository/postgres"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title BookStore API
// @version 1.0
// @description API для книжного онлайн-магазина
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@bookstore.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	db, err := postgres.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := postgres.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	bookRepo := postgres.NewBookRepository(db)
	authorRepo := postgres.NewAuthorRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	reviewRepo := postgres.NewReviewRepository(db)
	cartRepo := postgres.NewCartRepository(db)
	favoriteRepo := postgres.NewFavoriteRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	bookService := service.NewBookService(bookRepo)
	authorService := service.NewAuthorService(authorRepo, bookRepo)
	orderService := service.NewOrderService(orderRepo, cartRepo, bookRepo)
	reviewService := service.NewReviewService(reviewRepo, bookRepo)
	cartService := service.NewCartService(cartRepo, bookRepo)
	favoriteService := service.NewFavoriteService(favoriteRepo)

	// Initialize handlers
	h := handler.NewHandler(
		authService,
		bookService,
		authorService,
		orderService,
		reviewService,
		cartService,
		favoriteService,
		cfg.JWTSecret,
	)

	// Setup router
	router := h.SetupRoutes()

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

