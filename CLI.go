package oregontrail

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type CLI struct {
	store GameStore
	State GameState
	in    *bufio.Scanner
	out   io.Writer
}

func NewCLI(store GameStore, in io.Reader, out io.Writer) *CLI {
	return &CLI{
		store: store,
		State: GameState{},
		in:    bufio.NewScanner(in),
		out:   out,
	}
}

func (cli *CLI) InitSVT() {
	cli.State.Trip.FortAvailable = true
	cli.State.Flags.Injured = false
	cli.State.Flags.Ill = false
	cli.State.Flags.ClearedSouthPass = false
	cli.State.Flags.ClearedBlueMountains = false
	cli.State.Trip.Mileage = TotalRequiredMileage
	cli.State.Flags.SouthPassMileage = false
	cli.State.Trip.TurnNumber = 0
	cli.State.Player.Cash = InitialCash
}

// Game Prompts *******************************************************************************************************

func (cli *CLI) PromptShootingLevel() bool {
	cli.printf("HOW GOOD A SHOT ARE YOU WITH YOUR RIFLE?\n")
	cli.printf("  (1) ACE MARKSMAN\n")
	cli.printf("  (2) GOOD SHOT\n")
	cli.printf("  (3) FAIR TO MIDDLIN'\n")
	cli.printf("  (4) NEED MORE PRACTICE\n")
	cli.printf("  (5) SHAKY KNEES\n")
	cli.printf("ENTER ONE OF THE ABOVE: ")

	level, err := strconv.Atoi(strings.TrimSpace(cli.readLine()))
	//fmt.Printf("\n")
	//fmt.Fprint(cli.out, level)

	if err != nil || level < 1 || level > 5 {
		cli.printf("INVALID CHOICE\n")
		return false
	}
	cli.State.Player.ShootingLevel = level
	return true
}

func (cli *CLI) PromptInitialPurchases() bool {
	startingCash := cli.State.Player.Cash
	cli.printf("\nYOU HAVE $%d TO SPEND ON YOUR TRIP.\n", startingCash)

	spent := 0
	purchases := []func() (int, bool){
		cli.OxenPurchase,
		cli.FoodPurchase,
		cli.AmmoPurchase,
		cli.ClothingPurchase,
		cli.MiscPurchase,
	}

	for _, purchase := range purchases {
		amount, ok := purchase()
		if !ok {
			return false
		}
		spent += amount
	}

	remaining := startingCash - spent
	if remaining < 0 {
		cli.printf("YOU OVERSPENT ON YOUR INVENTORY\n")
		return false
	}

	cli.State.Player.Cash = remaining
	cli.printf("AFTER ALL YOUR PURCHASES, YOU NOW HAVE $%d LEFT.\n", remaining)
	return true
}

func (cli *CLI) OxenPurchase() (int, bool) {
	cli.printf("HOW MUCH DO YOU WANT TO SPEND ON YOUR OXEN TEAM? ")
	oxen, err := strconv.Atoi(strings.TrimSpace(cli.readLine()))
	if err != nil || oxen < 200 || oxen > 300 {
		cli.printf("AMOUNT MUST BE BETWEEN $200 AND $300\n")
		return 0, false
	}
	cli.State.Inventory.Oxen = oxen
	return oxen, true
}

func (cli *CLI) FoodPurchase() (int, bool) {
	cli.printf("HOW MUCH DO YOU WANT TO SPEND ON FOOD? ")
	food, err := strconv.Atoi(strings.TrimSpace(cli.readLine()))
	if err != nil || food < 100 || food > 200 {
		cli.printf("AMOUNT MUST BE BETWEEN $100 AND $200\n")
		return 0, false
	}
	cli.State.Inventory.Food = food
	return food, true
}

func (cli *CLI) AmmoPurchase() (int, bool) {
	cli.printf("HOW MUCH DO YOU WANT TO SPEND ON AMMO? ")
	ammo, err := strconv.Atoi(strings.TrimSpace(cli.readLine()))
	if err != nil || ammo < 50 || ammo > 100 {
		cli.printf("AMOUNT MUST BE BETWEEN $50 AND $100\n")
		return 0, false
	}
	cli.State.Inventory.Ammo = ammo
	return ammo, true
}

func (cli *CLI) ClothingPurchase() (int, bool) {
	cli.printf("HOW MUCH DO YOU WANT TO SPEND ON CLOTHING? ")
	clothing, err := strconv.Atoi(strings.TrimSpace(cli.readLine()))
	if err != nil || clothing < 50 || clothing > 100 {
		cli.printf("AMOUNT MUST BE BETWEEN $50 AND $100\n")
		return 0, false
	}
	cli.State.Inventory.Clothing = clothing
	return clothing, true
}

func (cli *CLI) MiscPurchase() (int, bool) {
	cli.printf("HOW MUCH DO YOU WANT TO SPEND ON MISCELLANEOUS ITEMS? ")
	misc, err := strconv.Atoi(strings.TrimSpace(cli.readLine()))
	if err != nil || misc < 50 || misc > 100 {
		cli.printf("AMOUNT MUST BE BETWEEN $50 AND $100\n")
		return 0, false
	}
	cli.State.Inventory.Miscellaneous = misc
	return misc, true
}

