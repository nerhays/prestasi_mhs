package service

import (
	"context"
	"errors"

	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository"
)

type StudentService struct {
	studentRepo repository.StudentRepository
	lecturerRepo repository.LecturerRepository
}

func NewStudentService(
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
) *StudentService {
	return &StudentService{
		studentRepo:  studentRepo,
		lecturerRepo: lecturerRepo,
	}
}

func (s *StudentService) GetProfileByUserID(userID string) (*model.Student, error) {
	return s.studentRepo.FindByUserID(userID)
}
func (s *StudentService) SetAdvisor(
	ctx context.Context,
	studentID string,
	lecturerID string,
) error {

	// 1. cek student
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return ErrStudentProfileNotFound
	}

	// 2. cek lecturer
	lect, err := s.lecturerRepo.FindByID(lecturerID)
	if err != nil {
		return ErrLecturerNotFound
	}

	// 3. pastikan role = Dosen Wali
	if lect.User.Role.Name != "Dosen Wali" {
		return errors.New("selected user is not a dosen wali")
	}

	// 4. update advisor
	return s.studentRepo.UpdateAdvisor(student.ID, lect.ID)
}
func (s *StudentService) GetAllStudents() ([]model.Student, error) {
	return s.studentRepo.FindAll()
}

func (s *StudentService) GetStudentByID(id string) (*model.Student, error) {
	return s.studentRepo.FindByID(id)
}


