package handler

import (
	"donor-service/app/internal/domain"
	"donor-service/app/internal/helper"

	"github.com/labstack/echo/v5"
)

type HospitalHandler struct {
	hospitalUsecase domain.HospitalUsecase
}

func NewHospitalHandler(
	hospitalUsecase domain.HospitalUsecase,
) *HospitalHandler {
	return &HospitalHandler{
		hospitalUsecase: hospitalUsecase,
	}
}

func (h *HospitalHandler) GetHospitals(
	c *echo.Context,
) error {
	city := c.QueryParam("city")

	if city == "" {
		return helper.BadRequest(
			c,
			"city is required",
		)
	}

	hospitals, err := h.hospitalUsecase.GetHospitals(city)
	if err != nil {
		return helper.InternalServerError(
			c,
			err.Error(),
		)
	}

	return helper.Success(
		c,
		200,
		"Success",
		hospitals,
	)
}
