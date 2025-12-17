package service

import (
	"testing"

	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetAllLecturers_Success(t *testing.T) {
	lectRepo := new(mocks.LecturerRepositoryMock)
	studentRepo := new(mocks.StudentRepositoryMock)

	svc := NewLecturerService(lectRepo, studentRepo)

	expected := []model.Lecturer{
		{ID: "lect-1"},
		{ID: "lect-2"},
	}

	lectRepo.On("FindAll").Return(expected, nil)

	res, err := svc.GetAllLecturers()

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "lect-1", res[0].ID)
}

func TestGetAdvisees_Success(t *testing.T) {
	lectRepo := new(mocks.LecturerRepositoryMock)
	studentRepo := new(mocks.StudentRepositoryMock)

	svc := NewLecturerService(lectRepo, studentRepo)

	lecturerID := "lect-1"
	students := []model.Student{
		{ID: "student-1"},
		{ID: "student-2"},
	}

	studentRepo.On("FindByAdvisorID", lecturerID).Return(students, nil)

	res, err := svc.GetAdvisees(lecturerID)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
}
