package utils

import (
	"github.com/rifkhia/lms-remake/internal/pkg"
	"time"
)

func ConvertTimes(times string) (string, pkg.CustomError) {
	t, err := time.Parse(time.RFC3339, times)
	if err != nil {
		return "", pkg.CustomError{
			Cause:   err,
			Code:    INTERNAL_SERVER_ERROR,
			Service: MODEL_SERVICE,
		}
	}

	timeOnly := t.Format("15:04")

	return timeOnly, pkg.CustomError{}
}

func ConvertDaysToInt(day string) string {
	dayMap := map[string]string{
		"monday":    "1",
		"tuesday":   "2",
		"wednesday": "3",
		"thursday":  "4",
		"friday":    "5",
		"saturday":  "6",
		"sunday":    "7",
	}

	return dayMap[day]
}

func ConvertIntToDay(day string) string {
	dayMap := map[string]string{
		"1": "monday",
		"2": "tuesday",
		"3": "wednesday",
		"4": "thursday",
		"5": "friday",
		"6": "saturday",
		"7": "sunday",
	}

	return dayMap[day]
}

func ContainsInt(array []int, element int) bool {
	for _, i := range array {
		if i == element {
			return true
		}
	}
	return false
}
