package handler

import (
	"donor-service/app/internal/domain"
	"donor-service/app/internal/helper"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type GeoapifyHandler struct {
	geoapifyUsecase domain.GeoapifyUsecase
}

func NewGeoapifyHandler(
	geoapifyUsecase domain.GeoapifyUsecase,
) *GeoapifyHandler {
	return &GeoapifyHandler{
		geoapifyUsecase: geoapifyUsecase,
	}
}

// GetGeoapifyHospitals godoc
// @Summary Get hospitals by city
// @Description Get list of hospitals from Geoapify based on city
// @Tags Geoapify
// @Produce json
// @Security BearerAuth
// @Param city query string true "City name" example(Bekasi)
// @Success 200 {object} helper.GeoapifyHospitalListSwaggoResponse "Success"
// @Failure 400 {object} helper.ErrorSwaggoResponse "Bad Request"
// @Failure 401 {object} helper.ErrorSwaggoResponse "Unauthorized"
// @Failure 500 {object} helper.ErrorSwaggoResponse "Internal Server Error"
// @Router /hospitals [get]
func (h *GeoapifyHandler) GetGeoapifyHospitals(
	c *echo.Context,
) error {
	city := c.QueryParam("city")

	if city == "" {
		return helper.BadRequest(
			c,
			"city is required",
		)
	}

	geoapifys, err := h.geoapifyUsecase.GetGeoapifyHospitals(city)
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
		geoapifys,
	)
}

// GetGeoapifyRoute godoc
// @Summary Get route for donor match
// @Description Get route information for a donor match using Geoapify
// @Tags Geoapify
// @Produce json
// @Security BearerAuth
// @Param id path string true "Donor Match ID"
// @Success 200 {object} helper.GeoapifyRouteSwaggoResponse "Success"
// @Failure 401 {object} helper.ErrorSwaggoResponse "Unauthorized"
// @Failure 422 {object} helper.ErrorSwaggoResponse "Unprocessable Content"
// @Failure 500 {object} helper.ErrorSwaggoResponse "Internal Server Error"
// @Router /donor-matches/{id}/route [get]
func (h *GeoapifyHandler) GetGeoapifyRoute(c *echo.Context) error {
	donorMatchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.Unprocessable(c, "invalid donor match id")
	}

	route, err := h.geoapifyUsecase.GetGeoapifyRoute(donorMatchID)
	if err != nil {
		return helper.InternalServerError(c, err.Error())
	}

	return helper.Success(c, 200, "Success", route)
}
