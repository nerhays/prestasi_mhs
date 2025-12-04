package service

import (
	"context"

	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository"
)

type AchievementService struct {
	achievementRepo repository.AchievementRepository
	studentRepo     repository.StudentRepository
	refRepo         repository.AchievementReferenceRepository
}

func NewAchievementService(
	achievementRepo repository.AchievementRepository,
	studentRepo repository.StudentRepository,
	refRepo repository.AchievementReferenceRepository,
) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
		studentRepo:     studentRepo,
		refRepo:         refRepo,
	}
}

// CreateAchievementForUser:
// - userID dari JWT
// - cari student by userID
// - insert achievement ke Mongo (pakai student.ID)
// - insert achievement_reference ke Postgres (status draft)
func (s *AchievementService) CreateAchievementForUser(ctx context.Context, userID string, ac *model.Achievement) (*model.Achievement, *model.AchievementReference, error) {
	// 1. Cari student berdasarkan userID
	student, err := s.studentRepo.FindByUserID(userID)
	if err != nil {
		return nil, nil, err
	}

	// 2. Set StudentID di dokumen Mongo = students.id (UUID)
	ac.StudentID = student.ID

	// 3. Insert ke Mongo
	createdAc, err := s.achievementRepo.Create(ctx, ac)
	if err != nil {
		return nil, nil, err
	}

	// 4. Insert reference ke Postgres (status: draft)
	ref, err := s.refRepo.CreateDraft(student.ID, createdAc.ID.Hex())
	if err != nil {
		return createdAc, nil, err // achievement sudah dibuat, tapi reference gagal
	}

	return createdAc, ref, nil
}

func (s *AchievementService) GetMyAchievements(ctx context.Context, userID string) ([]model.Achievement, error) {
	student, err := s.studentRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	return s.achievementRepo.FindByStudentID(ctx, student.ID)
}