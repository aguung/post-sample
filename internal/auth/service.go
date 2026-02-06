package auth

import (
	"errors"

	"post/internal/entity"
	pkgdb "post/internal/pkg/database"
	"post/internal/user"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Signup(input SignupInput) (*entity.User, error)
	Signin(input SigninInput) (string, string, error)
}

type service struct {
	userRepo   user.Repository
	jwtService JWTService
}

func NewService(userRepo user.Repository, jwtService JWTService) Service {
	return &service{userRepo, jwtService}
}

type SignupInput struct {
	Email    string      `json:"email" binding:"required,email"`
	Password string      `json:"password" binding:"required,min=6"`
	Role     entity.Role `json:"role"`
}

type SigninInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *service) Signup(input SignupInput) (*entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	// Default to user role if not provided or restricted
	if user.Role == 0 {
		user.Role = entity.RoleUser
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, pkgdb.ParseError(err)
	}

	return user, nil
}

func (s *service) Signin(input SigninInput) (string, string, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", "", errors.New("invalid email or password")
	}

	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}
