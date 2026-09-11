package usecase

import (
	"context"
	"fmt"
	"time"

	"donor-service/app/internal/domain"
	"donor-service/app/internal/dto"
	"donor-service/app/internal/entity"

	"github.com/google/uuid"
)

type donorMatchUsecase struct {
	donorMatchRepository domain.DonorMatchRepository
	bloodReqRepo         domain.BloodRequestRepository
	userServiceRepo      domain.UserServiceHttpRepo
	geoapifyRepo         domain.GeoapifyRepository
	donationRepository   domain.DonationRepository
}

func NewDonorMatchUsecase(
	donorMatchRepository domain.DonorMatchRepository,
	bloodReqRepo domain.BloodRequestRepository,
	userServiceRepo domain.UserServiceHttpRepo,
	geoapifyRepo domain.GeoapifyRepository,
	donationRepository domain.DonationRepository,
) domain.DonorMatchUsecase {
	return &donorMatchUsecase{
		donorMatchRepository: donorMatchRepository,
		bloodReqRepo:         bloodReqRepo,
		userServiceRepo:      userServiceRepo,
		geoapifyRepo:         geoapifyRepo,
		donationRepository:   donationRepository,
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

	bloodRequest, err := u.bloodReqRepo.GetById(
		ctx,
		donorMatch.BloodRequestID,
	)
	if err != nil {
		return nil, err
	}

	acceptedCount, err := u.donorMatchRepository.CountAccepted(
		ctx,
		donorMatch.BloodRequestID,
	)
	if err != nil {
		return nil, err
	}

	if acceptedCount >= int64(bloodRequest.Quantity) {
		return nil, domain.ErrBloodRequestFulfilled
	}

	donorProfile, err := u.userServiceRepo.GetDonorProfile(
		ctx,
		donorMatch.DonorID,
	)
	if err != nil {
		return nil, err
	}

	routeReq := dto.GeoapifyRoutingRequest{
		OriginLat:      donorProfile.Latitude,
		OriginLon:      donorProfile.Longitude,
		DestinationLat: bloodRequest.Latitude,
		DestinationLon: bloodRequest.Longitude,
	}

	route, err := u.geoapifyRepo.GetGeoapifyRoute(
		ctx,
		&routeReq,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	distanceKM := route.Distance

	if route.DistanceUnits == "meters" {
		distanceKM = route.Distance / 1000
	}

	if err := u.donorMatchRepository.UpdateDistance(
		ctx,
		donorMatch.ID,
		distanceKM,
	); err != nil {
		return nil, err
	}

	donorMatch.DistanceKM = distanceKM

	if err := u.donorMatchRepository.UpdateStatus(
		ctx,
		donorMatch.ID,
		"accepted",
	); err != nil {
		return nil, err
	}

	donorMatch.Status = "accepted"

	donation := &entity.Donation{
		ID:             uuid.New(),
		BloodRequestID: donorMatch.BloodRequestID,
		DonorID:        donorMatch.DonorID,
		DonorMatchID:   donorMatch.ID,
		Status:         "pending",
	}

	if err := u.donationRepository.Create(
		ctx,
		donation,
	); err != nil {
		return nil, err
	}

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

func (u *donorMatchUsecase) GetByDonorID(
	donorID uuid.UUID,
) ([]entity.DonorMatch, error) {
	ctx, cancel := context.WithTimeout(
		context.TODO(),
		donorMatchTimeout,
	)
	defer cancel()

	return u.donorMatchRepository.GetByDonorID(
		ctx,
		donorID,
	)
}
