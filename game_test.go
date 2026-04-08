package oregontrail

import "testing"

func TestGameConstants(t *testing.T) {
	t.Run("total mileage is 2040", func(t *testing.T) {
		got := TotalRequiredMileage
		want := 2040

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})
}
