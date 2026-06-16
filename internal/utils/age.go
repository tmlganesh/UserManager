package utils

import "time"

// CalculateAge returns a person's age in whole years based on their date of birth.
// It correctly handles:
//   - Birthday already passed this year → full year difference
//   - Birthday is today → full year difference (birthday counts)
//   - Birthday not yet reached → year difference minus one
//
// Month/day comparison is used instead of YearDay() to avoid
// leap-year boundary issues (e.g., Mar 1 in leap vs non-leap years).
func CalculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()

	// Check if the birthday has not yet occurred this year.
	if now.Month() < dob.Month() ||
		(now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}

	return years
}
