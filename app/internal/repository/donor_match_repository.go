package repository

import (
	"context"
	"errors"

	"donor-service/app/internal/domain"
	"donor-service/app/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type donorMatchRepository struct {
	db *gorm.DB
}

func NewDonorMatchRepository(db *gorm.DB) domain.DonorMatchRepository {
	return &donorMatchRepository{
		db: db,
	}
}

func (r *donorMatchRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.DonorMatch, error) {
	var donorMatch entity.DonorMatch

	if err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&donorMatch).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDonorMatchNotFound
		}

		return nil, err
	}

	return &donorMatch, nil
}

func (r *donorMatchRepository) GetByBloodRequestAndDonor(
	ctx context.Context,
	bloodRequestID uuid.UUID,
	donorID uuid.UUID,
) (*entity.DonorMatch, error) {
	var donorMatch entity.DonorMatch

	if err := r.db.
		WithContext(ctx).
		Where(
			"blood_request_id = ? AND donor_id = ?",
			bloodRequestID,
			donorID,
		).
		First(&donorMatch).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDonorMatchNotFound
		}

		return nil, err
	}

	return &donorMatch, nil
}

func (r *donorMatchRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) error {
	return r.db.
		WithContext(ctx).
		Model(&entity.DonorMatch{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

func (r *donorMatchRepository) Create(ctx context.Context, donor *entity.DonorMatch) error {
	return r.db.Create(donor).Error
}

func (r *donorMatchRepository) GetByRequesterID(
	ctx context.Context,
	requesterID uuid.UUID,
) ([]entity.DonorMatch, error) {
	var donorMatches []entity.DonorMatch

	subQuery := r.db.
		WithContext(ctx).
		Model(&entity.BloodRequest{}).
		Select("id").
		Where("requester_id = ?", requesterID)

	err := r.db.
		WithContext(ctx).
		Preload("BloodRequest").
		Where("blood_request_id IN (?)", subQuery).
		Find(&donorMatches).
		Error

	if err != nil {
		return nil, err
	}

	return donorMatches, nil
}

func (r *donorMatchRepository) CountAccepted(
	ctx context.Context,
	bloodRequestID uuid.UUID,
) (int64, error) {
	var count int64

	err := r.db.
		WithContext(ctx).
		Model(&entity.DonorMatch{}).
		Where(
			"blood_request_id = ? AND status = ?",
			bloodRequestID,
			"accepted",
		).
		Count(&count).
		Error

	return count, err
}

func (r *donorMatchRepository) UpdateDistance(
	ctx context.Context,
	id uuid.UUID,
	distanceKM float64,
) error {
	return r.db.
		WithContext(ctx).
		Model(&entity.DonorMatch{}).
		Where("id = ?", id).
		Update("distance_km", distanceKM).
		Error
}
