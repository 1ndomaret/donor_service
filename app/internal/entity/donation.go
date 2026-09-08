package entity

import (
	"time"

	"github.com/google/uuid"
)

type Donation struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	BloodRequestID uuid.UUID  `gorm:"type:uuid;index;not null" json:"blood_request_id"`
	DonorID        uuid.UUID  `gorm:"type:uuid;index;not null" json:"donor_id"`
	DonorMatchID   uuid.UUID  `gorm:"type:uuid;index;not null" json:"donor_match_id"`
	DonationDate   *time.Time `json:"donation_date"`
	Status         string     `gorm:"not null" json:"status"`
	ConfirmedBy    *uuid.UUID `gorm:"type:uuid" json:"confirmed_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
