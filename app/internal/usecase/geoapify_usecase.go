package usecase

import (
	"context"
	"donor-service/app/internal/domain"
	"donor-service/app/internal/dto"
	"time"

	"github.com/google/uuid"
)

type geoapifyUsecase struct {
	geoapifyRepo    domain.GeoapifyRepository
	donorMatchRepo  domain.DonorMatchRepository
	userServiceRepo domain.UserServiceHttpRepo
	bloodReqRepo    domain.BloodRequestRepository
}

func NewGeoapifyUsecase(
	geoapifyRepo domain.GeoapifyRepository,
	donorMatchRepo domain.DonorMatchRepository,
	userServiceRepo domain.UserServiceHttpRepo,
	bloodReqRepo domain.BloodRequestRepository,
) domain.GeoapifyUsecase {
	return &geoapifyUsecase{
		geoapifyRepo:    geoapifyRepo,
		donorMatchRepo:  donorMatchRepo,
		userServiceRepo: userServiceRepo,
		bloodReqRepo:    bloodReqRepo,
	}
}

var geoapifyTimeout = 10 * time.Second

func (u *geoapifyUsecase) GetGeoapifyHospitals(
	city string,
) ([]domain.GeoapifyHospital, error) {

	ctx, cancel := context.WithTimeout(
		context.TODO(),
		geoapifyTimeout,
	)
	defer cancel()

	return u.geoapifyRepo.GetGeoapifyHospitals(ctx, city)
}

func (u *geoapifyUsecase) GetGeoapifyRoute(donorMatchID uuid.UUID) (*domain.GeoapifyRoute, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), geoapifyTimeout)
	defer cancel()

	match, err := u.donorMatchRepo.GetByID(ctx, donorMatchID)
	if err != nil {
		return nil, err
	}

	bloodReq, err := u.bloodReqRepo.GetById(ctx, match.BloodRequestID)
	if err != nil {
		return nil, err
	}

	donor, err := u.userServiceRepo.GetDonorProfile(ctx, match.DonorID)
	if err != nil {
		return nil, err
	}

	req := dto.GeoapifyRoutingRequest{
		OriginLat:      donor.Latitude,
		OriginLon:      donor.Longitude,
		DestinationLat: bloodReq.Latitude,
		DestinationLon: bloodReq.Longitude,
	}

	return u.geoapifyRepo.GetGeoapifyRoute(ctx, &req)
}
