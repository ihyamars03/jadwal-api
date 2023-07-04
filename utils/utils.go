package utils

import "regexp"

func IsValidDay(day string) bool {
	all_days := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
	for _, d := range all_days {
		if d == day {
			return true
		}
	}
	return false
}

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
