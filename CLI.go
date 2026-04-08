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
}

func (cli *CLI) PromptShootingLevel() bool {
	cli.printf("HOW GOOD A SHOT ARE YOU WITH YOUR RIFLE?\n")
	cli.printf("  (1) ACE MARKSMAN\n")
	cli.printf("  (2) GOOD SHOT\n")
	cli.printf("  (3) FAIR TO MIDDLIN'\n")
	cli.printf("  (4) NEED MORE PRACTICE\n")
	cli.printf("  (5) SHAKY KNEES\n")
	cli.printf("ENTER ONE OF THE ABOVE: ")

	level, err := strconv.Atoi(strings.TrimSpace(cli.readLine()))
	fmt.Printf("\n")
	fmt.Fprint(cli.out, level)

	if err != nil || level < 1 || level > 5 {
		cli.printf("INVALID CHOICE\n")
		return false
	}
	cli.State.Player.ShootingLevel = level
	return true
}

// helper functions ***********************************

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func (cli *CLI) printf(format string, a ...interface{}) {
	fmt.Fprintf(cli.out, format, a...)
}
