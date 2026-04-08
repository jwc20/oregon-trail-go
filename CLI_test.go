package oregontrail_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	trail "github.com/jwc20/oregontrail"
)

func TestCLIInitSVT(t *testing.T) {
	store := &trail.StubGameStore{}
	out := &bytes.Buffer{}
	cli := trail.NewCLI(store, strings.NewReader(""), out)
	cli.InitSVT()

	assert.True(t, cli.State.Trip.FortAvailable, "FortAvailable should be true")
	assert.False(t, cli.State.Flags.Injured, "Injured should be false")
	assert.False(t, cli.State.Flags.ClearedBlueMtns, "ClearedBlueMtns should be false")
	assert.Equal(t, 2040, cli.State.Trip.Mileage, "Mileage should be 2040")
	assert.False(t, cli.State.Flags.SouthPassMileage, "SouthPassMileage should be false")
	assert.Equal(t, 0, cli.State.Trip.TurnNumber, "TurnNumber should be 0")
}
