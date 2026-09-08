package helper

import (
	"strings"
)

func IsValidBloodType(bloodType string) bool {
	var validBloodTypes = map[string]bool{
		"a+": true, "a-": true,
		"b+": true, "b-": true,
		"ab+": true, "ab-": true,
		"o+": true, "o-": true,
	}

	return validBloodTypes[strings.ToLower(bloodType)]
}
