package handler

import (
	"donor-service/app/internal/domain"
	"donor-service/app/internal/helper"
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type DonationHandler struct {
	donationUsecase domain.DonationUsecase
}

func NewDonationHandler(
	donationUsecase domain.DonationUsecase,
) *DonationHandler {
	return &DonationHandler{
		donationUsecase: donationUsecase,
	}
}

// Create godoc
// @Summary Create donation
// @Description Create a pending donation from an accepted donor match
// @Tags Donations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.DonationReq true "Donation data"
// @Success 201 {object} helper.DonationSwaggoResponse "Donation created"
// @Failure 400 {object} helper.ErrorSwaggoResponse "Bad Request"
// @Failure 401 {object} helper.ErrorSwaggoResponse "Unauthorized"
// @Failure 404 {object} helper.ErrorSwaggoResponse "Not Found"
// @Failure 500 {object} helper.ErrorSwaggoResponse "Internal Server Error"
// @Router /donations [post]
func (h *DonationHandler) Create(
	c *echo.Context,
) error {

	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(
			c,
			"invalid or missing token",
		)
	}

	// role, ok := c.Get("role").(string)
	// if !ok || role != "requester" {
	// 	return helper.Forbidden(
	// 		c,
	// 		"only requester can create donation",
	// 	)
	// }

	var req domain.DonationReq

	if err := c.Bind(&req); err != nil {
		return helper.BadRequest(c, err.Error())
	}

	donation, err := h.donationUsecase.Create(
		requesterID,
		&req,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidDonationInput):
			return helper.BadRequest(c, err.Error())

		case errors.Is(err, domain.ErrBloodReqNotFound):
			return helper.NotFound(c, err.Error())

		case errors.Is(err, domain.ErrDonorMatchNotFound):
			return helper.NotFound(c, err.Error())

		case errors.Is(err, domain.ErrDonorMatchNotAccepted):
			return helper.BadRequest(c, err.Error())

		default:
			return helper.InternalServerError(c, err.Error())
		}
	}

	return helper.Created(c, donation)
}

// Completed godoc
// @Summary Complete donation
// @Description Complete a pending donation and confirm it by the requester
// @Tags Donations
// @Produce json
// @Security BearerAuth
// @Param id path string true "Donation ID"
// @Success 200 {object} helper.DonationSwaggoResponse "Donation completed"
// @Failure 400 {object} helper.ErrorSwaggoResponse "Bad Request"
// @Failure 401 {object} helper.ErrorSwaggoResponse "Unauthorized"
// @Failure 404 {object} helper.ErrorSwaggoResponse "Not Found"
// @Failure 422 {object} helper.ErrorSwaggoResponse "Unprocessable Content"
// @Failure 500 {object} helper.ErrorSwaggoResponse "Internal Server Error"
// @Router /donations/{id}/complete [patch]
func (h *DonationHandler) Completed(
	c *echo.Context,
) error {

	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(
			c,
			"invalid or missing token",
		)
	}

	// role, ok := c.Get("role").(string)
	// if !ok || role != "requester" {
	// 	return helper.Forbidden(
	// 		c,
	// 		"only requester can complete donation",
	// 	)
	// }

	donationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.Unprocessable(
			c,
			"invalid donation id",
		)
	}

	donation, err := h.donationUsecase.Completed(
		requesterID,
		donationID,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDonationNotFound):
			return helper.NotFound(c, err.Error())

		case errors.Is(err, domain.ErrBloodReqNotFound):
			return helper.NotFound(c, err.Error())

		case errors.Is(err, domain.ErrInvalidDonationStatus):
			return helper.BadRequest(c, err.Error())

		default:
			return helper.InternalServerError(c, err.Error())
		}
	}

	return helper.Success(
		c,
		200,
		"Donation completed",
		donation,
	)
}
