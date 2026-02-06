package post

import (
	"post/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	Create(post *entity.Post) error
	FindAll() ([]entity.Post, error)
	FindByID(id uint) (*entity.Post, error)
	FindByUserID(userID uint) ([]entity.Post, error)
	Update(post *entity.Post) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(post *entity.Post) error {
	return r.db.Create(post).Error
}

func (r *repository) FindAll() ([]entity.Post, error) {
	var posts []entity.Post
	err := r.db.Preload("User").Joins("JOIN users ON posts.user_id = users.id").Where("users.deleted_at IS NULL").Find(&posts).Error
	return posts, err
}

func (r *repository) FindByID(id uint) (*entity.Post, error) {
	var post entity.Post
	err := r.db.First(&post, id).Error
	return &post, err
}

func (r *repository) FindByUserID(userID uint) ([]entity.Post, error) {
	var posts []entity.Post
	err := r.db.Where("user_id = ?", userID).Find(&posts).Error
	return posts, err
}

func (r *repository) Update(post *entity.Post) error {
	return r.db.Save(post).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entity.Post{}, id).Error
}
