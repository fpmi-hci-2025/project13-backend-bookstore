package service

import (
	"context"
	"testing"

	"github.com/fpmi-hci-2025/project13-backend-bookstore/internal/domain"
	"github.com/google/uuid"
)

// MockUserRepository implements repository.UserRepository for testing
type MockUserRepository struct {
	users map[string]*domain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return nil
}

func TestAuthService_Register(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewAuthService(repo, "test-secret")

	ctx := context.Background()

	// Test successful registration
	req := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := service.Register(ctx, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Token == "" {
		t.Error("Expected token to be generated")
	}

	if resp.User.Username != req.Username {
		t.Errorf("Expected username %s, got %s", req.Username, resp.User.Username)
	}

	if resp.User.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, resp.User.Email)
	}

	// Test duplicate registration
	_, err = service.Register(ctx, req)
	if err != domain.ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestAuthService_Login(t *testing.T) {
	repo := NewMockUserRepository()
	service := NewAuthService(repo, "test-secret")

	ctx := context.Background()

	// First register a user
	registerReq := &RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	_, err := service.Register(ctx, registerReq)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Test successful login
	loginReq := &LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := service.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Token == "" {
		t.Error("Expected token to be generated")
	}

	// Test login with wrong password
	loginReq.Password = "wrongpassword"
	_, err = service.Login(ctx, loginReq)
	if err != domain.ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}

	// Test login with non-existent email
	loginReq.Email = "nonexistent@example.com"
	_, err = service.Login(ctx, loginReq)
	if err != domain.ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

