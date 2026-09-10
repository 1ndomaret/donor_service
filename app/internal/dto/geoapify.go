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

type GeoapifyRoutingResponse struct {
	Features []struct {
		Properties struct {
			Distance      float64 `json:"distance"`
			DistanceUnits string  `json:"distance_units"`
			Time          float64 `json:"time"`
		} `json:"properties"`
	} `json:"features"`
}

type GeoapifyRoutingRequest struct {
	OriginLat      float64 `json:"origin_lat"`
	OriginLon      float64 `json:"origin_lon"`
	DestinationLat float64 `json:"destination_lat"`
	DestinationLon float64 `json:"destination_lon"`
}
