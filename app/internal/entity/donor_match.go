package entity

import (
	"time"

	"github.com/google/uuid"
)

type DonorMatch struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	BloodRequestID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_blood_request_donor" json:"blood_request_id"`
	DonorID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_blood_request_donor" json:"donor_id"`
	Status         string    `gorm:"not null" json:"status"`
	DistanceKM     float64   `json:"distance_km"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	BloodRequest *BloodRequest `gorm:"foreignKey:BloodRequestID;references:ID" json:"blood_request,omitempty"`
}
