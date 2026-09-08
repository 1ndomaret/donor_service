package scheduler

import (
	"log"

	"github.com/jasonlvhit/gocron"
)

type scheduler struct {
}

func NewScheduler() *scheduler {
	return &scheduler{}
}

func (s *scheduler) Start() error {
	log.Println("[CRON] Starting Cron Job...")
	gocron.Every(1).Hour().Do(s.FindDonorMatches)
	gocron.Start()

	return nil
}

func (s *scheduler) FindDonorMatches() {
	log.Println("[CRON] Finding Donor Matches...")
}
