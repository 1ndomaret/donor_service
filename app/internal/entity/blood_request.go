package entity

import (
	"time"

	"github.com/google/uuid"
)

type BloodRequest struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RequesterID        uuid.UUID `gorm:"type:uuid;index;not null" json:"requester_id"`
	BloodType          string    `gorm:"not null" json:"blood_type"`
	Quantity           int       `gorm:"not null" json:"quantity"`
	Urgency            string    `gorm:"not null" json:"urgency"`
	HospitalExternalID string    `gorm:"not null" json:"hospital_external_id"`
	HospitalName       string    `gorm:"not null" json:"hospital_name"`
	City               string    `gorm:"not null" json:"city"`
	Latitude           float64   `json:"latitude"`
	Longitude          float64   `json:"longitude"`
	Notes              string    `gorm:"not null" json:"notes"`
	Status             string    `gorm:"not null" json:"status"`
	NeededAt           time.Time `json:"needed_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	DonorMatches []DonorMatch `gorm:"foreignKey:BloodRequestID;references:ID" json:"donor_matches,omitempty"`
}
