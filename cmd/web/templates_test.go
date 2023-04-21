package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			"UTC",
			time.Date(2023, 04, 17, 20, 34, 58, 0, time.UTC),
			"17 Apr 2023 at 20:34",
		},
		{
			"CET",
			time.Date(2023, 04, 17, 20, 34, 58, 0, time.FixedZone("CET", 1*60*60)),
			"17 Apr 2023 at 19:34",
		},
		{
			"Empty",
			time.Time{},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := humanDate(tt.tm); got != tt.want {
				t.Errorf("humanDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
