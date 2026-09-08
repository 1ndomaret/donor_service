package domain

import (
	"context"
	"errors"

	"donor-service/app/internal/entity"

	"github.com/google/uuid"
)

var (
	ErrDonorMatchesNotFound = errors.New("donor matches not found")
	ErrForbidden            = errors.New("forbidden")
	ErrInvalidMatchStatus   = errors.New("invalid donor match status")
)

type DonorMatchesRepository interface {
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.DonorMatches, error)

	GetByBloodRequestAndDonor(
		ctx context.Context,
		bloodRequestID uuid.UUID,
		donorID uuid.UUID,
	) (*entity.DonorMatches, error)

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

type DonorMatchesUsecase interface {
	Invite(
		requesterID uuid.UUID,
		bloodRequestID uuid.UUID,
		donorID uuid.UUID,
	) (*entity.DonorMatches, error)

	Accept(
		donorID uuid.UUID,
		matchID uuid.UUID,
	) (*entity.DonorMatches, error)

	Decline(
		donorID uuid.UUID,
		matchID uuid.UUID,
	) (*entity.DonorMatches, error)
}
