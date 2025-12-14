package service

import (
	"errors"

	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository"
	"golang.org/x/crypto/bcrypt"
)
type UserService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
) *UserService {
	return &UserService{userRepo, roleRepo}
}
func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.userRepo.FindAll()
}
func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.userRepo.FindByID(id)
}
func (s *UserService) CreateUser(
	username, email, password, fullName, roleID string,
) error {

	// cek role valid
	if _, err := s.roleRepo.FindByID(roleID); err != nil {
		return errors.New("role not found")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
		RoleID:       roleID,
		IsActive:     true,
	}

	return s.userRepo.Create(user)
}
func (s *UserService) UpdateUser(
	id, username, email, fullName string,
) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	user.Username = username
	user.Email = email
	user.FullName = fullName

	return s.userRepo.Update(user)
}
func (s *UserService) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}
func (s *UserService) UpdateUserRole(userID, roleID string) error {
	if _, err := s.roleRepo.FindByID(roleID); err != nil {
		return errors.New("role not found")
	}

	return s.userRepo.UpdateRole(userID, roleID)
}
