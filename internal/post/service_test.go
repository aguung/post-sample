package post_test

import (
	"errors"
	"testing"
	"time"

	"post/internal/entity"
	"post/internal/post"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock of post.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(p *entity.Post) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockRepository) FindAll() ([]entity.Post, error) {
	args := m.Called()
	return args.Get(0).([]entity.Post), args.Error(1)
}

func (m *MockRepository) FindByID(id uint) (*entity.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockRepository) FindByUserID(userID uint) ([]entity.Post, error) {
	args := m.Called(userID)
	return args.Get(0).([]entity.Post), args.Error(1)
}

func (m *MockRepository) Update(p *entity.Post) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockCache is a mock of cache.Cache
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string) (any, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *MockCache) Set(key string, value any) {
	m.Called(key, value)
}

func (m *MockCache) Delete(key string) {
	m.Called(key)
}

func (m *MockCache) Purge() {
	m.Called()
}

func TestCreate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCache := new(MockCache)
		service := post.NewService(mockRepo, mockCache)
		userID := uint(1)
		input := post.CreatePostInput{Title: "Test Post", Content: "Content"}

		mockRepo.On("Create", mock.MatchedBy(func(p *entity.Post) bool {
			return p.UserID == userID && p.Title == input.Title && p.Content == input.Content
		})).Return(nil)

		// Expect cache invalidation
		mockCache.On("Delete", "all_posts").Return()

		result, err := service.Create(userID, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, input.Title, result.Title)
		mockRepo.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCache := new(MockCache)
		service := post.NewService(mockRepo, mockCache)
		userID := uint(1)
		input := post.CreatePostInput{Title: "Test Post", Content: "Content"}

		mockRepo.On("Create", mock.AnythingOfType("*entity.Post")).Return(errors.New("db error"))

		result, err := service.Create(userID, input)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
		// Cache should NOT be invalidated on failure
		mockCache.AssertNotCalled(t, "Delete", "all_posts")
	})
}

func TestGetAll(t *testing.T) {
	t.Run("Cache Hit", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCache := new(MockCache)
		service := post.NewService(mockRepo, mockCache)
		expectedPosts := []entity.Post{
			{ID: 1, Title: "Post 1", UserID: 1},
		}

		// Mock Cache Hit
		mockCache.On("Get", "all_posts").Return(expectedPosts, true)

		result, err := service.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, len(expectedPosts), len(result))
		mockCache.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "FindAll")
	})

	t.Run("Cache Miss (DB Success)", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCache := new(MockCache)
		service := post.NewService(mockRepo, mockCache)
		expectedPosts := []entity.Post{
			{ID: 1, Title: "Post 1", UserID: 1},
			{ID: 2, Title: "Post 2", UserID: 1},
		}

		// Mock Cache Miss
		mockCache.On("Get", "all_posts").Return(nil, false)

		// Mock DB Call
		mockRepo.On("FindAll").Return(expectedPosts, nil)

		// Mock Cache Set
		mockCache.On("Set", "all_posts", expectedPosts).Return()

		result, err := service.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, len(expectedPosts), len(result))
		mockRepo.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCache := new(MockCache)
		service := post.NewService(mockRepo, mockCache)
		postID := uint(1)
		expectedPost := &entity.Post{ID: postID, Title: "Test Post", CreatedAt: time.Now()}

		mockRepo.On("FindByID", postID).Return(expectedPost, nil)

		result, err := service.GetByID(postID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPost, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCache := new(MockCache)
		service := post.NewService(mockRepo, mockCache)
		postID := uint(2)

		mockRepo.On("FindByID", postID).Return(nil, errors.New("post not found"))

		result, err := service.GetByID(postID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
