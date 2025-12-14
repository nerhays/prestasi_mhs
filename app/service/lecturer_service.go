package service

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository"
)
type LecturerService struct {
	lecturerRepo repository.LecturerRepository
	studentRepo  repository.StudentRepository
}

func NewLecturerService(
	lecturerRepo repository.LecturerRepository,
	studentRepo repository.StudentRepository,
) *LecturerService {
	return &LecturerService{lecturerRepo, studentRepo}
}

func (s *LecturerService) GetAllLecturers() ([]model.Lecturer, error) {
	return s.lecturerRepo.FindAll()
}

func (s *LecturerService) GetAdvisees(lecturerID string) ([]model.Student, error) {
	return s.studentRepo.FindByAdvisorID(lecturerID)
}
