package service

import (
	"context"
	"testing"

	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateAchievementForUser_Success(t *testing.T) {
	studentRepo := new(mocks.StudentRepositoryMock)
	achRepo := new(mocks.AchievementRepositoryMock)
	refRepo := new(mocks.AchievementReferenceRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	lectRepo := new(mocks.LecturerRepositoryMock)
	logRepo := new(mocks.AchievementStatusLogRepositoryMock)

	svc := NewAchievementService(
		achRepo,
		studentRepo,
		refRepo,
		userRepo,
		lectRepo,
		logRepo,
	)

	userID := "user-1"
	student := &model.Student{ID: "student-1"}

	studentRepo.On("FindByUserID", userID).Return(student, nil)

	ach := &model.Achievement{
		Details: map[string]interface{}{
			"eventDate": "2025-01-01",
		},
	}

	achRepo.
		On("Create", context.Background(), ach).
		Run(func(args mock.Arguments) {
			a := args.Get(1).(*model.Achievement)
			a.ID = primitive.NewObjectID()
		}).
		Return(ach, nil)

	ref := &model.AchievementReference{ID: "ref-1"}

	refRepo.
		On("CreateDraft", student.ID, mock.Anything).
		Return(ref, nil)

	logRepo.
		On("Create", mock.Anything).
		Return(nil)

	ac, r, err := svc.CreateAchievementForUser(context.Background(), userID, ach)

	assert.NoError(t, err)
	assert.NotNil(t, ac)
	assert.NotNil(t, r)
	assert.Equal(t, student.ID, ac.StudentID)
}
func TestSubmitAchievement_Success(t *testing.T) {
	studentRepo := new(mocks.StudentRepositoryMock)
	refRepo := new(mocks.AchievementReferenceRepositoryMock)
	achRepo := new(mocks.AchievementRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	lectRepo := new(mocks.LecturerRepositoryMock)
	logRepo := new(mocks.AchievementStatusLogRepositoryMock)

	svc := NewAchievementService(
		achRepo, studentRepo, refRepo, userRepo, lectRepo, logRepo,
	)

	userID := "user-1"
	refID := "ref-1"

	student := &model.Student{ID: "student-1"}
	ref := &model.AchievementReference{
		ID:        refID,
		StudentID: student.ID,
		Status:    model.AchievementStatusDraft,
	}

	studentRepo.On("FindByUserID", userID).Return(student, nil)
	refRepo.On("GetByID", refID).Return(ref, nil)
	refRepo.On("Save", ref).Return(nil)

	logRepo.On(
		"Create",
		mock.AnythingOfType("*model.AchievementStatusLog"),
	).Return(nil)

	updated, err := svc.SubmitAchievement(context.Background(), userID, refID)

	assert.NoError(t, err)
	assert.Equal(t, model.AchievementStatusSubmitted, updated.Status)
}
func TestVerifyAchievement_ByAdvisor_Success(t *testing.T) {
	refRepo := new(mocks.AchievementReferenceRepositoryMock)
	studentRepo := new(mocks.StudentRepositoryMock)
	lectRepo := new(mocks.LecturerRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	achRepo := new(mocks.AchievementRepositoryMock)
	logRepo := new(mocks.AchievementStatusLogRepositoryMock)

	svc := NewAchievementService(
		achRepo, studentRepo, refRepo, userRepo, lectRepo, logRepo,
	)

	ref := &model.AchievementReference{
		ID:        "ref-1",
		Status:    model.AchievementStatusSubmitted,
		StudentID: "student-1",
	}

	student := &model.Student{
		ID:        "student-1",
		AdvisorID: "lect-1",
	}

	lect := &model.Lecturer{
		ID:     "lect-1",
		UserID: "user-lect",
	}

	verifier := &model.User{
		ID: "user-lect",
		Role: model.Role{
			Name: "Dosen Wali",
		},
	}

	refRepo.On("GetByID", ref.ID).Return(ref, nil)
	studentRepo.On("FindByID", ref.StudentID).Return(student, nil)
	lectRepo.On("FindByID", student.AdvisorID).Return(lect, nil)
	userRepo.On("FindByID", verifier.ID).Return(verifier, nil)
	refRepo.On("Save", ref).Return(nil)

	logRepo.On(
		"Create",
		mock.AnythingOfType("*model.AchievementStatusLog"),
	).Return(nil)
	
	updated, err := svc.VerifyAchievement(context.Background(), verifier.ID, ref.ID)

	assert.NoError(t, err)
	assert.Equal(t, model.AchievementStatusVerified, updated.Status)
}
