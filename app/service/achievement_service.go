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
	ErrUserNotFound     = errors.New("verifier user not found")
    ErrNotAdvisor       = errors.New("only the assigned academic advisor can verify or reject this achievement")
	ErrStudentNoAdvisor        = errors.New("student has no advisor assigned")
	ErrLecturerNotFound        = errors.New("lecturer record not found")
)

type AchievementService struct {
	achievementRepo repository.AchievementRepository
	studentRepo     repository.StudentRepository
	refRepo         repository.AchievementReferenceRepository
	userRepo        repository.UserRepository
	lecturerRepo    repository.LecturerRepository
}

func NewAchievementService(
	achievementRepo repository.AchievementRepository,
	studentRepo repository.StudentRepository,
	refRepo repository.AchievementReferenceRepository,
	userRepo repository.UserRepository,
	lecturerRepo    repository.LecturerRepository,
) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
		studentRepo:     studentRepo,
		refRepo:         refRepo,
		userRepo:        userRepo,
		lecturerRepo:    lecturerRepo,
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

    // Ambil student (menggunakan student.id dari ref)
	student, err := s.studentRepo.FindByID(ref.StudentID)
	if err != nil {
		return nil, ErrStudentProfileNotFound
	}

	// Pastikan student punya advisor
	if student.AdvisorID == "" {
		return nil, ErrStudentNoAdvisor
	}

	// Ambil lecturer record (students.advisor_id menyimpan lecturers.id)
	lect, err := s.lecturerRepo.FindByID(student.AdvisorID)
	if err != nil {
		return nil, ErrLecturerNotFound
	}

	// Ambil verifier user
	verifier, err := s.userRepo.FindByID(verifierUserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Jika verifier role = Dosen Wali, hanya boleh verifikasi jika lecturer.user_id == verifier.ID
	if verifier.Role.Name == "Dosen Wali" {
    if lect.UserID != verifier.ID {
        return nil, ErrNotAdvisor
    }
}
	// Admin bypass; other roles will be blocked by middleware/route

	// Cek status
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
	// Ambil student
	student, err := s.studentRepo.FindByID(ref.StudentID)
	if err != nil {
		return nil, ErrStudentProfileNotFound
	}

	if student.AdvisorID == "" {
		return nil, ErrStudentNoAdvisor
	}

	lect, err := s.lecturerRepo.FindByID(student.AdvisorID)
	if err != nil {
		return nil, ErrLecturerNotFound
	}

	verifier, err := s.userRepo.FindByID(verifierUserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if verifier.Role.Name == "Dosen Wali" {
    if lect.UserID != verifier.ID {
        return nil, ErrNotAdvisor
    }
}

	if ref.Status != model.AchievementStatusSubmitted {
		return nil, ErrInvalidStatus
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

func (s *AchievementService) GetDeletedAchievements(ctx context.Context, userID string) ([]model.Achievement, error) {
    student, err := s.studentRepo.FindByUserID(userID)
    if err != nil {
        return nil, ErrStudentProfileNotFound
    }

    return s.achievementRepo.FindDeletedByStudentID(ctx, student.ID)
}

func (s *AchievementService) GetBimbinganAchievements(ctx context.Context, verifierUserID string, page, perPage int, status *model.AchievementStatus) (int64, []map[string]interface{}, error) {
    // 1. find lecturer by user id
    lect, err := s.lecturerRepo.FindByUserID(verifierUserID) // if you don't have this, use FindByUserID; else FindByID after mapping
    if err != nil {
        return 0, nil, err
    }

    // 2. get students by advisor_id = lect.ID
    students, err := s.studentRepo.FindByAdvisorLecturerID(lect.ID)
    if err != nil {
        return 0, nil, err
    }
    if len(students) == 0 {
        return 0, []map[string]interface{}{}, nil
    }
    studentIDs := make([]string, 0, len(students))
    for _, sct := range students {
        studentIDs = append(studentIDs, sct.ID)
    }

    // 3. count total refs
    total, err := s.refRepo.CountByStudentIDs(studentIDs, status)
    if err != nil {
        return 0, nil, err
    }

    // 4. find refs paginated
    if page < 1 { page = 1 }
    if perPage < 1 { perPage = 10 }
    if perPage > 100 { perPage = 100 }
    offset := (page - 1) * perPage
    refs, err := s.refRepo.FindByStudentIDs(studentIDs, status, perPage, offset)
    if err != nil {
        return 0, nil, err
    }

    // 5. collect mongo ids
    mongoIDs := make([]string, 0, len(refs))
    for _, r := range refs {
        mongoIDs = append(mongoIDs, r.MongoAchievementID)
    }

    // 6. fetch mongo docs
    achievements, err := s.achievementRepo.FindByIDs(ctx, mongoIDs)
    if err != nil {
        return 0, nil, err
    }

    // 7. index achievements by id (hex)
    achMap := map[string]model.Achievement{}
    for _, a := range achievements {
        // a.ID is primitive.ObjectID; convert to hex string
        achMap[a.ID.Hex()] = a
    }

    // 8. combine refs + achievement doc into result rows
    results := make([]map[string]interface{}, 0, len(refs))
    for _, r := range refs {
        item := map[string]interface{}{"reference": r}
        if a, ok := achMap[r.MongoAchievementID]; ok {
            item["achievement"] = a
        } else {
            item["achievement"] = nil
        }
        results = append(results, item)
    }

    return total, results, nil
}
