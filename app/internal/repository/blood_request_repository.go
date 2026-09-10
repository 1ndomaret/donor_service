package repository

import (
	"context"
	"donor-service/app/internal/domain"
	"donor-service/app/internal/entity"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type bloodRequestRepository struct {
	db *gorm.DB
}

func NewBloodRequestRepository(db *gorm.DB) domain.BloodRequestRepository {
	return &bloodRequestRepository{
		db: db,
	}
}

func (r *bloodRequestRepository) Create(ctx context.Context, bloodReq *entity.BloodRequest) error {
	return r.db.Create(bloodReq).Error
}

func (r *bloodRequestRepository) FindAll(ctx context.Context, userID uuid.UUID) ([]entity.BloodRequest, error) {
	var bloodRequests []entity.BloodRequest

	err := r.db.WithContext(ctx).Where("requester_id = ?", userID).Find(&bloodRequests).Error
	if err != nil {
		return nil, err
	}

	return bloodRequests, nil
}

func (r *bloodRequestRepository) FindOne(ctx context.Context, userID uuid.UUID, bloodRequestID uuid.UUID) (*entity.BloodRequest, error) {
	var bloodRequest entity.BloodRequest

	if err := r.db.WithContext(ctx).Where("id = ? AND requester_id = ?", bloodRequestID, userID).First(&bloodRequest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBloodReqNotFound
		}
		return nil, err
	}

	return &bloodRequest, nil
}

func (r *bloodRequestRepository) Cancel(ctx context.Context, userID uuid.UUID, bloodRequestID uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&entity.BloodRequest{}).
		Where("id = ? AND requester_id = ?", bloodRequestID, userID).
		Where("status NOT IN ?", []string{"cancelled", "completed"}).
		Update("status", "cancelled")

	if result.Error != nil {
		fmt.Println(result.Error.Error())

		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrBloodReqNotFound
	}

	return nil
}

func (r *bloodRequestRepository) GetMatches(ctx context.Context, bloodRequestID uuid.UUID) ([]entity.DonorMatch, error) {
	var donorMatch []entity.DonorMatch

	err := r.db.WithContext(ctx).Where("blood_request_id = ?", bloodRequestID).Find(&donorMatch).Error
	if err != nil {
		return nil, err
	}

	return donorMatch, nil
}

func (r *bloodRequestRepository) GetById(ctx context.Context, bloodRequestID uuid.UUID) (*entity.BloodRequest, error) {
	var bloodRequest entity.BloodRequest

	if err := r.db.WithContext(ctx).Where("id = ?", bloodRequestID).First(&bloodRequest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBloodReqNotFound
		}
		return nil, err
	}

	return &bloodRequest, nil
}

func (r *bloodRequestRepository) CreateMatches(ctx context.Context, reqID uuid.UUID, donors []entity.DonorProfile) error {
	var matches []entity.DonorMatch

	for _, donor := range donors {
		matches = append(matches, entity.DonorMatch{
			ID:             uuid.New(),
			BloodRequestID: reqID,
			DonorID:        donor.ID,
			Status:         "invited",
			CreatedAt:      time.Now(),
		})
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&matches).Error
}

func (r *bloodRequestRepository) GetPendingReqs(ctx context.Context) ([]entity.BloodRequest, error) {
	var reqs []entity.BloodRequest

	if err := r.db.WithContext(ctx).Where("status = ?", "pending").Find(&reqs).Error; err != nil {
		return nil, err
	}

	return reqs, nil
}
