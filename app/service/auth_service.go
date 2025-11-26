package service

import (
	"errors"

	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository"
	"github.com/nerhays/prestasi_uas/utils"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

type LoginInput struct {
	Username string
	Password string
}

type LoginOutput struct {
	Token       string
	User        *model.User
	Permissions []model.Permission
}

func (s *AuthService) Login(input LoginInput) (*LoginOutput, error) {
	user, err := s.userRepo.FindByUsernameOrEmail(input.Username)
	if err != nil {
		return nil, errors.New("invalid_credentials")
	}

	if !user.IsActive {
		return nil, errors.New("user_inactive")
	}

	if !utils.CheckPassword(user.PasswordHash, input.Password) {
		return nil, errors.New("invalid_credentials")
	}

	perms, err := s.userRepo.GetPermissionsByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user, perms)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token:       token,
		User:        user,
		Permissions: perms,
	}, nil
}
