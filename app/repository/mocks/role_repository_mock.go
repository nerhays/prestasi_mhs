package mocks

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/stretchr/testify/mock"
)

type RoleRepositoryMock struct {
	mock.Mock
}

func (m *RoleRepositoryMock) FindAll() ([]model.Role, error) {
	args := m.Called()
	return args.Get(0).([]model.Role), args.Error(1)
}

func (m *RoleRepositoryMock) FindByID(id string) (*model.Role, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Role), args.Error(1)
}
