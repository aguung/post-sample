package profile

import (
	"post/internal/entity"
)

type Service interface {
	CreateOrUpdate(userID uint, input ProfileInput) (*entity.Profile, error)
	GetByUserID(userID uint) (*entity.Profile, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

type ProfileInput struct {
	Name string `json:"name" binding:"required"`
	Bio  string `json:"bio"`
}

func (s *service) CreateOrUpdate(userID uint, input ProfileInput) (*entity.Profile, error) {
	profile, err := s.repo.FindByUserID(userID)
	if err != nil {
		// Create new if not exists (assuming error means not found, handled better with specific error check usually)
		// Simpler logic:
		newProfile := &entity.Profile{
			UserID: userID,
			Name:   input.Name,
			Bio:    input.Bio,
		}
		if err := s.repo.Create(newProfile); err != nil {
			return nil, err
		}
		return newProfile, nil
	}

	// Update existing
	profile.Name = input.Name
	profile.Bio = input.Bio
	if err := s.repo.Update(profile); err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *service) GetByUserID(userID uint) (*entity.Profile, error) {
	return s.repo.FindByUserID(userID)
}
