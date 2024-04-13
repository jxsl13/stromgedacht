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

	value := *state.State
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

	for _, state := range *states.States {
		value := *state.State
		if 4 < value && value != -1 {
			t.Fatalf("state enum value is out of bounds: 1 <= value <= 4: %v", value)
		}
	}
}