func (cli *CLI) PromptEating() {
	cli.printf("DO YOU WANT TO EAT (1) POORLY (2) MODERATELY (3) WELL? ")
	choice, err := strconv.Atoi(strings.TrimSpace(cli.readLine()))
	if err != nil || choice < 1 || choice > 3 {
		choice = 2
	}
	cli.State.Trip.EatingChoice = choice

	foodUsed := 8 + 5*choice // TODO: check later if this is correct
	cli.State.Inventory.Food -= foodUsed
}

// trip logic *********************************************************************************************************
//1. continue
//2. hunt

func (cli *CLI) AdvanceMileage() {
	cli.State.Trip.PreviousMileage = cli.State.Trip.Mileage
	randomInt := GetRandomInt(10)
	miles := int(200 + (cli.State.Inventory.Oxen-220)/5 + randomInt*10)
	if cli.State.Trip.ActionChoice != 1 {
		miles /= 2
	}
	cli.State.Trip.Mileage += miles
}

func (cli *CLI) DateName() string {
	// Date logic is based on original Oregon Trail game
	dates := []string{
		"MARCH 29",
		"APRIL 12",
		"APRIL 26",
		"MAY 10",
		"MAY 24",
		"JUNE 7",
		"JUNE 21",
		"JULY 5",
		"JULY 19",
		"AUGUST 2",
		"AUGUST 16",
		"AUGUST 31",
		"SEPTEMBER 13",
		"SEPTEMBER 27",
		"OCTOBER 11",
		"OCTOBER 25",
		"NOVEMBER 8",
		"NOVEMBER 22",
		"DECEMBER 6",
		"DECEMBER 20",
	}

	weekdays := []string{
		"SATURDAY",
		"SUNDAY",
		"MONDAY",
		"TUESDAY",
		"WEDNESDAY",
		"THURSDAY",
		"FRIDAY",
	}

	turn := cli.State.Trip.TurnNumber

	if turn > len(dates) {
		return "WINTER"
	}

	idx := turn - 1
	date := dates[idx]
	weekday := weekdays[idx%len(weekdays)]

	return fmt.Sprintf("%s, %s, 1847", weekday, date)
}

func (cli *CLI) GenerateEvent() {
	r := GetRandomInt(100)

	switch {
	case r < 6:
		cli.printf("WAGON BREAKS DOWN — LOSS OF TIME AND SUPPLIES.\n")
		cli.State.Trip.Mileage -= 15 + GetRandomInt(10)
		cli.State.Inventory.Miscellaneous -= 8 + GetRandomInt(5)
	case r < 11:
		cli.printf("OX INJURED — LOSS OF TIME.\n")
		cli.State.Trip.Mileage -= 25
		cli.State.Inventory.Oxen -= 20
	case r < 15:
		cli.printf("BAD LUCK — YOUR DAUGHTER BROKE HER ARM.\n")
		cli.State.Flags.Injured = true
		cli.State.Inventory.Miscellaneous -= 5 + GetRandomInt(4)
	case r < 20:
		cli.printf("WILD ANIMALS ATTACK!\n")
		cli.State.Inventory.Ammo -= 10 + GetRandomInt(5)
		if cli.State.Inventory.Ammo < 0 {
			cli.printf("YOU RAN OUT OF BULLETS — THEY GOT SOME OF YOUR FOOD.\n")
			cli.State.Inventory.Food -= 30 + GetRandomInt(20)
		}
	case r < 25:
		cli.printf("COLD WEATHER — BRRRR!\n")
		if cli.State.Inventory.Clothing < 20 {
			cli.printf("YOU DON'T HAVE ENOUGH CLOTHING TO KEEP WARM.\n")
			cli.State.Flags.Ill = true
		}
	case r < 30:
		cli.printf("HEAVY RAINS — TIME LOST AND SUPPLIES DAMAGED.\n")
		cli.State.Trip.Mileage -= 10 + GetRandomInt(5)
		cli.State.Inventory.Food -= 10
		cli.State.Inventory.Ammo -= 5 + GetRandomInt(5)
		cli.State.Inventory.Miscellaneous -= 5 + GetRandomInt(5)
	case r < 33:
		cli.printf("BANDITS ATTACK!\n")
		cli.State.Inventory.Food -= 10 + GetRandomInt(10)
		cli.State.Player.Cash -= 10 + GetRandomInt(15)
		if cli.State.Player.Cash < 0 {
			cli.State.Player.Cash = 0
		}
	case r < 36:
		cli.printf("FIRE IN YOUR WAGON — LOSS OF SUPPLIES.\n")
		cli.State.Inventory.Food -= 40 + GetRandomInt(30)
		cli.State.Inventory.Ammo -= 20 + GetRandomInt(20)
		cli.State.Inventory.Miscellaneous -= 10 + GetRandomInt(10)
	case r < 40:
		cli.printf("HELPFUL INDIANS SHOW YOU WHERE TO FIND FOOD.\n")
		cli.State.Inventory.Food += 14 + GetRandomInt(5)
	default:
		// nothing happens
	}
}

// helper functions ***************************************************************************************************

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func (cli *CLI) printf(format string, a ...interface{}) {
	fmt.Fprintf(cli.out, format, a...)
}
