package usecase

import (
	"context"
	"time"

	"donor-service/app/internal/domain"
	"donor-service/app/internal/entity"

	"github.com/google/uuid"
)

type donorMatchesUsecase struct {
	donorMatchesRepository domain.DonorMatchesRepository
	bloodRequestReader     domain.BloodRequestReader
}

func NewDonorMatchesUsecase(
	donorMatchesRepository domain.DonorMatchesRepository,
	bloodRequestReader domain.BloodRequestReader,
) domain.DonorMatchesUsecase {
	return &donorMatchesUsecase{
		donorMatchesRepository: donorMatchesRepository,
		bloodRequestReader:     bloodRequestReader,
	}
}

var donorMatchesTimeout = 10 * time.Second

func (u *donorMatchesUsecase) Invite(
	requesterID uuid.UUID,
	bloodRequestID uuid.UUID,
	donorID uuid.UUID,
) (*entity.DonorMatches, error) {

	ctx, cancel := context.WithTimeout(
		context.TODO(),
		donorMatchesTimeout,
	)
	defer cancel()

	bloodRequest, err := u.bloodRequestReader.GetByID(
		ctx,
		bloodRequestID,
	)
	if err != nil {
		return nil, err
	}

	if bloodRequest.RequesterID != requesterID {
		return nil, domain.ErrForbidden
	}

	donorMatches, err := u.donorMatchesRepository.GetByBloodRequestAndDonor(
		ctx,
		bloodRequestID,
		donorID,
	)
	if err != nil {
		return nil, err
	}

	if donorMatches.Status == "accepted" ||
		donorMatches.Status == "declined" {
		return nil, domain.ErrInvalidMatchStatus
	}

	if donorMatches.Status == "invited" {
		return donorMatches, nil
	}

	if err := u.donorMatchesRepository.UpdateStatus(
		ctx,
		donorMatches.ID,
		"invited",
	); err != nil {
		return nil, err
	}

	donorMatches.Status = "invited"

	return donorMatches, nil
}

func (u *donorMatchesUsecase) Accept(
	donorID uuid.UUID,
	matchID uuid.UUID,
) (*entity.DonorMatches, error) {

	ctx, cancel := context.WithTimeout(
		context.TODO(),
		donorMatchesTimeout,
	)
	defer cancel()

	donorMatches, err := u.donorMatchesRepository.GetByID(
		ctx,
		matchID,
	)
	if err != nil {
		return nil, err
	}

	if donorMatches.DonorID != donorID {
		return nil, domain.ErrForbidden
	}

	if donorMatches.Status != "invited" {
		return nil, domain.ErrInvalidMatchStatus
	}

	if err := u.donorMatchesRepository.UpdateStatus(
		ctx,
		donorMatches.ID,
		"accepted",
	); err != nil {
		return nil, err
	}

	donorMatches.Status = "accepted"

	return donorMatches, nil
}

func (u *donorMatchesUsecase) Decline(
	donorID uuid.UUID,
	matchID uuid.UUID,
) (*entity.DonorMatches, error) {

	ctx, cancel := context.WithTimeout(
		context.TODO(),
		donorMatchesTimeout,
	)
	defer cancel()

	donorMatches, err := u.donorMatchesRepository.GetByID(
		ctx,
		matchID,
	)
	if err != nil {
		return nil, err
	}

	if donorMatches.DonorID != donorID {
		return nil, domain.ErrForbidden
	}

	if donorMatches.Status != "invited" {
		return nil, domain.ErrInvalidMatchStatus
	}

	if err := u.donorMatchesRepository.UpdateStatus(
		ctx,
		donorMatches.ID,
		"declined",
	); err != nil {
		return nil, err
	}

	donorMatches.Status = "declined"

	return donorMatches, nil
}
