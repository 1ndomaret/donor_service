package domain

import (
	"context"
	"donor-service/app/internal/entity"

	"github.com/google/uuid"
)

type UserServiceHttpRepo interface {
	SearchMatches(ctx context.Context, token string, req *SearchMatchesRequest) ([]entity.DonorProfile, error)
	GetDonorProfile(ctx context.Context, donor_id uuid.UUID) (*entity.DonorProfile, error)
}
