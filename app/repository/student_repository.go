package repository

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"gorm.io/gorm"
)

type StudentRepository interface {
	FindByUserID(userID string) (*model.Student, error)
	FindByID(id string) (*model.Student, error)
	FindByAdvisorLecturerID(lecturerID string) ([]model.Student, error)
	UpdateAdvisor(studentID, lecturerID string) error
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
func (r *studentRepository) FindByID(id string) (*model.Student, error) {
	var student model.Student
	if err := r.db.
		Preload("User.Role").
		Where("id = ?", id).
		First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}
func (r *studentRepository) FindByAdvisorLecturerID(lecturerID string) ([]model.Student, error) {
    var students []model.Student
    if err := r.db.Where("advisor_id = ?", lecturerID).Find(&students).Error; err != nil {
        return nil, err
    }
    return students, nil
}
func (r *studentRepository) UpdateAdvisor(studentID, lecturerID string) error {
	return r.db.
		Model(&model.Student{}).
		Where("id = ?", studentID).
		Update("advisor_id", lecturerID).
		Error
}

