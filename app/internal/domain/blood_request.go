package domain

import (
	"context"
	"donor-service/app/internal/entity"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrInvalidBloodType = errors.New("invalid blood type")
	ErrInvalidCoord     = errors.New("invalid latitude and longitude coords")
	ErrInvalidNeededAt  = errors.New("invalid needed at date")
	ErrBloodReqNotFound = errors.New("blood request not found")
)

type BloodRequestReq struct {
	BloodType          string    `json:"blood_type"`
	Quantity           int       `json:"quantity"`
	Urgency            string    `json:"urgency"`
	HospitalExternalID string    `json:"hospital_external_id"`
	HospitalName       string    `json:"hospital_name"`
	City               string    `json:"city"`
	Latitude           float64   `json:"latitude"`
	Longitude          float64   `json:"longitude"`
	Notes              string    `json:"notes"`
	NeededAt           time.Time `json:"needed_at"`
}

type SearchMatchesRequest struct {
	BloodType string `query:"blood_type"`
	City      string `query:"city"`
}

type BloodRequestRepository interface {
	Create(ctx context.Context, bloodReq *entity.BloodRequest) error
	FindAll(ctx context.Context, userID uuid.UUID) ([]entity.BloodRequest, error)
	FindOne(ctx context.Context, userID uuid.UUID, bloodRequestID uuid.UUID) (*entity.BloodRequest, error)
	Cancel(ctx context.Context, userID uuid.UUID, bloodRequestID uuid.UUID) error

	GetById(ctx context.Context, bloodRequestID uuid.UUID) (*entity.BloodRequest, error)
	GetMatches(ctx context.Context, bloodRequestID uuid.UUID) ([]entity.DonorMatch, error)
	CreateMatches(ctx context.Context, bloodRequestID uuid.UUID, donors []entity.DonorProfile) error
	GetPendingReqs(ctx context.Context) ([]entity.BloodRequest, error)
}

type BloodRequestUsecase interface {
	Create(userID uuid.UUID, req *BloodRequestReq) (*entity.BloodRequest, error)
	FindAll(userID uuid.UUID) ([]entity.BloodRequest, error)
	FindOne(userID uuid.UUID, bloodRequestID uuid.UUID) (*entity.BloodRequest, error)
	Cancel(userID uuid.UUID, bloodRequestID uuid.UUID) error

	GetMatches(bloodRequestID uuid.UUID) ([]entity.DonorMatch, error)
	SearchMatches(token string, userID uuid.UUID, bloodRequestID uuid.UUID) ([]entity.DonorProfile, error)
	ProcessDonorMatches() error
}

type BloodRequestHttpRepo interface {
	SearchMatches(ctx context.Context, token string, req *SearchMatchesRequest) ([]entity.DonorProfile, error)
}
