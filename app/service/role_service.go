package service

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository"
)

type RoleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) GetAllRoles() ([]model.Role, error) {
	return s.roleRepo.FindAll()
}
