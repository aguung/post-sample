package auth_test

import (
	"errors"
	"testing"

	"post/internal/auth"
	"post/internal/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock of user.Repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(u *entity.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]entity.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockJWTService is a mock of auth.JWTService
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(user *entity.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(user *entity.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func TestSignup(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockJWT := new(MockJWTService)
	service := auth.NewService(mockRepo, mockJWT)

	t.Run("Success", func(t *testing.T) {
		input := auth.SignupInput{
			Email:    "new@example.com",
			Password: "password",
		}

		mockRepo.On("FindByEmail", input.Email).Return(nil, errors.New("not found"))
		mockRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(nil)

		user, err := service.Signup(input)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, input.Email, user.Email)
		assert.NotEqual(t, input.Password, user.Password) // Password should be hashed
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmailAlreadyExists", func(t *testing.T) {
		input := auth.SignupInput{
			Email:    "existing@example.com",
			Password: "password",
		}
		existingUser := &entity.User{ID: 1, Email: input.Email}

		mockRepo.On("FindByEmail", input.Email).Return(existingUser, nil)

		user, err := service.Signup(input)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "email already registered", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestSignin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockJWT := new(MockJWTService)
	service := auth.NewService(mockRepo, mockJWT)

	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &entity.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     entity.RoleUser,
	}

	t.Run("Success", func(t *testing.T) {
		input := auth.SigninInput{
			Email:    "test@example.com",
			Password: "password",
		}

		mockRepo.On("FindByEmail", input.Email).Return(user, nil)
		mockJWT.On("GenerateToken", user).Return("mock_token", nil)
		mockJWT.On("GenerateRefreshToken", user).Return("mock_refresh_token", nil)

		token, refreshToken, err := service.Signin(input)

		assert.NoError(t, err)
		assert.Equal(t, "mock_token", token)
		assert.Equal(t, "mock_refresh_token", refreshToken)
		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		input := auth.SigninInput{
			Email:    "test@example.com",
			Password: "wrong_password",
		}

		mockRepo.On("FindByEmail", input.Email).Return(user, nil)

		token, refreshToken, err := service.Signin(input)

		assert.Error(t, err)
		assert.Equal(t, "", token)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, "invalid email or password", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		input := auth.SigninInput{
			Email:    "unknown@example.com",
			Password: "password",
		}

		mockRepo.On("FindByEmail", input.Email).Return(nil, errors.New("user not found"))

		token, refreshToken, err := service.Signin(input)

		assert.Error(t, err)
		assert.Equal(t, "", token)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, "invalid email or password", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
