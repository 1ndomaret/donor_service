package repository

import (
	"context"
	"donor-service/app/internal/domain"
	"donor-service/app/internal/entity"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type bloodRequestRepository struct {
	db *gorm.DB
}

func NewBloodRequestRepository(db *gorm.DB) domain.BloodRequestRepository {
	return &bloodRequestRepository{
		db: db,
	}
}

func (r *bloodRequestRepository) Create(ctx context.Context, req *entity.BloodRequest) error {
	return nil
}

func (r *bloodRequestRepository) FindAll(ctx context.Context) ([]entity.BloodRequest, error) {
	return nil, nil
}

func (r *bloodRequestRepository) FindOne(ctx context.Context, BloodRequestID uuid.UUID) (*entity.BloodRequest, error) {
	return nil, nil
}

func (r *bloodRequestRepository) Cancel(ctx context.Context, BloodRequestID uuid.UUID) error {
	return nil
}

func (r *bloodRequestRepository) FindMatches(ctx context.Context, BloodRequestID uuid.UUID) ([]entity.DonorMatches, error) {
	return nil, nil
}
