package domain

import "errors"

var (
	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Book errors
	ErrBookNotFound     = errors.New("book not found")
	ErrInsufficientStock = errors.New("insufficient stock")

	// Author errors
	ErrAuthorNotFound = errors.New("author not found")

	// Order errors
	ErrOrderNotFound = errors.New("order not found")
	ErrEmptyCart     = errors.New("cart is empty")

	// Review errors
	ErrReviewNotFound      = errors.New("review not found")
	ErrReviewAlreadyExists = errors.New("review already exists for this book")

	// Cart errors
	ErrCartItemNotFound = errors.New("cart item not found")

	// Favorite errors
	ErrFavoriteNotFound      = errors.New("favorite not found")
	ErrFavoriteAlreadyExists = errors.New("book already in favorites")

	// General errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

