package mocks

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/stretchr/testify/mock"
)

type AchievementReferenceRepositoryMock struct {
	mock.Mock
}

func (m *AchievementReferenceRepositoryMock) CreateDraft(studentID, mongoID string) (*model.AchievementReference, error) {
	args := m.Called(studentID, mongoID)
	return args.Get(0).(*model.AchievementReference), args.Error(1)
}

func (m *AchievementReferenceRepositoryMock) GetByID(id string) (*model.AchievementReference, error) {
	args := m.Called(id)
	return args.Get(0).(*model.AchievementReference), args.Error(1)
}

func (m *AchievementReferenceRepositoryMock) Save(ref *model.AchievementReference) error {
	args := m.Called(ref)
	return args.Error(0)
}

func (m *AchievementReferenceRepositoryMock) FindByStudentID(studentID string) ([]model.AchievementReference, error) {
	args := m.Called(studentID)
	return args.Get(0).([]model.AchievementReference), args.Error(1)
}

func (m *AchievementReferenceRepositoryMock) FindByStudentIDs(ids []string, status *model.AchievementStatus, limit, offset int) ([]model.AchievementReference, error) {
	args := m.Called(ids, status, limit, offset)
	return args.Get(0).([]model.AchievementReference), args.Error(1)
}

func (m *AchievementReferenceRepositoryMock) FindAll(offset, limit int, status *string) ([]model.AchievementReference, int64, error) {
	args := m.Called(offset, limit, status)
	return args.Get(0).([]model.AchievementReference), args.Get(1).(int64), args.Error(2)
}

func (m *AchievementReferenceRepositoryMock) CountByStudentIDs(ids []string, status *model.AchievementStatus) (int64, error) {
	args := m.Called(ids, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *AchievementReferenceRepositoryMock) CountByStatus() (map[string]int64, error) {
	args := m.Called()
	return args.Get(0).(map[string]int64), args.Error(1)
}
