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

// Invite godoc
// @Summary Invite donor
// @Description Invite a donor to a blood request
// @Tags Donor Matches
// @Produce json
// @Security BearerAuth
// @Param id path string true "Blood Request ID"
// @Param donorId path string true "Donor ID"
// @Success 200 {object} helper.DonorMatchSwaggoResponse "Donor invited successfully"
// @Failure 400 {object} helper.ErrorSwaggoResponse "Bad Request"
// @Failure 401 {object} helper.ErrorSwaggoResponse "Unauthorized"
// @Failure 403 {object} helper.ErrorSwaggoResponse "Forbidden"
// @Failure 404 {object} helper.ErrorSwaggoResponse "Not Found"
// @Failure 500 {object} helper.ErrorSwaggoResponse "Internal Server Error"
// @Router /blood-requests/{id}/invite/{donorId} [post]
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

// Accept godoc
// @Summary Accept donor invitation
// @Description Accept donor invitation, validate blood request quantity, update distance, and create a pending donation
// @Tags Donor Matches
// @Produce json
// @Security BearerAuth
// @Param id path string true "Donor Match ID"
// @Success 200 {object} helper.DonorMatchSwaggoResponse "Invitation accepted"
// @Failure 400 {object} helper.ErrorSwaggoResponse "Bad Request"
// @Failure 401 {object} helper.ErrorSwaggoResponse "Unauthorized"
// @Failure 403 {object} helper.ErrorSwaggoResponse "Forbidden"
// @Failure 404 {object} helper.ErrorSwaggoResponse "Not Found"
// @Failure 500 {object} helper.ErrorSwaggoResponse "Internal Server Error"
// @Router /donor-matches/{id}/accept [patch]
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

		case errors.Is(err, domain.ErrBloodRequestFulfilled):
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

// Decline godoc
// @Summary Decline donor invitation
// @Description Decline a donor match invitation
// @Tags Donor Matches
// @Produce json
// @Security BearerAuth
// @Param id path string true "Donor Match ID"
// @Success 200 {object} helper.DonorMatchSwaggoResponse "Invitation declined"
// @Failure 400 {object} helper.ErrorSwaggoResponse "Bad Request"
// @Failure 401 {object} helper.ErrorSwaggoResponse "Unauthorized"
// @Failure 403 {object} helper.ErrorSwaggoResponse "Forbidden"
// @Failure 404 {object} helper.ErrorSwaggoResponse "Not Found"
// @Failure 500 {object} helper.ErrorSwaggoResponse "Internal Server Error"
// @Router /donor-matches/{id}/decline [patch]
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

// GetByDonorID godoc
// @Summary Get donor matches for requester
// @Description Get donor matches from blood requests owned by the authenticated requester
// @Tags Donor Matches
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helper.DonorMatchListSwaggoResponse "Success"
// @Failure 401 {object} helper.ErrorSwaggoResponse "Unauthorized"
// @Failure 500 {object} helper.ErrorSwaggoResponse "Internal Server Error"
// @Router /donor-matches [get]
func (h *DonorMatchHandler) GetByDonorID(c *echo.Context) error {
	donorID, ok := c.Get("donor_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(c, "invalid or missing token")
	}

	donorMatches, err := h.donorMatchUsecase.GetByDonorID(donorID)
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
