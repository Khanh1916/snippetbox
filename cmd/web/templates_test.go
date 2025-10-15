package main

import (
	"testing"
	"time"
)

func TestHumaDate(t *testing.T) {
	tm := time.Date(2025, 10, 15, 14, 20, 0, 0, time.UTC)
	hd := humanDate(tm)

	if hd != "15 Oct 2025 at 14:20" {
		t.Errorf("got %q; want %q", hd, "15 Oct 2025 at 14:20")
	}
}
