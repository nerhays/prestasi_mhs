package repository

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"gorm.io/gorm"
)

type StudentRepository interface {
	FindByUserID(userID string) (*model.Student, error)
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) FindByUserID(userID string) (*model.Student, error) {
	var student model.Student
	if err := r.db.
		Preload("User.Role").
		Where("user_id = ?", userID).
		First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}
