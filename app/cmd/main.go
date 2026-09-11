package main

import (
	"donor-service/app/internal/config"
	"donor-service/app/internal/entity"
	"donor-service/app/internal/handler"
	"donor-service/app/internal/repository"
	httpRepo "donor-service/app/internal/repository/http"
	"donor-service/app/internal/router"
	"donor-service/app/internal/scheduler"
	"donor-service/app/internal/usecase"

	_ "donor-service/docs"

	"log"

	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

// @title BloodConnect Donor Service API
// @version 1.0
// @description REST API for BloodConnect donor, blood request, donation, and hospital services.
// @host localhost:1324
// @BasePath /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(
		&entity.BloodRequest{},
		&entity.Donation{},
		&entity.DonorMatch{},
	); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	servicesConfig := config.ServicesSecret()

	userServiceHttpRepo := httpRepo.NewUserServiceHttpRepo(servicesConfig)

	bloodReqRepository := repository.NewBloodRequestRepository(db)
	bloodReqUsecase := usecase.NewBloodRequestUsecase(bloodReqRepository, userServiceHttpRepo)
	bloodReqHandler := handler.NewBloodRequestHandler(bloodReqUsecase)

	geoapifyRepository := httpRepo.NewGeoapifyHttpRepo(servicesConfig)

	donorMatchRepository := repository.NewDonorMatchRepository(db)
	donationRepository := repository.NewDonationRepository(db)

	donorMatchUsecase := usecase.NewDonorMatchUsecase(donorMatchRepository, bloodReqRepository, userServiceHttpRepo, geoapifyRepository, donationRepository)
	donorMatchHandler := handler.NewDonorMatchHandler(donorMatchUsecase)

	donationUsecase := usecase.NewDonationUsecase(donationRepository, donorMatchRepository, bloodReqRepository)
	donationHandler := handler.NewDonationHandler(donationUsecase)

	geoapifyUsecase := usecase.NewGeoapifyUsecase(geoapifyRepository, donorMatchRepository, userServiceHttpRepo, bloodReqRepository)
	geoapifyHandler := handler.NewGeoapifyHandler(geoapifyUsecase)

	scheduler := scheduler.NewScheduler(bloodReqUsecase)
	if err := scheduler.Start(); err != nil {
		log.Println(err.Error())
	}

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	router.Register(e,
		bloodReqHandler,
		donorMatchHandler,
		donationHandler,
		geoapifyHandler,
	)

	if err := e.Start(":1324"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
