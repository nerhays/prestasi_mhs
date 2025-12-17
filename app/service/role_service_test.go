package service

import (
	"testing"

	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetAllRoles_Success(t *testing.T) {
	roleRepo := new(mocks.RoleRepositoryMock)

	svc := NewRoleService(roleRepo)

	roles := []model.Role{
		{ID: "r1", Name: "Admin"},
		{ID: "r2", Name: "Mahasiswa"},
	}

	roleRepo.On("FindAll").Return(roles, nil)

	res, err := svc.GetAllRoles()

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "Admin", res[0].Name)
}
