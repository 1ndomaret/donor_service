package domain

import (
	"context"
	"errors"
	"time"

	"donor-service/app/internal/entity"

	"github.com/google/uuid"
)

var (
	ErrDonationNotFound      = errors.New("donation not found")
	ErrInvalidDonationInput  = errors.New("invalid donation input")
	ErrInvalidDonationStatus = errors.New("invalid donation status")
	ErrDonorMatchNotAccepted = errors.New("donor match is not accepted")
)

type DonationReq struct {
	BloodRequestID uuid.UUID `json:"blood_request_id"`
	DonorMatchID   uuid.UUID `json:"donor_match_id"`
}

type DonationRepository interface {
	Create(
		ctx context.Context,
		donation *entity.Donation,
	) error

	FindOne(
		ctx context.Context,
		donation uuid.UUID,
	) (*entity.Donation, error)

	Completed(
		ctx context.Context,
		donationID uuid.UUID,
		confirmedBy uuid.UUID,
		donationDate time.Time,
	) error
}

type DonationUsecase interface {
	Create(
		requesterID uuid.UUID,
		req *DonationReq,
	) (*entity.Donation, error)

	Complete(
		requesterID uuid.UUID,
		donationID uuid.UUID,
	) (*entity.Donation, error)
}
