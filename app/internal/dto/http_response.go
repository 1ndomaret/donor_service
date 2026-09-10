package dto

import "donor-service/app/internal/entity"

type MatchSearchResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Data    []entity.DonorProfile `json:"data"`
}

type GetDonorResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Data    *entity.DonorProfile `json:"data"`
}
