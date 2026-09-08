package entity

import (
	"time"

	"github.com/google/uuid"
)

type DonorProfile struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"user_id"`
	BloodType        string     `json:"blood_type"`
	City             string     `json:"city"`
	Latitude         float64    `json:"latitude"`
	Longitude        float64    `json:"longitude"`
	IsAvailable      bool       `json:"is_available"`
	LastDonationDate *time.Time `json:"last_donation_date,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
