package usecase

import (
	"context"
	"donor-service/app/internal/domain"
	"donor-service/app/internal/entity"
	"donor-service/app/internal/helper"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type bloodReqUsecase struct {
	bloodReqRepo    domain.BloodRequestRepository
	userServiceRepo domain.UserServiceHttpRepo
}

func NewBloodRequestUsecase(bloodReqRepo domain.BloodRequestRepository, userServiceRepo domain.UserServiceHttpRepo) domain.BloodRequestUsecase {
	return &bloodReqUsecase{
		bloodReqRepo:    bloodReqRepo,
		userServiceRepo: userServiceRepo,
	}
}

var timeOut = 10 * time.Second
var cronTimeOut = 5 * time.Minute
var serviceTimeOut = 15 * time.Second

func (u *bloodReqUsecase) Create(userID uuid.UUID, req *domain.BloodRequestReq) (*entity.BloodRequest, error) {
	if req.BloodType == "" ||
		req.Quantity <= 0 ||
		req.Urgency == "" ||
		req.HospitalExternalID == "" ||
		req.HospitalName == "" ||
		req.City == "" ||
		req.NeededAt.IsZero() {
		return nil, domain.ErrInvalidInput
	}

	if !helper.IsValidBloodType(req.BloodType) {
		return nil, domain.ErrInvalidBloodType
	}

	if !helper.IsValidCoord(req.Latitude, req.Longitude) {
		return nil, domain.ErrInvalidCoord
	}

	if req.NeededAt.Before(time.Now()) {
		return nil, domain.ErrInvalidNeededAt
	}

	bloodReq := &entity.BloodRequest{
		ID:                 uuid.New(),
		RequesterID:        userID,
		BloodType:          req.BloodType,
		Quantity:           req.Quantity,
		Urgency:            req.Urgency,
		HospitalExternalID: req.HospitalExternalID,
		HospitalName:       req.HospitalName,
		City:               req.City,
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		Notes:              req.Notes,
		Status:             "pending",
		NeededAt:           req.NeededAt,
	}

	ctx, cancel := context.WithTimeout(context.TODO(), timeOut)
	defer cancel()

	if err := u.bloodReqRepo.Create(ctx, bloodReq); err != nil {
		return nil, err
	}

	filter := domain.SearchMatchesRequest{
		BloodType: bloodReq.BloodType,
		City:      bloodReq.City,
	}

	matches, err := u.userServiceRepo.SearchMatches(
		ctx,
		"",
		&filter,
	)
	if err != nil {
		log.Println("failed to search donor matches:", err)
		return bloodReq, nil
	}

	if len(matches) > 0 {
		if err := u.bloodReqRepo.CreateMatches(
			ctx,
			bloodReq.ID,
			matches,
		); err != nil {
			return nil, err
		}
	}

	return bloodReq, nil
}

func (u *bloodReqUsecase) FindAll(userID uuid.UUID) ([]entity.BloodRequest, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeOut)
	defer cancel()

	return u.bloodReqRepo.FindAll(ctx, userID)
}

func (u *bloodReqUsecase) FindOne(userID uuid.UUID, BloodRequestID uuid.UUID) (*entity.BloodRequest, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeOut)
	defer cancel()

	return u.bloodReqRepo.FindOne(ctx, userID, BloodRequestID)
}

func (u *bloodReqUsecase) Cancel(userID uuid.UUID, BloodRequestID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.TODO(), timeOut)
	defer cancel()

	return u.bloodReqRepo.Cancel(ctx, userID, BloodRequestID)
}

func (u *bloodReqUsecase) GetMatches(BloodRequestID uuid.UUID) ([]entity.DonorMatch, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeOut)
	defer cancel()

	_, err := u.bloodReqRepo.GetById(ctx, BloodRequestID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *bloodReqUsecase) SearchMatches(token string, userID uuid.UUID, BloodRequestID uuid.UUID) ([]entity.DonorProfile, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), serviceTimeOut)
	defer cancel()

	bloodReq, err := u.bloodReqRepo.FindOne(ctx, userID, BloodRequestID)
	if err != nil {
		return nil, err
	}

	filter := domain.SearchMatchesRequest{
		BloodType: bloodReq.BloodType,
		City:      bloodReq.City,
	}

	matches, err := u.userServiceRepo.SearchMatches(ctx, token, &filter)
	if err != nil {
		return nil, err
	}

	if len(matches) > 0 {
		err = u.bloodReqRepo.CreateMatches(ctx, bloodReq.ID, matches)
		if err != nil {
			return nil, err
		}
	}

	return matches, nil
}

func (u *bloodReqUsecase) ProcessDonorMatches() error {
	ctx, cancel := context.WithTimeout(context.Background(), cronTimeOut)
	defer cancel()

	pendingReqs, err := u.bloodReqRepo.GetPendingReqs(ctx)
	if err != nil {
		return err
	}
	if len(pendingReqs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	for _, req := range pendingReqs {
		wg.Add(1)

		go func(request entity.BloodRequest) {
			defer wg.Done()

			searchReq := domain.SearchMatchesRequest{
				BloodType: request.BloodType,
				City:      request.City,
			}

			matches, err := u.userServiceRepo.SearchMatches(ctx, "", &searchReq)
			if err != nil {
				log.Println(err)
				return
			}

			if len(matches) > 0 {
				err = u.bloodReqRepo.CreateMatches(ctx, request.ID, matches)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}(req)
	}
	wg.Wait()

	return nil
}
