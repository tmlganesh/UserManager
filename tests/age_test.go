package tests

import (
	"testing"
	"time"

	"github.com/ganesh/ainyx/internal/utils"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "birthday already passed this year",
			dob:      time.Date(2000, now.Month()-1, 15, 0, 0, 0, 0, time.UTC),
			expected: now.Year() - 2000,
		},
		{
			name:     "birthday is today",
			dob:      time.Date(2000, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			expected: now.Year() - 2000,
		},
		{
			name:     "birthday not yet reached this year",
			dob:      time.Date(2000, now.Month()+1, 15, 0, 0, 0, 0, time.UTC),
			expected: now.Year() - 2000 - 1,
		},
		{
			name:     "born on Jan 1st",
			dob:      time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: now.Year() - 1995,
		},
		{
			name: "born on Dec 31st",
			dob:  time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: func() int {
				age := now.Year() - 1990
				if now.Month() < 12 || (now.Month() == 12 && now.Day() < 31) {
					age--
				}
				return age
			}(),
		},
		{
			name: "leap year baby - Feb 29",
			dob:  time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: func() int {
				age := now.Year() - 2000
				if now.Month() < 2 || (now.Month() == 2 && now.Day() < 29) {
					age--
				}
				return age
			}(),
		},
		{
			name:     "newborn - born today",
			dob:      now,
			expected: 0,
		},
		{
			name:     "one year old - born exactly one year ago",
			dob:      now.AddDate(-1, 0, 0),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.CalculateAge(tt.dob)
			if got != tt.expected {
				t.Errorf("CalculateAge(%v) = %d, want %d", tt.dob.Format("2006-01-02"), got, tt.expected)
			}
		})
	}
}
