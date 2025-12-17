package mocks

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/stretchr/testify/mock"
)

type LecturerRepositoryMock struct {
	mock.Mock
}

func (m *LecturerRepositoryMock) FindByID(id string) (*model.Lecturer, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Lecturer), args.Error(1)
}

func (m *LecturerRepositoryMock) FindByUserID(userID string) (*model.Lecturer, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.Lecturer), args.Error(1)
}

func (m *LecturerRepositoryMock) FindAll() ([]model.Lecturer, error) {
	args := m.Called()
	return args.Get(0).([]model.Lecturer), args.Error(1)
}
