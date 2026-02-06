package user_test

import (
	"errors"
	"testing"

	"post/internal/entity"
	"post/internal/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock of user.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(u *entity.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockRepository) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockRepository) FindByID(id uint) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockRepository) FindAll() ([]entity.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetByID(t *testing.T) {
	mockRepo := new(MockRepository)
	service := user.NewService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedUser := &entity.User{ID: 1, Email: "test@example.com"}
		mockRepo.On("FindByID", uint(1)).Return(expectedUser, nil)

		result, err := service.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.On("FindByID", uint(2)).Return(nil, errors.New("user not found"))

		result, err := service.GetByID(2)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByEmail(t *testing.T) {
	mockRepo := new(MockRepository)
	service := user.NewService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedUser := &entity.User{ID: 1, Email: "test@example.com"}
		mockRepo.On("FindByEmail", "test@example.com").Return(expectedUser, nil)

		result, err := service.GetByEmail("test@example.com")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})
}
