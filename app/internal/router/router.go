package router

import (
	"donor-service/app/internal/handler"
	"donor-service/app/internal/middleware"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Echo,
	bloodReqHandler *handler.BloodRequestHandler,
	donorMatchHandler *handler.DonorMatchHandler,
	donationHandler *handler.DonationHandler,
	geoapifyHandler *handler.GeoapifyHandler,
) {
	api := e.Group("/api/v1")

	private := api.Group("")
	private.Use(echojwt.WithConfig(middleware.JwtConfig()), middleware.ParseJwtClaims)

	private.POST("/blood-requests", bloodReqHandler.Create)
	private.GET("/blood-requests", bloodReqHandler.FindAll)
	private.GET("/blood-requests/:id", bloodReqHandler.FindOne)
	private.PUT("/blood-requests/:id/cancel", bloodReqHandler.Cancel)

	private.GET("/blood-requests/:id/matches", bloodReqHandler.GetMatches)

	private.POST("/blood-requests/:id/invite/:donorId", donorMatchHandler.Invite)
	private.PATCH("/donor-matches/:id/accept", donorMatchHandler.Accept)
	private.PATCH("/donor-matches/:id/decline", donorMatchHandler.Decline)
	private.POST("/donations", donationHandler.Create)
	private.PATCH("/donations/:id/complete", donationHandler.Completed)

	private.GET("/hospitals", geoapifyHandler.GetGeoapifyHospitals)
	private.GET("/donor-matches/:id/route", geoapifyHandler.GetGeoapifyRoute)

	/*
		- Bikin blood-request
			- dia ngecek juga donor profile yang sama blood-typenya, isAvailable, city
				- kalau nemu dia akan bikin donor_matches

		- Kalau nemu match, dia bisa nge invite user-user yang cocok

		- Buat pendonor dia bisa accept/decline
			- Auto decline setelah seminggu
	*/

	// TODOS:
	// Go Routine match checker buat blood-request jalan setiap jam
	// - Get donors by blood type, availability, city(?)

	// A
	// POST /blood-requests <-- autorun match checker pertama after create
	// GET  /blood-requests
	// GET  /blood-requests/:id
	// PUT  /blood-requests/:id/cancel

	// GET /blood-requests/:id/matches

	// B
	// POST /blood-requests/:id/invite/:donorId

	// PATCH /donor-matches/:id/accept
	// PATCH /donor-matches/:id/decline

	// POST  /donations <-- yang bisa create si requester
	// PATCH /donations/:id/complete <-- yang bisa complete si requester

	// GET /geoapifys

	/*
		http://localhost:1323/api/v1/users/donor-profile/search?blood_type=A+&city=Jakarta
		- masih case sensitive

		http://localhost:1324/api/v1/blood-requests
		- masih nge-get cancelled requests

		http://localhost:1324/api/v1/blood-requests
		- auto create donor matches

		http://localhost:1324/api/v1/donor-matches/:id/accept
		- cuma bisa accept kalau BloodRequest.DonorMatches.length < BloodRequest.Quantity
		- auto update DonorMatch.DistanceKm

		http://localhost:1324/api/v1/donor-matches/:id/accept
		- auto create Donation setelah accept

		http://localhost:1324/api/v1/donor_match/
		- get donor match for requester
	*/
}
