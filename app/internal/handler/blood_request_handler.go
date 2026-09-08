package handler

import (
	"donor-service/app/internal/domain"

	"github.com/labstack/echo/v5"
)

type BloodRequestHandler struct {
	bloodReqUse domain.BloodRequestUsecase
}

func NewBloodRequestHandler(bloodReqUse domain.BloodRequestUsecase) *BloodRequestHandler {
	return &BloodRequestHandler{
		bloodReqUse: bloodReqUse,
	}
}

func (h *BloodRequestHandler) Create(c *echo.Context) error {
	return nil
}
func (h *BloodRequestHandler) FindAll(c *echo.Context) error {
	return nil
}
func (h *BloodRequestHandler) FindOne(c *echo.Context) error {
	return nil
}
func (h *BloodRequestHandler) Cancel(c *echo.Context) error {
	return nil
}
func (h *BloodRequestHandler) FindMatches(c *echo.Context) error {
	return nil
}
