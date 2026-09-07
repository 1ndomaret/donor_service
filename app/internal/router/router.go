package router

import (
	"donor-service/app/internal/middleware"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Echo,
) {
	api := e.Group("/api/v1")

	private := api.Group("")
	private.Use(echojwt.WithConfig(middleware.JwtConfig()), middleware.ParseJwtClaims)

	/*
		- Bikin blood-request
			- dia ngecek juga donor profile yang sama blood-typenya, isAvailable, city
				- kalau nemu dia akan bikin donor_matches

		- Kalau nemu match, dia bisa nge invite user-user yang cocok

		- Buat pendonor dia bisa accept/decline
			- Auto decline setelah seminggu
	*/

	// TODOS:
	// GOROUTINE buat cari matches jalan setiap jam

	// POST /blood-requests
	// GET  /blood-requests
	// GET  /blood-requests/:id
	// PUT  /blood-requests/:id/cancel

	// GET /blood-requests/:id/matches

	// POST /blood-requests/:id/invite/:donorId

	// PATCH /donor-matches/:id/accept
	// PATCH /donor-matches/:id/decline

	// POST  /donations <-- yang bisa create si requester
	// PATCH /donations/:id/complete <-- yang bisa complete si requester

	// GET /hospitals
}
