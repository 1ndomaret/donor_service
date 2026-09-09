package usecase

import (
	"context"
	"donor-service/app/internal/domain"
	"time"
)

type hospitalUsecase struct {
	hospitalRepo domain.HospitalRepository
}

func NewHospitalUsecase(
	hospitalRepo domain.HospitalRepository,
) domain.HospitalUsecase {
	return &hospitalUsecase{
		hospitalRepo: hospitalRepo,
	}
}

var hospitalTimeout = 10 * time.Second

func (u *hospitalUsecase) GetHospitals(
	city string,
) ([]domain.Hospital, error) {

	ctx, cancel := context.WithTimeout(
		context.TODO(),
		hospitalTimeout,
	)
	defer cancel()

	return u.hospitalRepo.GetHospitals(ctx, city)
}
