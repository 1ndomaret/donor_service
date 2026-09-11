package http

import (
	"context"
	"donor-service/app/internal/config"
	"donor-service/app/internal/domain"
	"donor-service/app/internal/dto"
	"donor-service/app/internal/entity"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/go-resty/resty/v2"
)

type userServiceHttpRepo struct {
	client       *resty.Client
	serviceToken string
}

func NewUserServiceHttpRepo(cfg config.ServicesConfig) domain.UserServiceHttpRepo {
	client := resty.New()

	baseURL := fmt.Sprintf("%s:%s", cfg.UserURL, cfg.UserPort)

	client.SetBaseURL(baseURL)
	client.SetTimeout(10 * time.Second)

	return &userServiceHttpRepo{
		client:       client,
		serviceToken: cfg.ServiceToken,
	}
}

func (r *userServiceHttpRepo) SearchMatches(ctx context.Context, token string, req *domain.SearchMatchesRequest) ([]entity.DonorProfile, error) {
	var apiRes dto.MatchSearchResponse

	queryParams := map[string]string{}
	if req.BloodType != "" {
		queryParams["blood_type"] = req.BloodType
	}
	if req.City != "" {
		queryParams["city"] = req.City
	}

	authToken := token
	if authToken == "" {
		authToken = r.serviceToken
	}

	res, err := r.client.R().
		SetHeader("Authorization", authToken).SetContext(ctx).SetQueryParams(queryParams).
		SetResult(&apiRes).
		Get("/api/v1/users/donor-profile/search")
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, errors.New(res.Status())
	}

	if len(apiRes.Data) == 0 {
		return []entity.DonorProfile{}, nil
	}

	return apiRes.Data, nil
}

func (r *userServiceHttpRepo) GetDonorProfile(ctx context.Context, donor_id uuid.UUID) (*entity.DonorProfile, error) {
	var apiRes dto.GetDonorResponse
	res, err := r.client.R().
		SetHeader("Authorization", r.serviceToken).SetContext(ctx).
		SetResult(&apiRes).
		Get("/api/v1/users/donor-profile/" + donor_id.String())
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, errors.New(res.Status())
	}

	return apiRes.Data, nil
}
