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

	"github.com/go-resty/resty/v2"
)

type bloodRequestHttpRepo struct {
	client *resty.Client
}

func NewBloodRequestHttpRepo(cfg config.ServicesConfig) domain.BloodRequestHttpRepo {
	client := resty.New()

	baseURL := fmt.Sprintf("%s:%s", cfg.UserURL, cfg.UserPort)
	fmt.Println(baseURL)

	client.SetBaseURL(baseURL)
	client.SetTimeout(10 * time.Second)

	return &bloodRequestHttpRepo{
		client: client,
	}
}

func (r *bloodRequestHttpRepo) SearchMatches(ctx context.Context, token string, req *domain.SearchMatchesRequest) ([]entity.DonorProfile, error) {
	var apiRes dto.MatchSearchResponse

	queryParams := map[string]string{}
	if req.BloodType != "" {
		queryParams["blood_type"] = req.BloodType
	}
	if req.City != "" {
		queryParams["city"] = req.City
	}

	res, err := r.client.R().SetHeader("Authorization", token).SetContext(ctx).SetQueryParams(queryParams).SetResult(&apiRes).Get("/api/v1/users/donor-profile/search")
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
