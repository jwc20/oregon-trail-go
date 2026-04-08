package oregontrail

import (
	"bufio"
	"io"
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
	cli.State.Trip.Mileage = 0
	cli.State.Flags.SouthPassMileage = false
	cli.State.Trip.TurnNumber = 0
}

// helper functions ***********************************

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
