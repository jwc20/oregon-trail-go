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
	cli.State.Trip.Mileage = 0
	cli.State.Flags.SouthPassMileage = false
	cli.State.Trip.TurnNumber = 0
	cli.State.Player.Cash = InitialCash
}

// shooting prompt ****************************************************************************************************

func (cli *CLI) PromptShootingLevel() bool {
	cli.printf("HOW GOOD A SHOT ARE YOU WITH YOUR RIFLE?\n")
	cli.printf("  (1) ACE MARKSMAN\n")
	cli.printf("  (2) GOOD SHOT\n")
	cli.printf("  (3) FAIR TO MIDDLIN'\n")
	cli.printf("  (4) NEED MORE PRACTICE\n")
	cli.printf("  (5) SHAKY KNEES\n")
	cli.printf("ENTER ONE OF THE ABOVE: ")

	line, ok := cli.readLine()
	if !ok {
		return false
	}
	level, err := strconv.Atoi(strings.TrimSpace(line))
	//fmt.Printf("\n")
	//fmt.Fprint(cli.out, level)

	if err != nil || level < 1 || level > 5 {
		cli.printf("INVALID CHOICE\n")
		return false
	}
	cli.State.Player.ShootingLevel = level
	return true
}

// initial purchase logic *********************************************************************************************

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
	for {
		cli.printf("HOW MUCH DO YOU WANT TO SPEND ON YOUR OXEN TEAM? ")

		line, ok := cli.readLine()
		if !ok {
			return 0, false
		}
		oxen, err := strconv.Atoi(strings.TrimSpace(line))

		if err != nil || oxen < 200 || oxen > 300 {
			cli.printf("AMOUNT MUST BE BETWEEN $200 AND $300\n")
			continue
		}
		cli.State.Inventory.Oxen = oxen
		return oxen, true
	}
}

func (cli *CLI) FoodPurchase() (int, bool) {
	for {
		cli.printf("HOW MUCH DO YOU WANT TO SPEND ON FOOD? ")

		line, ok := cli.readLine()
		if !ok {
			return 0, false
		}
		food, err := strconv.Atoi(strings.TrimSpace(line))

		if err != nil || food < 100 || food > 200 {
			cli.printf("AMOUNT MUST BE BETWEEN $100 AND $200\n")
			continue
		}
		cli.State.Inventory.Food = food
		return food, true
	}
}

func (cli *CLI) AmmoPurchase() (int, bool) {
	for {
		cli.printf("HOW MUCH DO YOU WANT TO SPEND ON AMMO? ")

		line, ok := cli.readLine()
		if !ok {
			return 0, false
		}
		ammo, err := strconv.Atoi(strings.TrimSpace(line))

		if err != nil || ammo < 50 || ammo > 100 {
			cli.printf("AMOUNT MUST BE BETWEEN $50 AND $100\n")
			continue
		}
		cli.State.Inventory.Ammo = ammo
		return ammo, true
	}
}

func (cli *CLI) ClothingPurchase() (int, bool) {
	for {
		cli.printf("HOW MUCH DO YOU WANT TO SPEND ON CLOTHING? ")

		line, ok := cli.readLine()
		if !ok {
			return 0, false
		}
		clothing, err := strconv.Atoi(strings.TrimSpace(line))

		if err != nil || clothing < 50 || clothing > 100 {
			cli.printf("AMOUNT MUST BE BETWEEN $50 AND $100\n")
			continue
		}
		cli.State.Inventory.Clothing = clothing
		return clothing, true
	}
}

func (cli *CLI) MiscPurchase() (int, bool) {
	for {
		cli.printf("HOW MUCH DO YOU WANT TO SPEND ON MISCELLANEOUS ITEMS? ")

		line, ok := cli.readLine()
		if !ok {
			return 0, false
		}
		misc, err := strconv.Atoi(strings.TrimSpace(line))

		if err != nil || misc < 50 || misc > 100 {
			cli.printf("AMOUNT MUST BE BETWEEN $50 AND $100\n")
			continue
		}
		cli.State.Inventory.Miscellaneous = misc
		return misc, true
	}
}

// eating logic *********************************************************************************************************

