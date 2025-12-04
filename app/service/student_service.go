package service

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository"
)

type StudentService struct {
	studentRepo repository.StudentRepository
}

func NewStudentService(studentRepo repository.StudentRepository) *StudentService {
	return &StudentService{studentRepo: studentRepo}
}

func (s *StudentService) GetProfileByUserID(userID string) (*model.Student, error) {
	return s.studentRepo.FindByUserID(userID)
}
