package usecase

import (
	"context"
	"donor-service/app/internal/domain"
	"time"
)

type geoapifyUsecase struct {
	geoapifyRepo domain.GeoapifyRepository
}

func NewGeoapifyUsecase(
	geoapifyRepo domain.GeoapifyRepository,
) domain.GeoapifyUsecase {
	return &geoapifyUsecase{
		geoapifyRepo: geoapifyRepo,
	}
}

var geoapifyTimeout = 10 * time.Second

func (u *geoapifyUsecase) GetGeoapifys(
	city string,
) ([]domain.Geoapify, error) {

	ctx, cancel := context.WithTimeout(
		context.TODO(),
		geoapifyTimeout,
	)
	defer cancel()

	return u.geoapifyRepo.GetGeoapifys(ctx, city)
}
