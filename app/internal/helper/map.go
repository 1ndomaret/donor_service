package helper

func IsValidCoord(lat, lon float64) bool {
	if lat < -90.0 || lat > 90.0 {
		return false
	}
	if lon < -180.0 || lon > 180.0 {
		return false
	}
	return true
}
