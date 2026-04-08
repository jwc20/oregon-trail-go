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
	cli.State.Flags.ClearedBlueMtns = false
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
	//cash := cli.State.Player.Cash
	cli.printf("\nYOU HAVE $%d TO SPEND ON YOUR TRIP.\n", cli.State.Player.Cash)

	cli.OxenPurchase()
	cli.FoodPurchase()
	cli.AmmoPurchase()
	cli.ClothingPurchase()
	cli.MiscPurchase()

	return true
}

func (cli *CLI) OxenPurchase() bool {
	cli.printf("HOW MUCH DO YOU WANT TO SPEND ON YOUR OXEN TEAM? ")
	oxen, err := strconv.Atoi(strings.TrimSpace(cli.readLine()))
	if err != nil || oxen < 200 || oxen > 300 {
		cli.printf("AMOUNT MUST BE BETWEEN $200 AND $300\n")
		return false
	}
	cli.State.Inventory.Oxen = oxen
	cli.State.Player.Cash -= oxen
	return true
}

func (cli *CLI) FoodPurchase() bool {
	cli.printf("HOW MUCH DO YOU WANT TO SPEND ON FOOD? ")
	food, err := strconv.Atoi(strings.TrimSpace(cli.readLine()))
	if err != nil || food < 100 || food > 200 {
		cli.printf("AMOUNT MUST BE BETWEEN $100 AND $200\n")
		return false
	}
	cli.State.Inventory.Food = food
	cli.State.Player.Cash -= food
	return true
}

func (cli *CLI) AmmoPurchase() bool {
	cli.printf("HOW MUCH DO YOU WANT TO SPEND ON AMMO? ")
	ammo, err := strconv.Atoi(strings.TrimSpace(cli.readLine()))
	if err != nil || ammo < 50 || ammo > 100 {
		cli.printf("AMOUNT MUST BE BETWEEN $50 AND $100\n")
		return false
	}
	cli.State.Inventory.Ammo = ammo
	cli.State.Player.Cash -= ammo
	return true
}

func (cli *CLI) ClothingPurchase() bool {
	cli.printf("HOW MUCH DO YOU WANT TO SPEND ON CLOTHING? ")
	clothing, err := strconv.Atoi(strings.TrimSpace(cli.readLine()))
	if err != nil || clothing < 50 || clothing > 100 {
		cli.printf("AMOUNT MUST BE BETWEEN $50 AND $100\n")
		return false
	}
	cli.State.Inventory.Clothing = clothing
	cli.State.Player.Cash -= clothing
	return true
}

func (cli *CLI) MiscPurchase() bool {
	cli.printf("HOW MUCH DO YOU WANT TO SPEND ON MISCELLANEOUS ITEMS? ")
	misc, err := strconv.Atoi(strings.TrimSpace(cli.readLine()))
	if err != nil || misc < 50 || misc > 100 {
		cli.printf("AMOUNT MUST BE BETWEEN $50 AND $100\n")
		return false
	}
	cli.State.Inventory.Miscellaneous = misc
	cli.State.Player.Cash -= misc
	return true
}

// helper functions ***************************************************************************************************

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func (cli *CLI) printf(format string, a ...interface{}) {
	fmt.Fprintf(cli.out, format, a...)
}
