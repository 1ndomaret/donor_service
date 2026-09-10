package handler

import (
	"donor-service/app/internal/domain"
	"donor-service/app/internal/helper"

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

func (h *GeoapifyHandler) GetGeoapifys(
	c *echo.Context,
) error {
	city := c.QueryParam("city")

	if city == "" {
		return helper.BadRequest(
			c,
			"city is required",
		)
	}

	geoapifys, err := h.geoapifyUsecase.GetGeoapifys(city)
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
