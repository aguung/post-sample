package profile_test

import (
	"errors"
	"testing"

	"post/internal/entity"
	"post/internal/profile"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock of profile.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(p *entity.Profile) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockRepository) Update(p *entity.Profile) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockRepository) FindByUserID(userID uint) (*entity.Profile, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Profile), args.Error(1)
}

func TestCreateOrUpdate(t *testing.T) {
	t.Run("CreateNew", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := profile.NewService(mockRepo)
		userID := uint(1)
		input := profile.ProfileInput{Name: "New User", Bio: "Hello"}

		mockRepo.On("FindByUserID", userID).Return(nil, errors.New("not found"))
		mockRepo.On("Create", mock.MatchedBy(func(p *entity.Profile) bool {
			return p.UserID == userID && p.Name == input.Name && p.Bio == input.Bio
		})).Return(nil)

		result, err := service.CreateOrUpdate(userID, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, input.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateExisting", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := profile.NewService(mockRepo)
		userID := uint(1)
		input := profile.ProfileInput{Name: "Updated User", Bio: "Updated Bio"}
		existingProfile := &entity.Profile{UserID: userID, Name: "Old Name", Bio: "Old Bio"}

		mockRepo.On("FindByUserID", userID).Return(existingProfile, nil)
		mockRepo.On("Update", mock.MatchedBy(func(p *entity.Profile) bool {
			return p.UserID == userID && p.Name == input.Name && p.Bio == input.Bio
		})).Return(nil)

		result, err := service.CreateOrUpdate(userID, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, input.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByUserID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := profile.NewService(mockRepo)
		userID := uint(1)
		expectedProfile := &entity.Profile{UserID: userID, Name: "Test User"}

		mockRepo.On("FindByUserID", userID).Return(expectedProfile, nil)

		result, err := service.GetByUserID(userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedProfile, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := profile.NewService(mockRepo)
		userID := uint(2)
		mockRepo.On("FindByUserID", userID).Return(nil, errors.New("profile not found"))

		result, err := service.GetByUserID(userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
