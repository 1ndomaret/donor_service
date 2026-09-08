package handler

import (
	"donor-service/app/internal/domain"
	"donor-service/app/internal/helper"
	"errors"

	"github.com/google/uuid"

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
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(c, "invalid or missing token")
	}

	var req domain.BloodRequestReq
	if err := c.Bind(&req); err != nil {
		return helper.BadRequest(c, err.Error())
	}

	bloodReq, err := h.bloodReqUse.Create(userID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			return helper.BadRequest(c, domain.BloodRequestReq{})
		}

		if errors.Is(err, domain.ErrInvalidCoord) || errors.Is(err, domain.ErrInvalidNeededAt) {
			return helper.Unprocessable(c, err.Error())
		}

		if errors.Is(err, domain.ErrBloodReqNotFound) {
			return helper.NotFound(c, err.Error())
		}

		return helper.InternalServerError(c, err.Error())
	}

	return helper.Created(c, bloodReq)
}

func (h *BloodRequestHandler) FindAll(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(c, "invalid or missing token")
	}

	bloodReqs, err := h.bloodReqUse.FindAll(userID)
	if err != nil {
		return helper.InternalServerError(c, err.Error())
	}
	return helper.Success(c, 200, "Success", bloodReqs)
}

func (h *BloodRequestHandler) FindOne(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(c, "invalid or missing token")
	}

	bloodReqID := c.Param("id")

	reqUUID, err := uuid.Parse(bloodReqID)
	if err != nil {
		return helper.Unprocessable(c, "invalid blood request id")
	}

	bloodReq, err := h.bloodReqUse.FindOne(userID, reqUUID)
	if err != nil {
		if errors.Is(err, domain.ErrBloodReqNotFound) {
			return helper.NotFound(c, err.Error())
		}

		return helper.InternalServerError(c, err.Error())
	}
	return helper.Success(c, 200, "Success", bloodReq)
}

func (h *BloodRequestHandler) Cancel(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(c, "invalid or missing token")
	}

	bloodReqID := c.Param("id")

	reqUUID, err := uuid.Parse(bloodReqID)
	if err != nil {
		return helper.Unprocessable(c, "invalid blood request id")
	}

	if err := h.bloodReqUse.Cancel(userID, reqUUID); err != nil {
		if errors.Is(err, domain.ErrBloodReqNotFound) {
			return helper.NotFound(c, err.Error())
		}

		return helper.InternalServerError(c, err.Error())
	}
	return helper.Success(c, 200, "Success", "Request Cancelled")
}

func (h *BloodRequestHandler) FindMatches(c *echo.Context) error {
	return nil
}
