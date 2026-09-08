package usecase

import (
	"donor-service/app/internal/domain"
	"donor-service/app/internal/entity"

	"github.com/google/uuid"
)

type bloodReqUsecase struct {
	bloodReqRepo domain.BloodRequestRepository
}

func NewBloodRequestUsecase(bloodReqRepo domain.BloodRequestRepository) domain.BloodRequestUsecase {
	return &bloodReqUsecase{
		bloodReqRepo: bloodReqRepo,
	}
}

func (u *bloodReqUsecase) Create(req *domain.BloodRequestReq) (*entity.BloodRequest, error) {
	return nil, nil
}

func (u *bloodReqUsecase) FindAll() ([]entity.BloodRequest, error) {
	return nil, nil
}

func (u *bloodReqUsecase) FindOne(BloodRequestID uuid.UUID) (*entity.BloodRequest, error) {
	return nil, nil
}

func (u *bloodReqUsecase) Cancel(BloodRequestID uuid.UUID) error {
	return nil
}

func (u *bloodReqUsecase) FindMatches(BloodRequestID uuid.UUID) ([]entity.DonorMatches, error) {
	return nil, nil
}
