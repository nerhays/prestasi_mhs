package service

import (
	"testing"

	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllUsers_Success(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	roleRepo := new(mocks.RoleRepositoryMock)

	svc := NewUserService(userRepo, roleRepo)

	expectedUsers := []model.User{
		{ID: "u1", Username: "user1"},
		{ID: "u2", Username: "user2"},
	}

	userRepo.On("FindAll").Return(expectedUsers, nil)

	users, err := svc.GetAllUsers()

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "u1", users[0].ID)
	assert.Equal(t, "u2", users[1].ID)

	userRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	roleRepo := new(mocks.RoleRepositoryMock)

	svc := NewUserService(userRepo, roleRepo)

	roleRepo.On("FindByID", "role-1").Return(&model.Role{ID: "role-1"}, nil)
	userRepo.On("Create", mock.Anything).Return(nil)

	err := svc.CreateUser(
		"john",
		"john@mail.com",
		"password",
		"John Doe",
		"role-1",
	)

	assert.NoError(t, err)
}

func TestUpdateUserRole_Success(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	roleRepo := new(mocks.RoleRepositoryMock)

	svc := NewUserService(userRepo, roleRepo)

	roleRepo.On("FindByID", "role-2").Return(&model.Role{}, nil)
	userRepo.On("UpdateRole", "user-1", "role-2").Return(nil)

	err := svc.UpdateUserRole("user-1", "role-2")

	assert.NoError(t, err)
}
