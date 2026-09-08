package domain

import (
	"context"
	"errors"

	"donor-service/app/internal/entity"

	"github.com/google/uuid"
)

var (
	ErrDonorMatchNotFound = errors.New("donor match not found")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidMatchStatus = errors.New("invalid donor match status")
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
}
