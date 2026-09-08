package repository

import (
	"context"
	"errors"
	"time"

	"donor-service/app/internal/domain"
	"donor-service/app/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type donationRepository struct {
	db *gorm.DB
}

func NewDonationRepository(
	db *gorm.DB,
) domain.DonationRepository {
	return &donationRepository{
		db: db,
	}
}

func (r *donationRepository) Create(
	ctx context.Context,
	donation *entity.Donation,
) error {
	return r.db.
		WithContext(ctx).
		Create(donation).
		Error
}

func (r *donationRepository) FindOne(
	ctx context.Context,
	donationID uuid.UUID,
) (*entity.Donation, error) {
	var donation entity.Donation

	if err := r.db.
		WithContext(ctx).
		Where("id = ?", donationID).
		First(&donation).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDonationNotFound
		}

		return nil, err
	}

	return &donation, nil
}

func (r *donationRepository) Complete(
	ctx context.Context,
	donationID uuid.UUID,
	confirmedBy uuid.UUID,
	donationDate time.Time,
) error {

	result := r.db.
		WithContext(ctx).
		Model(&entity.Donation{}).
		Where("id = ?", donationID).
		Where("status = ?", "pending").
		Updates(map[string]interface{}{
			"status":        "completed",
			"donation_date": donationDate,
			"confirmed_by":  confirmedBy,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrInvalidDonationStatus
	}

	return nil
}
