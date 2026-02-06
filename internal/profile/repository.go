package profile

import (
	"post/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	Create(profile *entity.Profile) error
	Update(profile *entity.Profile) error
	FindByUserID(userID uint) (*entity.Profile, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(profile *entity.Profile) error {
	return r.db.Create(profile).Error
}

func (r *repository) Update(profile *entity.Profile) error {
	return r.db.Save(profile).Error
}

func (r *repository) FindByUserID(userID uint) (*entity.Profile, error) {
	var profile entity.Profile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}
