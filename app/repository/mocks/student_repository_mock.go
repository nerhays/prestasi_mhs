package mocks

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/stretchr/testify/mock"
)

type StudentRepositoryMock struct {
	mock.Mock
}

func (m *StudentRepositoryMock) FindByUserID(userID string) (*model.Student, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *StudentRepositoryMock) FindByID(id string) (*model.Student, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *StudentRepositoryMock) FindAll() ([]model.Student, error) {
	args := m.Called()
	return args.Get(0).([]model.Student), args.Error(1)
}

func (m *StudentRepositoryMock) FindByAdvisorID(advisorID string) ([]model.Student, error) {
	args := m.Called(advisorID)
	return args.Get(0).([]model.Student), args.Error(1)
}

func (m *StudentRepositoryMock) FindByAdvisorLecturerID(lecturerID string) ([]model.Student, error) {
	args := m.Called(lecturerID)
	return args.Get(0).([]model.Student), args.Error(1)
}

func (m *StudentRepositoryMock) SetAdvisor(studentID, advisorID string) error {
	args := m.Called(studentID, advisorID)
	return args.Error(0)
}
func (m *StudentRepositoryMock) UpdateAdvisor(studentID, advisorID string) error {
	args := m.Called(studentID, advisorID)
	return args.Error(0)
}
