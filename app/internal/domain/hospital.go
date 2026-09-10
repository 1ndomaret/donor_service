package domain

import "context"

type Hospital struct {
	ExternalID string  `json:"external_id"`
	Name       string  `json:"name"`
	City       string  `json:"city"`
	Address    string  `json:"address"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type HospitalRepository interface {
	GetHospitals(
		ctx context.Context,
		city string,
	) ([]Hospital, error)
}

type HospitalUsecase interface {
	GetHospitals(
		city string,
	) ([]Hospital, error)
}
