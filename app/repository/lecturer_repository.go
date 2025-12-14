package repository

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"gorm.io/gorm"
)

type LecturerRepository interface {
	FindByID(id string) (*model.Lecturer, error)
	FindByUserID(userID string) (*model.Lecturer, error)
}

type lecturerRepository struct {
	db *gorm.DB
}

func NewLecturerRepository(db *gorm.DB) LecturerRepository {
	return &lecturerRepository{db: db}
}

func (r *lecturerRepository) FindByID(id string) (*model.Lecturer, error) {
	var lect model.Lecturer
	if err := r.db.
		Preload("User").
		Preload("User.Role").
		Where("id = ?", id).
		First(&lect).Error; err != nil {
		return nil, err
	}
	return &lect, nil
}
func (r *lecturerRepository) FindByUserID(userID string) (*model.Lecturer, error) {
	var lect model.Lecturer
	if err := r.db.
		Preload("User").
		Preload("User.Role").
		Where("user_id = ?", userID).
		First(&lect).Error; err != nil {
		return nil, err
	}
	return &lect, nil
}
