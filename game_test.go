package oregontrail

import "testing"

func TestInitSVT(t *testing.T) {
	state := &GameState{}
	InitSVT(state)

	if !state.Trip.FortAvailable {
		t.Error("FortAvailable should be true")
	}
	if state.Flags.Injured {
		t.Error("Injured should be false")
	}
	if state.Flags.Ill {
		t.Error("Ill should be false")
	}
	if state.Flags.ClearedSouthPass {
		t.Error("ClearedSouthPass should be false")
	}
	if state.Flags.ClearedBlueMtns {
		t.Error("ClearedBlueMtns should be false")
	}
	if state.Trip.Mileage != 0 {
		t.Errorf("Mileage: got %d want 0", state.Trip.Mileage)
	}
	if state.Flags.SouthPassMileage {
		t.Error("SouthPassMileage should be false")
	}
	if state.Trip.TurnNumber != 0 {
		t.Errorf("TurnNumber: got %d want 0", state.Trip.TurnNumber)
	}
}
