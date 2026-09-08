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
	"log"

	"github.com/labstack/echo/v5"
)

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

	bloodReqRepository := repository.NewBloodRequestRepository(db)
	bloodReqHttpRepo := httpRepo.NewBloodRequestHttpRepo(config.ServicesSecret())
	bloodReqUsecase := usecase.NewBloodRequestUsecase(bloodReqRepository, bloodReqHttpRepo)
	bloodReqHandler := handler.NewBloodRequestHandler(bloodReqUsecase)

	donorMatchRepository := repository.NewDonorMatchRepository(db)
	donorMatchUsecase := usecase.NewDonorMatchUsecase(donorMatchRepository, bloodReqRepository)
	donorMatchHandler := handler.NewDonorMatchHandler(donorMatchUsecase)
	donationRepository := repository.NewDonationRepository(db)
	donationUsecase := usecase.NewDonationUsecase(donationRepository, donorMatchRepository, bloodReqRepository)
	donationHandler := handler.NewDonationHandler(donationUsecase)

	scheduler := scheduler.NewScheduler()
	if err := scheduler.Start(); err != nil {
		log.Println(err.Error())
	}

	e := echo.New()

	router.Register(e,
		bloodReqHandler,
		donorMatchHandler,
		donationHandler,
	)

	if err := e.Start(":1324"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
