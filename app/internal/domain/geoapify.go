package domain

import "context"

type Geoapify struct {
	ExternalID string  `json:"external_id"`
	Name       string  `json:"name"`
	City       string  `json:"city"`
	Address    string  `json:"address"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type GeoapifyRepository interface {
	GetGeoapifys(
		ctx context.Context,
		city string,
	) ([]Geoapify, error)
}

type GeoapifyUsecase interface {
	GetGeoapifys(
		city string,
	) ([]Geoapify, error)
}
