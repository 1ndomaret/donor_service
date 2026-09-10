package scheduler

import (
	"donor-service/app/internal/domain"
	"log"

	"github.com/jasonlvhit/gocron"
)

type scheduler struct {
	bloodReqUse domain.BloodRequestUsecase
}

func NewScheduler(bloodReqUse domain.BloodRequestUsecase) *scheduler {
	return &scheduler{
		bloodReqUse: bloodReqUse,
	}
}

func (s *scheduler) Start() error {
	gocron.Every(1).Hour().Do(s.FindDonorMatch)
	gocron.Start()

	return nil
}

func (s *scheduler) FindDonorMatch() {
	log.Println("[CRON] Finding Donor Match...")
	if err := s.bloodReqUse.ProcessDonorMatches(); err != nil {
		log.Printf("[CRON] ERROR: Failed to process donor matches: %v\n", err)
	}
}
