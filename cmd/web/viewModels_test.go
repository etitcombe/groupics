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
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "UTC",
			tm:   time.Date(2020, time.September, 30, 10, 30, 0, 0, time.UTC),
			want: "30 Sep 2020 at 10:30",
		},
		{
			name: "CET",
			tm:   time.Date(2020, time.September, 30, 10, 30, 0, 0, time.FixedZone("EST", -4*60*60)),
			want: "30 Sep 2020 at 14:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := humanDate(tt.tm)

			if got != tt.want {
				t.Errorf("want %q; got %q", tt.want, got)
			}
		})
	}
}
