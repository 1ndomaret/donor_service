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

type BloodRequestRepository interface {
	Create(ctx context.Context, req *entity.BloodRequest) error
	FindAll(ctx context.Context, userID uuid.UUID) ([]entity.BloodRequest, error)
	FindOne(ctx context.Context, userID uuid.UUID, bloodRequestID uuid.UUID) (*entity.BloodRequest, error)
	Cancel(ctx context.Context, userID uuid.UUID, bloodRequestID uuid.UUID) error

	FindMatches(ctx context.Context, bloodRequestID uuid.UUID) ([]entity.DonorMatches, error)
}

type BloodRequestUsecase interface {
	Create(userID uuid.UUID, req *BloodRequestReq) (*entity.BloodRequest, error)
	FindAll(userID uuid.UUID) ([]entity.BloodRequest, error)
	FindOne(userID uuid.UUID, bloodRequestID uuid.UUID) (*entity.BloodRequest, error)
	Cancel(userID uuid.UUID, bloodRequestID uuid.UUID) error

	FindMatches(bloodRequestID uuid.UUID) ([]entity.DonorMatches, error)
}
