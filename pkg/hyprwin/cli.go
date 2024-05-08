package hyprwin

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

const Usage string = `Usage:
    hyprwin DISPATCHER DIRECTION

Flag:
    --help          Show this message

Dispatchers:
    movefocus       Moves the focus in a direction

    movewindow      Moves the active window in a direction or to a monitor.
                    For floating windows, moves the window to the screen edge in that direction
Directions:
    l,r,u,d         For left, right, up, down
    mon:<monitor>   Only for movefocus dispatcher`

var (
	ErrHelpRequested       = errors.New("help requested")
	ErrNotEnoughArgs       = errors.New("not enough arguments, expected 2")
	ErrTooManyArgs         = errors.New("too many arguments, expected 2")
	ErrIncorrectDispatcher = errors.New("incorrect dispatcher received")
	ErrIncorrectDirection  = errors.New("incorrect direction received")
)

type (
	dispatcher string
	direction  string
)

var (
	dispatchers = []dispatcher{
		dispatcher("movefocus"),
		dispatcher("movewindow"),
	}
	directions = []direction{
		direction("l"),
		direction("r"),
		direction("u"),
		direction("d"),
	}
)

func (dp dispatcher) Str() string {
	return string(dp)
}

func (dp dispatcher) IsValid() bool {
	return slices.Contains(dispatchers, dp)
}

func (dir direction) Str() string {
	return string(dir)
}

func (dir direction) ToMonitor() bool {
	return strings.HasPrefix(dir.Str(), "mon:")
}

func (dir direction) IsValid(dp dispatcher) bool {
	isBaseDir := slices.Contains(directions, dir)
	switch dp {
	case dispatcher("movefocus"):
		return isBaseDir
	case dispatcher("movewindow"):
		return isBaseDir || dir.ToMonitor()
	}
	return false
}

type command struct {
	dispatcher dispatcher
	direction  direction
}

func helpRequested(args []string) bool {
	for _, help := range []string{"h", "-h", "--help", "help"} {
		if slices.Contains(args, help) {
			return true
		}
	}
	return false
}

func HandleCli() (cmd *command, err error) {
	args := os.Args[1:]

	if helpRequested(args) {
		return nil, ErrHelpRequested
	}

	if len(args) < 2 {
		fmt.Print(Usage)
		return nil, ErrNotEnoughArgs
	}

	if len(args) != 2 {
		return nil, ErrTooManyArgs
	}

	dp, dir := dispatcher(args[0]), direction(args[1])
	if !dp.IsValid() {
		return nil, ErrIncorrectDispatcher
	}
	if !dir.IsValid(dp) {
		return nil, ErrIncorrectDirection
	}

	return &command{dp, dir}, nil
}
