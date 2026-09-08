package domain

import (
	"context"
	"donor-service/app/internal/entity"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput = errors.New("invalid input")
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
	Status             string    `json:"status"`
	NeededAt           time.Time `json:"needed_at"`
}

type BloodRequestRepository interface {
	Create(ctx context.Context, req *entity.BloodRequest) error
	FindAll(ctx context.Context) ([]entity.BloodRequest, error)
	FindOne(ctx context.Context, BloodRequestID uuid.UUID) (*entity.BloodRequest, error)
	Cancel(ctx context.Context, BloodRequestID uuid.UUID) error

	FindMatches(ctx context.Context, BloodRequestID uuid.UUID) ([]entity.DonorMatches, error)
}

type BloodRequestUsecase interface {
	Create(req *BloodRequestReq) (*entity.BloodRequest, error)
	FindAll() ([]entity.BloodRequest, error)
	FindOne(BloodRequestID uuid.UUID) (*entity.BloodRequest, error)
	Cancel(BloodRequestID uuid.UUID) error

	FindMatches(BloodRequestID uuid.UUID) ([]entity.DonorMatches, error)
}
