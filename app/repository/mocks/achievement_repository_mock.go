package mocks

import (
	"context"

	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/stretchr/testify/mock"
)

type AchievementRepositoryMock struct {
	mock.Mock
}

func (m *AchievementRepositoryMock) Create(ctx context.Context, ac *model.Achievement) (*model.Achievement, error) {
	args := m.Called(ctx, ac)
	return args.Get(0).(*model.Achievement), args.Error(1)
}

func (m *AchievementRepositoryMock) FindByStudentID(ctx context.Context, studentID string) ([]model.Achievement, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).([]model.Achievement), args.Error(1)
}

func (m *AchievementRepositoryMock) SoftDelete(ctx context.Context, mongoID string) error {
	args := m.Called(ctx, mongoID)
	return args.Error(0)
}

func (m *AchievementRepositoryMock) FindDeletedByStudentID(ctx context.Context, studentID string) ([]model.Achievement, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).([]model.Achievement), args.Error(1)
}

func (m *AchievementRepositoryMock) FindByIDs(ctx context.Context, ids []string) ([]model.Achievement, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]model.Achievement), args.Error(1)
}

func (m *AchievementRepositoryMock) AddAttachment(ctx context.Context, mongoID string, att model.Attachment) error {
	args := m.Called(ctx, mongoID, att)
	return args.Error(0)
}

func (m *AchievementRepositoryMock) FindByID(ctx context.Context, id string) (*model.Achievement, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Achievement), args.Error(1)
}

func (m *AchievementRepositoryMock) CountByType(ctx context.Context) (map[string]int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *AchievementRepositoryMock) FindAll(ctx context.Context) ([]model.Achievement, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Achievement), args.Error(1)
}

func (m *AchievementRepositoryMock) Update(ctx context.Context, id string, payload *model.Achievement) (*model.Achievement, error) {
	args := m.Called(ctx, id, payload)
	return args.Get(0).(*model.Achievement), args.Error(1)
}