func (cli *CLI) PromptEating() {
	cli.printf("DO YOU WANT TO EAT (1) POORLY (2) MODERATELY (3) WELL? ")

	line, _ := cli.readLine()
	choice, err := strconv.Atoi(strings.TrimSpace(line))

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

func (cli *CLI) HandleAilment() bool {
	if cli.State.Inventory.Miscellaneous < 5 {
		cli.printf("YOU DON'T HAVE ENOUGH SUPPLIES TO GET WELL.\n")
		if cli.State.Flags.Ill {
			cli.printf("YOU DIED OF PNEUMONIA.\n")
		} else {
			cli.printf("YOU DIED OF YOUR INJURIES.\n")
		}
		return false
	}
	cli.State.Inventory.Miscellaneous -= 5 + GetRandomInt(5)
	cli.State.Flags.Ill = false
	cli.State.Flags.Injured = false
	cli.printf("YOU USED MEDICINE AND RESTED.\n")
	return true
}

func (cli *CLI) GameLoop() {
	for cli.State.Trip.Mileage < TotalRequiredMileage {
		cli.State.Trip.TurnNumber++
		cli.State.Trip.CurrentDate = cli.State.Trip.TurnNumber

		cli.printf("\n---------------------------------------------\n")
		cli.printf("TURN %d — %s\n", cli.State.Trip.TurnNumber, cli.DateName())
		cli.printf("TOTAL MILEAGE: %d\n", cli.State.Trip.Mileage)
		cli.printf("FOOD: %d  AMMO: %d  CLOTHING: %d  MISC: %d\n",
			cli.State.Inventory.Food, cli.State.Inventory.Ammo,
			cli.State.Inventory.Clothing, cli.State.Inventory.Miscellaneous)
		cli.printf("CASH: $%d\n", cli.State.Player.Cash)
		cli.printf("---------------------------------------------\n")

		if cli.State.Inventory.Food < 0 {
			cli.printf("YOU RAN OUT OF FOOD AND STARVED TO DEATH.\n")
			cli.PrintFinalMessage()
			return
		}

		if !cli.PromptTurnAction() {
			return
		}

		cli.PromptEating()
		cli.AdvanceMileage()
		cli.GenerateEvent()

		if cli.State.Flags.Injured || cli.State.Flags.Ill {
			if !cli.HandleAilment() {
				cli.PrintFinalMessage()
				return
			}
		}
	}

	cli.printf("\n*** CONGRATULATIONS! YOU MADE IT TO OREGON CITY! ***\n")
	cli.PrintFinalMessage()
}

func (cli *CLI) PlaySVT() {
	cli.printf("DO YOU NEED INSTRUCTIONS (YES/NO)? ")
	answer, _ := cli.readLine()
	if strings.ToUpper(answer) == "YES" {
		cli.printIntro()
	}

	cli.InitSVT()

	if !cli.PromptShootingLevel() {
		return
	}

	if !cli.PromptInitialPurchases() {
		return
	}

	cli.GameLoop()
}

func (cli *CLI) printIntro() {
	cli.printf("\n")
	cli.printf("THIS PROGRAM SIMULATES A TRIP OVER THE OREGON TRAIL FROM\n")
	cli.printf("INDEPENDENCE, MISSOURI TO OREGON CITY, OREGON IN 1847.\n")
	cli.printf("YOUR FAMILY OF FIVE WILL COVER THE 2040 MILE OREGON TRAIL\n")
	cli.printf("IN 5-6 MONTHS --- IF YOU MAKE IT ALIVE.\n")
	cli.printf("\n")
	cli.printf("YOU HAD SAVED $900 TO SPEND FOR THE TRIP, AND YOU'VE JUST\n") // TODO: this is hard coded
	cli.printf("PAID $200 FOR A WAGON.\n")
	cli.printf("YOU WILL NEED TO SPEND THE REST OF YOUR MONEY ON THE\n")
	cli.printf("FOLLOWING ITEMS:\n")
	cli.printf("\n")
	cli.printf("  OXEN - YOU CAN SPEND $200-$300 ON YOUR TEAM.\n")
	cli.printf("  FOOD - THE MORE YOU HAVE, THE LESS CHANCE OF GETTING SICK.\n")
	cli.printf("  AMMUNITION - $1 BUYS A GENEROUS SUPPLY OF BULLETS.\n")
	cli.printf("  CLOTHING - IMPORTANT FOR COLD WEATHER IN THE MOUNTAINS.\n")
	cli.printf("  MISCELLANEOUS SUPPLIES - MEDICINE AND REPAIR ITEMS.\n")
	cli.printf("\n")
}

// helper functions ***************************************************************************************************

//func (cli *CLI) readLine() string {
//	cli.in.Scan()
//	return cli.in.Text()
//}

func (cli *CLI) readLine() (string, bool) {
	if !cli.in.Scan() {
		return "", false
	}
	return cli.in.Text(), true
}

func (cli *CLI) printf(format string, a ...interface{}) {
	fmt.Fprintf(cli.out, format, a...)
}

func (cli *CLI) PromptTurnAction() bool {
	cli.printf("DO YOU WANT TO (1) CONTINUE ON TRAIL (2) HUNT? ")

	line, ok := cli.readLine()
	if !ok {
		return false
	}
	choice, err := strconv.Atoi(strings.TrimSpace(line))

	if err != nil {
		choice = 1
	}
	cli.State.Trip.ActionChoice = choice
	return true
}

func (cli *CLI) PrintFinalMessage() {
	cli.printf("\n--- FINAL STATS ---\n")
	cli.printf("TOTAL MILEAGE: %d / %d\n", cli.State.Trip.Mileage, TotalRequiredMileage)
	cli.printf("TURNS TAKEN: %d\n", cli.State.Trip.TurnNumber)
	cli.printf("CASH REMAINING: $%d\n", cli.State.Player.Cash)
}
