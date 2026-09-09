package dto

type GeoapifyGeocodingResponse struct {
	Results []struct {
		PlaceID string `json:"place_id"`
	} `json:"results"`
}

type GeoapifyPlacesResponse struct {
	Features []struct {
		Properties struct {
			PlaceID   string  `json:"place_id"`
			Name      string  `json:"name"`
			City      string  `json:"city"`
			Formatted string  `json:"formatted"`
			Latitude  float64 `json:"lat"`
			Longitude float64 `json:"lon"`
		} `json:"properties"`
	} `json:"features"`
}
