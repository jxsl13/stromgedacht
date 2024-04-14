package client

import (
	"testing"
	"time"
)

func TestClient_GetNow(t *testing.T) {
	c, _ := New()
	state, err := c.GetNow("68309") // Mannheim
	if err != nil {
		t.Fatal(err)
	}

	if state == nil {
		t.Fatal("now state is nil")
	}

	value := state.State
	if 4 < value && value != -1 {
		t.Fatalf("state enum value is out of bounds: 1 <= value <= 4 or equal to -1: %v", value)
	}
}

func TestClient_GetStates(t *testing.T) {
	c, _ := New()

	now := time.Now().Truncate(time.Second)
	hour, min, sec := now.Clock()

	offset := time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second
	t.Logf("%v", offset)
	lb := now.AddDate(0, 0, -4).Add(-offset)
	ub := now.AddDate(0, 0, 3).Add(-offset - time.Second)

	states, err := c.GetStates("68309", lb, ub) // Mannheim
	if err != nil {
		t.Fatal(err)
	}
	if states == nil {
		t.Fatalf("states state is nil for: from %s to %s", lb, ub)
	}

	if states.States == nil {
		t.Fatalf("states.States is nil for: from %s to %s", lb, ub)
	}

	for _, state := range states.States {
		value := state.State
		if 4 < value && value != -1 {
			t.Fatalf("state enum value is out of bounds: 1 <= value <= 4: %v", value)
		}
	}
}

func TestClient_GetForecast(t *testing.T) {
	c, _ := New()

	now := time.Now().In(time.UTC)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lb := today.AddDate(0, 0, -2)
	ub := today.AddDate(0, 0, -1)
	minDate := lb
	maxDate := ub.AddDate(0, 0, 1)

	forecast, err := c.GetForecast("68309", lb, ub)
	if err != nil {
		t.Fatal(err)
	}

	if forecast == nil {
		t.Fatalf("forecast is nil for: from %s to %s", lb, ub)
	}
	if forecast.Load == nil {
		t.Fatalf("forecast.Load is nil for: from %s to %s", lb, ub)
	}
	if forecast.RenewableEnergy == nil {
		t.Fatalf("forecast.RenewableEnergy is nil for: from %s to %s", lb, ub)
	}
	if forecast.ResidualLoad == nil {
		t.Fatalf("forecast.ResidualLoad is nil for: from %s to %s", lb, ub)
	}
	if forecast.SuperGreenThreshold == nil {
		t.Fatalf("forecast.SuperGreenThreshold is nil for: from %s to %s", lb, ub)
	}

	for _, load := range forecast.Load {
		dt := load.DateTime
		if dt.Before(minDate) || dt.After(maxDate) {
			t.Fatalf("load value at %v is out of temporal bounds [%v, %v]", dt, minDate, maxDate)
		}
	}
}
