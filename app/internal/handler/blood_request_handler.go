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

// @Tags Blood Requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.BloodRequestReq true "Blood request data"
// @Success 201 {object} helper.BloodRequestSwaggoResponse "Blood request created"
// @Failure 400 {object} helper.ErrorSwaggoResponse "Bad Request"
// @Failure 401 {object} helper.ErrorSwaggoResponse "Unauthorized"
// @Failure 403 {object} helper.ErrorSwaggoResponse "Forbidden"
// @Failure 500 {object} helper.ErrorSwaggoResponse "Internal Server Error"
// @Router /blood-requests [post]
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

// FindAll godoc
// @Summary Get blood requests
// @Description Get blood requests owned by the authenticated requester
// @Tags Blood Requests
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helper.BloodRequestListSwaggoResponse "Success"
// @Failure 401 {object} helper.ErrorSwaggoResponse "Unauthorized"
// @Failure 403 {object} helper.ErrorSwaggoResponse "Forbidden"
// @Failure 500 {object} helper.ErrorSwaggoResponse "Internal Server Error"
// @Router /blood-requests [get]
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

// FindOne godoc
// @Summary Get blood request detail
// @Description Get a blood request by ID
// @Tags Blood Requests
// @Produce json
// @Security BearerAuth
// @Param id path string true "Blood Request ID"
// @Success 200 {object} helper.BloodRequestSwaggoResponse "Success"
// @Failure 401 {object} helper.ErrorSwaggoResponse "Unauthorized"
// @Failure 404 {object} helper.ErrorSwaggoResponse "Not Found"
// @Failure 422 {object} helper.ErrorSwaggoResponse "Unprocessable Content"
// @Failure 500 {object} helper.ErrorSwaggoResponse "Internal Server Error"
// @Router /blood-requests/{id} [get]
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

// Cancel godoc
// @Summary Cancel blood request
// @Description Cancel a blood request owned by the authenticated requester
// @Tags Blood Requests
// @Produce json
// @Security BearerAuth
// @Param id path string true "Blood Request ID"
// @Success 200 {object} helper.BloodRequestSwaggoResponse "Blood request cancelled"
// @Failure 400 {object} helper.ErrorSwaggoResponse "Bad Request"
// @Failure 401 {object} helper.ErrorSwaggoResponse "Unauthorized"
// @Failure 404 {object} helper.ErrorSwaggoResponse "Not Found"
// @Failure 422 {object} helper.ErrorSwaggoResponse "Unprocessable Content"
// @Failure 500 {object} helper.ErrorSwaggoResponse "Internal Server Error"
// @Router /blood-requests/{id}/cancel [put]
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

// GetMatches godoc
// @Summary Get donor matches
// @Description Get donor matches for a blood request
// @Tags Blood Requests
// @Produce json
// @Security BearerAuth
// @Param id path string true "Blood Request ID"
// @Success 200 {object} helper.DonorMatchListSwaggoResponse "Success"
// @Failure 401 {object} helper.ErrorSwaggoResponse "Unauthorized"
// @Failure 404 {object} helper.ErrorSwaggoResponse "Not Found"
// @Failure 422 {object} helper.ErrorSwaggoResponse "Unprocessable Content"
// @Failure 500 {object} helper.ErrorSwaggoResponse "Internal Server Error"
// @Router /blood-requests/{id}/matches [get]
func (h *BloodRequestHandler) GetMatches(c *echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return helper.Unauthorized(c, "missing token")
	}

	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return helper.Unauthorized(c, "invalid token")
	}

	bloodReqID := c.Param("id")
	reqUUID, err := uuid.Parse(bloodReqID)
	if err != nil {
		return helper.Unprocessable(c, "invalid blood request id")
	}

	profiles, err := h.bloodReqUse.SearchMatches(token, userID, reqUUID)
	if err != nil {
		return helper.InternalServerError(c, err.Error())
	}

	return helper.Success(c, 200, "Success", profiles)
}
