package handler

import (
	"donor-service/app/internal/domain"
	"donor-service/app/internal/helper"
	"errors"

	"github.com/google/uuid"

	"github.com/labstack/echo/v5"
)

type DonorMatchesHandler struct {
	donorMatchesUsecase domain.DonorMatchesUsecase
}

func NewDonorMatchesHandler(
	donorMatchesUsecase domain.DonorMatchesUsecase,
) *DonorMatchesHandler {
	return &DonorMatchesHandler{
		donorMatchesUsecase: donorMatchesUsecase,
	}
}

func (h *DonorMatchesHandler) Invite(c *echo.Context) error {
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

	donorMatches, err := h.donorMatchesUsecase.Invite(
		requesterID,
		bloodRequestID,
		donorID,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDonorMatchesNotFound):
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
		donorMatches,
	)
}

func (h *DonorMatchesHandler) Accept(c *echo.Context) error {
	donorID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(c, "invalid or missing token")
	}

	matchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "invalid donor match id")
	}

	donorMatches, err := h.donorMatchesUsecase.Accept(
		donorID,
		matchID,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDonorMatchesNotFound):
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
		donorMatches,
	)
}

func (h *DonorMatchesHandler) Decline(c *echo.Context) error {
	donorID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(c, "invalid or missing token")
	}

	matchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "invalid donor match id")
	}

	donorMatches, err := h.donorMatchesUsecase.Decline(
		donorID,
		matchID,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDonorMatchesNotFound):
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
		"Invitation declined",
		donorMatches,
	)
}
