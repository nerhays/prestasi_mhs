package mocks

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) FindByUsernameOrEmail(usernameOrEmail string) (*model.User, error) {
	args := m.Called(usernameOrEmail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepositoryMock) GetPermissionsByUserID(userID string) ([]model.Permission, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Permission), args.Error(1)
}

func (m *UserRepositoryMock) FindByID(id string) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepositoryMock) FindAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *UserRepositoryMock) Create(user *model.User) error {
	return nil
}
func (m *UserRepositoryMock) Update(user *model.User) error {
	return nil
}
func (m *UserRepositoryMock) Delete(id string) error {
	return nil
}
func (m *UserRepositoryMock) UpdateRole(userID, roleID string) error {
	return nil
}
