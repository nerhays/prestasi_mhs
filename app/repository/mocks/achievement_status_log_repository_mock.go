package mocks

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/stretchr/testify/mock"
)

type AchievementStatusLogRepositoryMock struct {
	mock.Mock
}

func (m *AchievementStatusLogRepositoryMock) Create(log *model.AchievementStatusLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *AchievementStatusLogRepositoryMock) FindByReferenceID(refID string) ([]model.AchievementStatusLog, error) {
	args := m.Called(refID)
	return args.Get(0).([]model.AchievementStatusLog), args.Error(1)
}
