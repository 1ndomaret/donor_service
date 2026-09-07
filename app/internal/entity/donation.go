package entity

import (
	"time"

	"github.com/google/uuid"
)

type Donation struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	BloodRequestID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"blood_request_id"`
	DonorID        uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"donor_id"`
	DonationDate   time.Time `json:"donation_date"`
	Status         string    `gorm:"not null" json:"status"`
	Confirmedby    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"confirmed_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
