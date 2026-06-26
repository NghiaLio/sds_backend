package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mobile-crud-backend/model"
	"mobile-crud-backend/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService manages user authentication.
type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) error
	Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
}

type authServiceImpl struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authServiceImpl{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *authServiceImpl) Register(ctx context.Context, req model.RegisterRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// Check user uniqueness
	existingID, _, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if existingID > 0 {
		return errors.New("username is already taken")
	}

	// Hash password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	return s.userRepo.Create(ctx, req.Username, string(hashedPassword))
}

func (s *authServiceImpl) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	id, passwordHash, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, errors.New("invalid username or password")
	}

	// Compare plaintext password with hash
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT AccessToken (expires in 10 minutes)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  id,
		"username": req.Username,
		"exp":      time.Now().Add(10 * time.Minute).Unix(),
	})

	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	return &model.LoginResponse{
		AccessToken: tokenStr,
		Username:    req.Username,
	}, nil
}
