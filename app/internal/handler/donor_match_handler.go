package handler

import (
	"donor-service/app/internal/domain"
	"donor-service/app/internal/helper"
	"errors"

	"github.com/google/uuid"

	"github.com/labstack/echo/v5"
)

type DonorMatchHandler struct {
	donorMatchUsecase domain.DonorMatchUsecase
}

func NewDonorMatchHandler(
	donorMatchUsecase domain.DonorMatchUsecase,
) *DonorMatchHandler {
	return &DonorMatchHandler{
		donorMatchUsecase: donorMatchUsecase,
	}
}

func (h *DonorMatchHandler) Invite(c *echo.Context) error {
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(c, "invalid or missing token")
	}

	bloodRequestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "invalid blood request id")
	}

	donorID, err := uuid.Parse(c.Param("donorId"))
	if err != nil {
		return helper.BadRequest(c, "invalid donor id")
	}

	donorMatch, err := h.donorMatchUsecase.Invite(
		requesterID,
		bloodRequestID,
		donorID,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDonorMatchNotFound):
			return helper.NotFound(c, err.Error())

		case errors.Is(err, domain.ErrForbidden):
			return helper.Forbidden(c, err.Error())

		case errors.Is(err, domain.ErrInvalidMatchStatus):
			return helper.BadRequest(c, err.Error())

		default:
			return helper.InternalServerError(c, err.Error())
		}
	}

	return helper.Success(
		c,
		200,
		"Donor invited successfully",
		donorMatch,
	)
}

func (h *DonorMatchHandler) Accept(c *echo.Context) error {
	donorID, ok := c.Get("donor_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(c, "invalid or missing token")
	}

	matchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "invalid donor match id")
	}

	donorMatch, err := h.donorMatchUsecase.Accept(
		donorID,
		matchID,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDonorMatchNotFound):
			return helper.NotFound(c, err.Error())

		case errors.Is(err, domain.ErrForbidden):
			return helper.Forbidden(c, err.Error())

		case errors.Is(err, domain.ErrInvalidMatchStatus):
			return helper.BadRequest(c, err.Error())

		default:
			return helper.InternalServerError(c, err.Error())
		}
	}

	return helper.Success(
		c,
		200,
		"Invitation accepted",
		donorMatch,
	)
}

func (h *DonorMatchHandler) Decline(c *echo.Context) error {
	donorID, ok := c.Get("donor_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(c, "invalid or missing token")
	}

	matchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "invalid donor match id")
	}

	donorMatch, err := h.donorMatchUsecase.Decline(
		donorID,
		matchID,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDonorMatchNotFound):
			return helper.NotFound(c, err.Error())

		case errors.Is(err, domain.ErrForbidden):
			return helper.Forbidden(c, err.Error())

		case errors.Is(err, domain.ErrInvalidMatchStatus):
			return helper.BadRequest(c, err.Error())

		case errors.Is(err, domain.ErrBloodRequestFulfilled):
			return helper.BadRequest(c, err.Error())

		default:
			return helper.InternalServerError(c, err.Error())
		}
	}

	return helper.Success(
		c,
		200,
		"Invitation declined",
		donorMatch,
	)
}

func (h *DonorMatchHandler) GetByRequesterID(c *echo.Context) error {
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(c, "invalid or missing token")
	}

	donorMatches, err := h.donorMatchUsecase.GetByRequesterID(requesterID)
	if err != nil {
		return helper.InternalServerError(c, err.Error())
	}

	return helper.Success(
		c,
		200,
		"Success",
		donorMatches,
	)
}
