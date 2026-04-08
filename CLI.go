package oregontrail

import (
	"bufio"
	"io"
)

type CLI struct {
	State GameState
	in    *bufio.Scanner
	out   io.Writer
}

func NewCLI(in io.Reader, out io.Writer) *CLI {
	return &CLI{
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
