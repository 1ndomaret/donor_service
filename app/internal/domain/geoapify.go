package domain

import (
	"context"
	"donor-service/app/internal/dto"

	"github.com/google/uuid"
)

type GeoapifyHospital struct {
	ExternalID string  `json:"external_id"`
	Name       string  `json:"name"`
	City       string  `json:"city"`
	Address    string  `json:"address"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type GeoapifyRoute struct {
	Distance      float64 `json:"distance"`
	DistanceUnits string  `json:"distance_units"`
	Time          float64 `json:"time"`
	Mode          string  `json:"mode"`
}

type GeoapifyRepository interface {
	GetGeoapifyHospitals(
		ctx context.Context,
		city string,
	) ([]GeoapifyHospital, error)

	GetGeoapifyRoute(ctx context.Context, req *dto.GeoapifyRoutingRequest) (*GeoapifyRoute, error)
}

type GeoapifyUsecase interface {
	GetGeoapifyHospitals(
		city string,
	) ([]GeoapifyHospital, error)

	GetGeoapifyRoute(donorMatchID uuid.UUID) (*GeoapifyRoute, error)
}
