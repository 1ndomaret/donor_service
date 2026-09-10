package http

import (
	"context"
	"donor-service/app/internal/config"
	"donor-service/app/internal/domain"
	"donor-service/app/internal/dto"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type geoapifyHttpRepo struct {
	client *resty.Client
	apiKey string
}

func NewGeoapifyHttpRepo(
	cfg config.ServicesConfig,
) domain.GeoapifyRepository {

	client := resty.New()
	client.SetBaseURL("https://api.geoapify.com")
	client.SetTimeout(10 * time.Second)

	return &geoapifyHttpRepo{
		client: client,
		apiKey: cfg.GeoapifyAPIKey,
	}
}

func (r *geoapifyHttpRepo) GetGeoapifys(
	ctx context.Context,
	city string,
) ([]domain.Geoapify, error) {

	// 1. Cari kota dulu untuk mendapatkan place_id
	var geoRes dto.GeoapifyGeocodingResponse

	res, err := r.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"text":   city,
			"type":   "city",
			"filter": "countrycode:id",
			"limit":  "1",
			"format": "json",
			"apiKey": r.apiKey,
		}).
		SetResult(&geoRes).
		Get("/v1/geocode/search")

	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, errors.New(res.Status())
	}

	if len(geoRes.Results) == 0 {
		return []domain.Geoapify{}, nil
	}

	placeID := geoRes.Results[0].PlaceID

	// 2. Cari rumah sakit di area kota tersebut
	var placesRes dto.GeoapifyPlacesResponse

	res, err = r.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"categories": "healthcare.geoapify",
			"filter":     fmt.Sprintf("place:%s", placeID),
			"limit":      "20",
			"lang":       "id",
			"apiKey":     r.apiKey,
		}).
		SetResult(&placesRes).
		Get("/v2/places")

	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, errors.New(res.Status())
	}

	geoapifys := make([]domain.Geoapify, 0)

	for _, feature := range placesRes.Features {
		prop := feature.Properties

		geoapifys = append(geoapifys, domain.Geoapify{
			ExternalID: prop.PlaceID,
			Name:       prop.Name,
			City:       prop.City,
			Address:    prop.Formatted,
			Latitude:   prop.Latitude,
			Longitude:  prop.Longitude,
		})
	}

	return geoapifys, nil
}
