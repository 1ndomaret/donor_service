package helper

import (
	"donor-service/app/internal/domain"
	"donor-service/app/internal/entity"

	"github.com/labstack/echo/v5"
)

type ErrorResponse struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
	Detail    any    `json:"detail,omitempty"`
}

type Response struct {
	Status  string         `json:"status"`
	Message string         `json:"message,omitempty"`
	Data    any            `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

type DonorMatchSwaggoResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message,omitempty"`
	Data    entity.DonorMatch `json:"data"`
}

type DonorMatchListSwaggoResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message,omitempty"`
	Data    []entity.DonorMatch `json:"data"`
}

type BloodRequestSwaggoResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message,omitempty"`
	Data    entity.BloodRequest `json:"data"`
}

type BloodRequestListSwaggoResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message,omitempty"`
	Data    []entity.BloodRequest `json:"data"`
}

type DonationSwaggoResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message,omitempty"`
	Data    entity.Donation `json:"data"`
}

type GeoapifyHospitalListSwaggoResponse struct {
	Status  string                    `json:"status"`
	Message string                    `json:"message,omitempty"`
	Data    []domain.GeoapifyHospital `json:"data"`
}

type GeoapifyRouteSwaggoResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message,omitempty"`
	Data    domain.GeoapifyRoute `json:"data"`
}

type ErrorDetailSwaggoResponse struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
	Detail    string `json:"detail,omitempty"`
}

type ErrorSwaggoResponse struct {
	Status string                    `json:"status"`
	Error  ErrorDetailSwaggoResponse `json:"error"`
}

func Success(c *echo.Context, httpCode int, msg string, data any) error {
	return c.JSON(httpCode, Response{
		Status:  "Success",
		Message: msg,
		Data:    data,
	})
}

func Error(c *echo.Context, httpCode int, errCode, msg string, data any) error {
	return c.JSON(httpCode, Response{
		Status: "Error",
		Error: &ErrorResponse{
			ErrorCode: errCode,
			Message:   msg,
			Detail:    data,
		},
	})
}

func Created(c *echo.Context, data any) error {
	return Success(c, 201, "Created", data)
}

func BadRequest(c *echo.Context, data any) error {
	return Error(c, 400, "BAD_REQUEST", "Bad Request", data)
}

func Unauthorized(c *echo.Context, data any) error {
	return Error(c, 401, "UNAUTHORIZED", "Unauthorized", data)
}

func Forbidden(c *echo.Context, data any) error {
	return Error(c, 403, "FORBIDDEN", "Forbidden", data)
}

func NotFound(c *echo.Context, data any) error {
	return Error(c, 404, "NOT_FOUND", "Not Found", data)
}

func Unprocessable(c *echo.Context, data any) error {
	return Error(c, 422, "UNPROCESSABLE_CONTENT", "Unprocessable Content", data)
}

func InternalServerError(c *echo.Context, data any) error {
	return Error(c, 500, "INTERNAL_SERVER_ERROR", "Internal Server Error", data)
}
