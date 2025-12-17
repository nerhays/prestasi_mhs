package service_test

import (
	"errors"
	"testing"

	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository/mocks"
	"github.com/nerhays/prestasi_uas/app/service"
	"github.com/nerhays/prestasi_uas/utils"
	"github.com/stretchr/testify/assert"
)

func TestLoginSuccess(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)

	hashed, _ := utils.HashPassword("password123")

	user := &model.User{
		ID:           "user-1",
		Username:     "admin",
		PasswordHash: hashed,
		IsActive:     true,
		Role: model.Role{
			Name: "Admin",
		},
	}

	perms := []model.Permission{
		{Name: "read"},
	}

	userRepo.On("FindByUsernameOrEmail", "admin").
		Return(user, nil)
	userRepo.On("GetPermissionsByUserID", "user-1").
		Return(perms, nil)

	authSvc := service.NewAuthService(userRepo)

	res, err := authSvc.Login(service.LoginInput{
		Username: "admin",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, "admin", res.User.Username)
	userRepo.AssertExpectations(t)
}

func TestLoginWrongPassword(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)

	hashed, _ := utils.HashPassword("password123")

	userRepo.On("FindByUsernameOrEmail", "admin").
		Return(&model.User{
			ID:           "user-1",
			PasswordHash: hashed,
			IsActive:     true,
		}, nil)

	authSvc := service.NewAuthService(userRepo)

	res, err := authSvc.Login(service.LoginInput{
		Username: "admin",
		Password: "wrong",
	})

	assert.Nil(t, res)
	assert.EqualError(t, err, "invalid_credentials")
}

func TestLoginUserInactive(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)

	userRepo.On("FindByUsernameOrEmail", "admin").
		Return(&model.User{
			ID:       "user-1",
			IsActive: false,
		}, nil)

	authSvc := service.NewAuthService(userRepo)

	res, err := authSvc.Login(service.LoginInput{
		Username: "admin",
		Password: "password",
	})

	assert.Nil(t, res)
	assert.EqualError(t, err, "user_inactive")
}

func TestRefreshTokenSuccess(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)

	hashed, _ := utils.HashPassword("password")

	user := &model.User{
		ID:           "user-1",
		PasswordHash: hashed,
		IsActive:     true,
		Role: model.Role{Name: "Admin"},
	}

	perms := []model.Permission{{Name: "read"}}
	token, _ := utils.GenerateToken(user, perms)

	userRepo.On("FindByID", "user-1").Return(user, nil)
	userRepo.On("GetPermissionsByUserID", "user-1").Return(perms, nil)

	authSvc := service.NewAuthService(userRepo)

	newToken, err := authSvc.RefreshToken(token)

	assert.NoError(t, err)
	assert.NotEmpty(t, newToken)
}

func TestGetProfileSuccess(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)

	user := &model.User{
		ID:       "user-1",
		Username: "admin",
	}

	userRepo.On("FindByID", "user-1").Return(user, nil)

	authSvc := service.NewAuthService(userRepo)

	res, err := authSvc.GetProfile("user-1")

	assert.NoError(t, err)
	assert.Equal(t, "admin", res.Username)
}

func TestGetProfileNotFound(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)

	userRepo.On("FindByID", "x").
		Return(nil, errors.New("not found"))

	authSvc := service.NewAuthService(userRepo)

	res, err := authSvc.GetProfile("x")

	assert.Nil(t, res)
	assert.EqualError(t, err, "user not found")
}
