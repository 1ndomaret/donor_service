package domain

import (
	"context"
	"errors"

	"donor-service/app/internal/entity"

	"github.com/google/uuid"
)

var (
	ErrDonorMatchNotFound    = errors.New("donor match not found")
	ErrForbidden             = errors.New("forbidden")
	ErrInvalidMatchStatus    = errors.New("invalid donor match status")
	ErrBloodRequestFulfilled = errors.New("blood request quantity already fulfilled")
)

type DonorMatchRepository interface {
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.DonorMatch, error)

	GetByBloodRequestAndDonor(
		ctx context.Context,
		bloodRequestID uuid.UUID,
		donorID uuid.UUID,
	) (*entity.DonorMatch, error)

	UpdateStatus(
		ctx context.Context,
		id uuid.UUID,
		status string,
	) error

	Create(ctx context.Context, donor *entity.DonorMatch) error

	CountAccepted(
		ctx context.Context,
		bloodRequestID uuid.UUID,
	) (int64, error)

	UpdateDistance(
		ctx context.Context,
		id uuid.UUID,
		distanceKM float64,
	) error

	GetByDonorID(
		ctx context.Context,
		requesterID uuid.UUID,
	) ([]entity.DonorMatch, error)
}

type BloodRequestReader interface {
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.BloodRequest, error)
}

type DonorMatchUsecase interface {
	Invite(
		requesterID uuid.UUID,
		bloodRequestID uuid.UUID,
		donorID uuid.UUID,
	) (*entity.DonorMatch, error)

	Accept(
		donorID uuid.UUID,
		matchID uuid.UUID,
	) (*entity.DonorMatch, error)

	Decline(
		donorID uuid.UUID,
		matchID uuid.UUID,
	) (*entity.DonorMatch, error)

	GetByDonorID(
		requesterID uuid.UUID,
	) ([]entity.DonorMatch, error)
}
