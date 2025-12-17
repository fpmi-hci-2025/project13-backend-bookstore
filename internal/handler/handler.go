package handler

import (
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/handler/middleware"
	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/fpmi-hci-2025/project13-backend-bookstore/docs"
)

type Handler struct {
	authService     *service.AuthService
	bookService     *service.BookService
	authorService   *service.AuthorService
	orderService    *service.OrderService
	reviewService   *service.ReviewService
	cartService     *service.CartService
	favoriteService *service.FavoriteService
	jwtSecret       string
}

func NewHandler(
	authService *service.AuthService,
	bookService *service.BookService,
	authorService *service.AuthorService,
	orderService *service.OrderService,
	reviewService *service.ReviewService,
	cartService *service.CartService,
	favoriteService *service.FavoriteService,
	jwtSecret string,
) *Handler {
	return &Handler{
		authService:     authService,
		bookService:     bookService,
		authorService:   authorService,
		orderService:    orderService,
		reviewService:   reviewService,
		cartService:     cartService,
		favoriteService: favoriteService,
		jwtSecret:       jwtSecret,
	}
}

func (h *Handler) SetupRoutes() *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", h.healthCheck)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
		}

		// Books routes (public read, protected write)
		books := v1.Group("/books")
		{
			books.GET("", h.GetBooks)
			books.GET("/:id", h.GetBook)
			books.GET("/search", h.SearchBooks)
			books.GET("/:id/reviews", h.GetBookReviews)

			// Protected
			booksProtected := books.Group("")
			booksProtected.Use(middleware.Auth(h.jwtSecret))
			{
				booksProtected.POST("", h.CreateBook)
				booksProtected.PUT("/:id", h.UpdateBook)
				booksProtected.DELETE("/:id", h.DeleteBook)
			}
		}

		// Authors routes (public read, protected write)
		authors := v1.Group("/authors")
		{
			authors.GET("", h.GetAuthors)
			authors.GET("/:id", h.GetAuthor)

			// Protected
			authorsProtected := authors.Group("")
			authorsProtected.Use(middleware.Auth(h.jwtSecret))
			{
				authorsProtected.POST("", h.CreateAuthor)
				authorsProtected.PUT("/:id", h.UpdateAuthor)
				authorsProtected.DELETE("/:id", h.DeleteAuthor)
			}
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.Auth(h.jwtSecret))
		{
			// User profile
			protected.GET("/me", h.GetProfile)
			protected.PUT("/me", h.UpdateProfile)

			// Cart
			cart := protected.Group("/cart")
			{
				cart.GET("", h.GetCart)
				cart.POST("", h.AddToCart)
				cart.PUT("/:id", h.UpdateCartItem)
				cart.DELETE("/:id", h.RemoveFromCart)
				cart.DELETE("", h.ClearCart)
			}

			// Orders
			orders := protected.Group("/orders")
			{
				orders.GET("", h.GetOrders)
				orders.POST("", h.CreateOrder)
				orders.GET("/:id", h.GetOrder)
				orders.POST("/:id/cancel", h.CancelOrder)
			}

			// Reviews
			reviews := protected.Group("/reviews")
			{
				reviews.POST("", h.CreateReview)
				reviews.PUT("/:id", h.UpdateReview)
				reviews.DELETE("/:id", h.DeleteReview)
			}

			// Favorites
			favorites := protected.Group("/favorites")
			{
				favorites.GET("", h.GetFavorites)
				favorites.POST("", h.AddFavorite)
				favorites.DELETE("/:book_id", h.RemoveFavorite)
			}
		}
	}

	return router
}

func (h *Handler) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "bookstore-api",
	})
}

