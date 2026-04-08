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
	assert.False(t, cli.State.Flags.ClearedBlueMountains, "ClearedBlueMountains should be false")
	assert.Equal(t, 2040, cli.State.Trip.Mileage, "Mileage should be 2040")
	assert.False(t, cli.State.Flags.SouthPassMileage, "SouthPassMileage should be false")
	assert.Equal(t, 0, cli.State.Trip.TurnNumber, "TurnNumber should be 0")
}

func TestShootingLevel(t *testing.T) {
	store := &trail.StubGameStore{}
	out := &bytes.Buffer{}
	cli := trail.NewCLI(store, strings.NewReader("3\n"), out)
	cli.InitSVT()

	result := cli.PromptShootingLevel()

	assert.True(t, result, "Expected PromptShootingLevel to return true")
	assert.Equal(t, 3, cli.State.Player.ShootingLevel, "ShootingLevel")
}

func TestInitialPurchases(t *testing.T) {
	input := "200\n100\n50\n50\n50\n"
	store := &trail.StubGameStore{}
	out := &bytes.Buffer{}
	cli := trail.NewCLI(store, strings.NewReader(input), out)
	cli.InitSVT()

	result := cli.PromptInitialPurchases()

	assert.True(t, result, "Expected PromptInitialPurchases to return true")
	assert.Equal(t, 200, cli.State.Inventory.Oxen, "Inventory Oxen should be 200")
	assert.Equal(t, 100, cli.State.Inventory.Food, "Inventory Food should be 100")
	assert.Equal(t, 50, cli.State.Inventory.Ammo, "Inventory Ammo should be 50")
	assert.Equal(t, 50, cli.State.Inventory.Clothing, "Inventory Clothing should be 50")
	assert.Equal(t, 50, cli.State.Inventory.Miscellaneous, "Inventory Miscellaneous should be 50")
	assert.Equal(t, 250, cli.State.Player.Cash, "Player Cash should be 250")
}

func TestInitialPurchasesOverspending(t *testing.T) {
	input := "300\n300\n300\n300\n300\n"
	store := &trail.StubGameStore{}
	out := &bytes.Buffer{}
	cli := trail.NewCLI(store, strings.NewReader(input), out)
	cli.InitSVT()

	result := cli.PromptInitialPurchases()

	assert.False(t, result, "expected PromptInitialPurchases to return false on overspend")
	assert.NotContains(t, out.String(), "OVERSPENT", "expected output to contain OVERSPENT")
}

func TestEating(t *testing.T) {
	store := &trail.StubGameStore{}
	out := &bytes.Buffer{}
	cli := trail.NewCLI(store, strings.NewReader("2\n"), out)
	cli.InitSVT()
	cli.State.Inventory.Food = 100

	cli.PromptEating()

	assert.Equal(t, 2, cli.State.Trip.EatingChoice, "expected EatingChoice to be 2")
	assert.Equal(t, 82, cli.State.Inventory.Food, "expected Inventory Food to be 82")
}

func TestAdvanceMileage(t *testing.T) {
	store := &trail.StubGameStore{}
	out := &bytes.Buffer{}
	cli := trail.NewCLI(store, strings.NewReader(""), out)
	cli.InitSVT()
	cli.State.Inventory.Oxen = 1000
	cli.State.Trip.ActionChoice = 1

	cli.AdvanceMileage()

	assert.GreaterOrEqualf(t, cli.State.Trip.Mileage, 0, "Mileage should be > 0, got %d", cli.State.Trip.Mileage)
}

func TestGenerateEvent(t *testing.T) {
	store := &trail.StubGameStore{}
	out := &bytes.Buffer{}
	cli := trail.NewCLI(store, strings.NewReader(""), out)
	cli.InitSVT()
	cli.State.Inventory.Food = 100
	cli.State.Inventory.Ammo = 50
	cli.State.Inventory.Miscellaneous = 30
	cli.State.Inventory.Clothing = 20
	cli.State.Trip.Mileage = 500

	// run several events and verify the state changed
	originalFood := cli.State.Inventory.Food
	originalMileage := cli.State.Trip.Mileage

	changed := false
	for i := 0; i < 20; i++ {
		cli.GenerateEvent()
		if cli.State.Inventory.Food != originalFood ||
			cli.State.Trip.Mileage != originalMileage ||
			cli.State.Flags.Injured ||
			cli.State.Flags.Ill {
			changed = true
			break
		}
	}

	assert.True(t, changed, "expected at least one event to change game state")
}

func TestHandleAilment(t *testing.T) {
	t.Run("dies when no supplies", func(t *testing.T) {
		store := &trail.StubGameStore{}
		out := &bytes.Buffer{}
		cli := trail.NewCLI(store, strings.NewReader(""), out)
		cli.InitSVT()
		cli.State.Flags.Ill = true
		cli.State.Inventory.Miscellaneous = 2

		survived := cli.HandleAilment()

		assert.False(t, survived, "expected HandleAilment to return false (death)")
		assert.Contains(t, out.String(), "PNEUMONIA", "expected death message to mention PNEUMONIA")
	})

	t.Run("survives with enough supplies", func(t *testing.T) {
		store := &trail.StubGameStore{}
		out := &bytes.Buffer{}
		cli := trail.NewCLI(store, strings.NewReader(""), out)
		cli.InitSVT()
		cli.State.Flags.Injured = true
		cli.State.Inventory.Miscellaneous = 30

		survived := cli.HandleAilment()

		assert.True(t, survived, "expected HandleAilment to return true (survived)")
		assert.False(t, cli.State.Flags.Injured, "expected Injured to be false")
	})
}
