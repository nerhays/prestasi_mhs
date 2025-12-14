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
func (s *AuthService) RefreshToken(oldToken string) (string, error) {
	claims, err := utils.ParseToken(oldToken)
	if err != nil {
		return "", errors.New("invalid or expired token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil || !user.IsActive {
		return "", errors.New("user not found or inactive")
	}

	perms, err := s.userRepo.GetPermissionsByUserID(user.ID)
	if err != nil {
		return "", err
	}

	return utils.GenerateToken(user, perms)
}
func (s *AuthService) GetProfile(userID string) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
