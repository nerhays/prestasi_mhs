package service

import (
	"context"
	"testing"

	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetProfileByUserID_Success(t *testing.T) {
	studentRepo := new(mocks.StudentRepositoryMock)
	lectRepo := new(mocks.LecturerRepositoryMock)

	svc := NewStudentService(studentRepo, lectRepo)

	userID := "user-1"
	student := &model.Student{ID: "student-1"}

	studentRepo.On("FindByUserID", userID).Return(student, nil)

	res, err := svc.GetProfileByUserID(userID)

	assert.NoError(t, err)
	assert.Equal(t, "student-1", res.ID)
}

func TestSetAdvisor_Success(t *testing.T) {
	studentRepo := new(mocks.StudentRepositoryMock)
	lectRepo := new(mocks.LecturerRepositoryMock)

	svc := NewStudentService(studentRepo, lectRepo)

	student := &model.Student{ID: "student-1"}
	lecturer := &model.Lecturer{
		ID: "lect-1",
		User: model.User{
			Role: model.Role{Name: "Dosen Wali"},
		},
	}

	studentRepo.On("FindByID", "student-1").Return(student, nil)
	lectRepo.On("FindByID", "lect-1").Return(lecturer, nil)
	studentRepo.On("UpdateAdvisor", "student-1", "lect-1").Return(nil)

	err := svc.SetAdvisor(context.Background(), "student-1", "lect-1")

	assert.NoError(t, err)
}

func TestSetAdvisor_NotDosenWali(t *testing.T) {
	studentRepo := new(mocks.StudentRepositoryMock)
	lectRepo := new(mocks.LecturerRepositoryMock)

	svc := NewStudentService(studentRepo, lectRepo)

	student := &model.Student{ID: "student-1"}
	lecturer := &model.Lecturer{
		ID: "lect-1",
		User: model.User{
			Role: model.Role{Name: "Mahasiswa"},
		},
	}

	studentRepo.On("FindByID", "student-1").Return(student, nil)
	lectRepo.On("FindByID", "lect-1").Return(lecturer, nil)

	err := svc.SetAdvisor(context.Background(), "student-1", "lect-1")

	assert.Error(t, err)
}
