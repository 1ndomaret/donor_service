package usecase

import (
	"context"
	"time"

	"donor-service/app/internal/domain"
	"donor-service/app/internal/entity"

	"github.com/google/uuid"
)

type donorMatchUsecase struct {
	donorMatchRepository domain.DonorMatchRepository
	bloodReqRepo         domain.BloodRequestRepository
	userServiceRepo      domain.UserServiceHttpRepo
}

func NewDonorMatchUsecase(
	donorMatchRepository domain.DonorMatchRepository,
	bloodReqRepo domain.BloodRequestRepository,
	userServiceRepo domain.UserServiceHttpRepo,
) domain.DonorMatchUsecase {
	return &donorMatchUsecase{
		donorMatchRepository: donorMatchRepository,
		bloodReqRepo:         bloodReqRepo,
		userServiceRepo:      userServiceRepo,
	}
}

var donorMatchTimeout = 10 * time.Second

func (u *donorMatchUsecase) Invite(
	requesterID uuid.UUID,
	bloodRequestID uuid.UUID,
	donorID uuid.UUID,
) (*entity.DonorMatch, error) {

	ctx, cancel := context.WithTimeout(
		context.TODO(),
		donorMatchTimeout,
	)
	defer cancel()

	bloodRequest, err := u.bloodReqRepo.FindOne(
		ctx,
		requesterID,
		bloodRequestID,
	)
	if err != nil {
		return nil, err
	}

	if bloodRequest.RequesterID != requesterID {
		return nil, domain.ErrForbidden
	}

	donorMatch, err := u.donorMatchRepository.GetByBloodRequestAndDonor(
		ctx,
		bloodRequestID,
		donorID,
	)
	if err != nil {
		return nil, err
	}

	if donorMatch.Status == "accepted" ||
		donorMatch.Status == "declined" {
		return nil, domain.ErrInvalidMatchStatus
	}

	if donorMatch.Status == "invited" {
		return donorMatch, nil
	}

	if err := u.donorMatchRepository.UpdateStatus(
		ctx,
		donorMatch.ID,
		"invited",
	); err != nil {
		return nil, err
	}

	donorMatch.Status = "invited"

	return donorMatch, nil
}

func (u *donorMatchUsecase) Accept(
	donorID uuid.UUID,
	matchID uuid.UUID,
) (*entity.DonorMatch, error) {

	ctx, cancel := context.WithTimeout(
		context.TODO(),
		donorMatchTimeout,
	)
	defer cancel()

	donorMatch, err := u.donorMatchRepository.GetByID(
		ctx,
		matchID,
	)
	if err != nil {
		return nil, err
	}

	if donorMatch.DonorID != donorID {
		return nil, domain.ErrForbidden
	}

	if donorMatch.Status != "invited" {
		return nil, domain.ErrInvalidMatchStatus
	}

	if err := u.donorMatchRepository.UpdateStatus(
		ctx,
		donorMatch.ID,
		"accepted",
	); err != nil {
		return nil, err
	}

	donorMatch.Status = "accepted"

	return donorMatch, nil
}

func (u *donorMatchUsecase) Decline(
	donorID uuid.UUID,
	matchID uuid.UUID,
) (*entity.DonorMatch, error) {

	ctx, cancel := context.WithTimeout(
		context.TODO(),
		donorMatchTimeout,
	)
	defer cancel()

	donorMatch, err := u.donorMatchRepository.GetByID(
		ctx,
		matchID,
	)
	if err != nil {
		return nil, err
	}

	if donorMatch.DonorID != donorID {
		return nil, domain.ErrForbidden
	}

	if donorMatch.Status != "invited" {
		return nil, domain.ErrInvalidMatchStatus
	}

	if err := u.donorMatchRepository.UpdateStatus(
		ctx,
		donorMatch.ID,
		"declined",
	); err != nil {
		return nil, err
	}

	donorMatch.Status = "declined"

	return donorMatch, nil
}
