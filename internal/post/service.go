package post

import (
	"log"
	"post/internal/entity"
	"post/internal/pkg/cache"
)

type Service interface {
	Create(userID uint, input CreatePostInput) (*entity.Post, error)
	GetAll() ([]entity.Post, error)
	GetByID(id uint) (*entity.Post, error)
	GetByUserID(userID uint) ([]entity.Post, error)
	Delete(id uint) error
}

type service struct {
	repo  Repository
	cache cache.Cache
}

func NewService(repo Repository, cache cache.Cache) Service {
	return &service{repo, cache}
}

type CreatePostInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (s *service) Create(userID uint, input CreatePostInput) (*entity.Post, error) {
	post := &entity.Post{
		UserID:  userID,
		Title:   input.Title,
		Content: input.Content,
	}
	if err := s.repo.Create(post); err != nil {
		return nil, err
	}
	// Invalidate Cache
	s.cache.Delete("all_posts")
	return post, nil
}

func (s *service) GetAll() ([]entity.Post, error) {
	// Check Cache
	if val, ok := s.cache.Get("all_posts"); ok {
		log.Println("Hit Cache: all_posts")
		return val.([]entity.Post), nil
	}

	// Fetch DB
	posts, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Set Cache
	s.cache.Set("all_posts", posts)
	log.Println("Miss Cache: all_posts (Set)")
	return posts, nil
}

func (s *service) GetByID(id uint) (*entity.Post, error) {
	return s.repo.FindByID(id)
}

func (s *service) GetByUserID(userID uint) ([]entity.Post, error) {
	return s.repo.FindByUserID(userID)
}

func (s *service) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	// Invalidate Cache
	s.cache.Delete("all_posts")
	return nil
}
