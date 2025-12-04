package service

import (
	"context"
	"errors"
	"time"

	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository"
)
var (
	ErrRefNotFound          = errors.New("achievement_reference_not_found")
	ErrStudentProfileNotFound = errors.New("student_profile_not_found")
	ErrInvalidStatus        = errors.New("invalid_status_transition")
	ErrNotOwner             = errors.New("not_owner")
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

func (s *AchievementService) SubmitAchievement(ctx context.Context, userID, refID string) (*model.AchievementReference, error) {
	student, err := s.studentRepo.FindByUserID(userID)
	if err != nil {
		return nil, ErrStudentProfileNotFound
	}

	ref, err := s.refRepo.GetByID(refID)
	if err != nil {
		return nil, ErrRefNotFound
	}

	// milik mahasiswa yang login
	if ref.StudentID != student.ID {
		return nil, ErrNotOwner
	}

	if ref.Status != model.AchievementStatusDraft {
		return nil, ErrInvalidStatus
	}

	now := time.Now()
	ref.Status = model.AchievementStatusSubmitted
	ref.SubmittedAt = &now

	if err := s.refRepo.Save(ref); err != nil {
		return nil, err
	}

	return ref, nil
}

func (s *AchievementService) VerifyAchievement(ctx context.Context, verifierUserID, refID string) (*model.AchievementReference, error) {
	ref, err := s.refRepo.GetByID(refID)
	if err != nil {
		return nil, ErrRefNotFound
	}

	if ref.Status != model.AchievementStatusSubmitted {
		return nil, ErrInvalidStatus
	}

	now := time.Now()
	ref.Status = model.AchievementStatusVerified
	ref.VerifiedAt = &now
	ref.VerifiedBy = &verifierUserID
	ref.RejectionNote = nil

	if err := s.refRepo.Save(ref); err != nil {
		return nil, err
	}

	return ref, nil
}

func (s *AchievementService) RejectAchievement(ctx context.Context, verifierUserID, refID, note string) (*model.AchievementReference, error) {
	ref, err := s.refRepo.GetByID(refID)
	if err != nil {
		return nil, ErrRefNotFound
	}

	if ref.Status != model.AchievementStatusSubmitted {
		return nil, ErrInvalidStatus
	}

	now := time.Now()
	ref.Status = model.AchievementStatusRejected
	ref.VerifiedAt = &now
	ref.VerifiedBy = &verifierUserID
	ref.RejectionNote = &note

	if err := s.refRepo.Save(ref); err != nil {
		return nil, err
	}

	return ref, nil
}

func (s *AchievementService) DeleteDraftAchievement(ctx context.Context, userID, refID string) error {
	student, err := s.studentRepo.FindByUserID(userID)
	if err != nil {
		return ErrStudentProfileNotFound
	}

	ref, err := s.refRepo.GetByID(refID)
	if err != nil {
		return ErrRefNotFound
	}

	if ref.StudentID != student.ID {
		return ErrNotOwner
	}

	if ref.Status != model.AchievementStatusDraft {
		return ErrInvalidStatus
	}

	// 1. Soft delete di Mongo
	if err := s.achievementRepo.SoftDelete(ctx, ref.MongoAchievementID); err != nil {
		return err
	}

	// 2. Update status reference di Postgres
	ref.Status = model.AchievementStatusDeleted
	return s.refRepo.Save(ref)
}

