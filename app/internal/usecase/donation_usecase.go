package usecase

import (
	"context"
	"time"

	"donor-service/app/internal/domain"
	"donor-service/app/internal/entity"

	"github.com/google/uuid"
)

type donationUsecase struct {
	donationRepo     domain.DonationRepository
	donorMatchesRepo domain.DonorMatchesRepository
	bloodReqRepo     domain.BloodRequestRepository
}

func NewDonationUsecase(
	donationRepo domain.DonationRepository,
	donorMatchesRepo domain.DonorMatchesRepository,
	bloodReqRepo domain.BloodRequestRepository,
) domain.DonationUsecase {
	return &donationUsecase{
		donationRepo:     donationRepo,
		donorMatchesRepo: donorMatchesRepo,
		bloodReqRepo:     bloodReqRepo,
	}
}

var donationTimeout = 10 * time.Second

func (u *donationUsecase) Create(
	requesterID uuid.UUID,
	req *domain.DonationReq,
) (*entity.Donation, error) {

	if req == nil ||
		req.BloodRequestID == uuid.Nil ||
		req.DonorMatchID == uuid.Nil {
		return nil, domain.ErrInvalidDonationInput
	}

	ctx, cancel := context.WithTimeout(
		context.TODO(),
		donationTimeout,
	)
	defer cancel()

	_, err := u.bloodReqRepo.FindOne(
		ctx,
		requesterID,
		req.BloodRequestID,
	)
	if err != nil {
		return nil, err
	}

	donorMatch, err := u.donorMatchesRepo.GetByID(
		ctx,
		req.DonorMatchID,
	)
	if err != nil {
		return nil, err
	}

	if donorMatch.BloodRequestID != req.BloodRequestID {
		return nil, domain.ErrInvalidDonationInput
	}

	if donorMatch.Status != "accepted" {
		return nil, domain.ErrDonorMatchNotAccepted
	}

	donation := &entity.Donation{
		ID:             uuid.New(),
		BloodRequestID: req.BloodRequestID,
		DonorID:        donorMatch.DonorID,
		DonorMatchID:   donorMatch.ID,
		Status:         "pending",
	}

	if err := u.donationRepo.Create(
		ctx,
		donation,
	); err != nil {
		return nil, err
	}

	return donation, nil
}

func (u *donationUsecase) Complete(
	requesterID uuid.UUID,
	donationID uuid.UUID,
) (*entity.Donation, error) {

	ctx, cancel := context.WithTimeout(
		context.TODO(),
		donationTimeout,
	)
	defer cancel()

	donation, err := u.donationRepo.FindOne(
		ctx,
		donationID,
	)
	if err != nil {
		return nil, err
	}

	_, err = u.bloodReqRepo.FindOne(
		ctx,
		requesterID,
		donation.BloodRequestID,
	)
	if err != nil {
		return nil, err
	}

	if donation.Status != "pending" {
		return nil, domain.ErrInvalidDonationStatus
	}

	now := time.Now()

	if err := u.donationRepo.Complete(
		ctx,
		donation.ID,
		requesterID,
		now,
	); err != nil {
		return nil, err
	}

	donation.Status = "completed"
	donation.DonationDate = &now
	donation.ConfirmedBy = &requesterID

	return donation, nil
}
