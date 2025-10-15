package main

import (
	"testing"
	"time"

	"github.com/Khanh1916/snippetbox/internal/assert"
)

func TestHumaDate(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2025, 10, 15, 14, 20, 0, 0, time.UTC),
			want: "15 Oct 2025 at 14:20",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2025, 10, 15, 13, 20, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "15 Oct 2025 at 13:20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			// Use the new assert.Equal() helper to compare the expected and actual values.
			assert.Equal(t, hd, tt.want)
		})
	}
}
