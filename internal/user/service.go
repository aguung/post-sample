package user

import (
	"post/internal/entity"
)

type Service interface {
	GetByID(id uint) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetAll() ([]entity.User, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetByID(id uint) (*entity.User, error) {
	return s.repo.FindByID(id)
}

func (s *service) GetByEmail(email string) (*entity.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *service) GetAll() ([]entity.User, error) {
	return s.repo.FindAll()
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(id)
}
