package repository

import (
	"context"
	"errors"

	"donor-service/app/internal/domain"
	"donor-service/app/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type donorMatchesRepository struct {
	db *gorm.DB
}

func NewDonorMatchesRepository(db *gorm.DB) domain.DonorMatchesRepository {
	return &donorMatchesRepository{
		db: db,
	}
}

func (r *donorMatchesRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.DonorMatches, error) {
	var donorMatches entity.DonorMatches

	if err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&donorMatches).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDonorMatchesNotFound
		}

		return nil, err
	}

	return &donorMatches, nil
}

func (r *donorMatchesRepository) GetByBloodRequestAndDonor(
	ctx context.Context,
	bloodRequestID uuid.UUID,
	donorID uuid.UUID,
) (*entity.DonorMatches, error) {
	var donorMatches entity.DonorMatches

	if err := r.db.
		WithContext(ctx).
		Where(
			"blood_request_id = ? AND donor_id = ?",
			bloodRequestID,
			donorID,
		).
		First(&donorMatches).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDonorMatchesNotFound
		}

		return nil, err
	}

	return &donorMatches, nil
}

func (r *donorMatchesRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) error {
	return r.db.
		WithContext(ctx).
		Model(&entity.DonorMatches{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}
