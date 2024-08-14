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
    --version       Show program version

Dispatchers:
    movefocus       Moves the focus in a direction

    movewindow      Moves the active window in a direction or to a monitor.
                    For floating windows, moves the window to the screen edge in that direction
Directions:
    l,r,u,d         For left, right, up, down
    mon:<monitor>   Only for movefocus dispatcher`

var (
	ErrHelpRequested       = errors.New("help requested")
	ErrVersionRequested    = errors.New("version requested")
	ErrNotEnoughArgs       = errors.New("not enough arguments, expected 2")
	ErrTooManyArgs         = errors.New("too many arguments, expected 2")
	ErrIncorrectDispatcher = errors.New("incorrect dispatcher received")
	ErrIncorrectDirection  = errors.New("incorrect direction received")
)

type (
	DispatcherCmd string
	DirectionArg  string
)

var (
	dispatchers = []DispatcherCmd{
		DispatcherCmd("movefocus"),
		DispatcherCmd("movewindow"),
	}
	directions = []DirectionArg{
		DirectionArg("l"),
		DirectionArg("r"),
		DirectionArg("u"),
		DirectionArg("d"),
	}
)

func (dp DispatcherCmd) Str() string {
	return string(dp)
}

func (dp DispatcherCmd) IsValid() bool {
	return slices.Contains(dispatchers, dp)
}

func (dir DirectionArg) Str() string {
	return string(dir)
}

func (dir DirectionArg) ToMonitor() bool {
	return strings.HasPrefix(dir.Str(), "mon:")
}

func (dir DirectionArg) IsValid(dp DispatcherCmd) bool {
	isBaseDir := slices.Contains(directions, dir)
	switch dp {
	case DispatcherCmd("movefocus"):
		return isBaseDir
	case DispatcherCmd("movewindow"):
		return isBaseDir || dir.ToMonitor()
	}
	return false
}

type CommandRequest struct {
	dispatcher DispatcherCmd
	direction  DirectionArg
}

func helpRequested(args []string) bool {
	for _, help := range []string{"h", "-h", "--help", "help"} {
		if slices.Contains(args, help) {
			return true
		}
	}
	return false
}

func versionRequested(args []string) bool {
	for _, v := range []string{"v", "-v", "--version", "version"} {
		if slices.Contains(args, v) {
			return true
		}
	}
	return false
}

func HandleCli() (cmd *CommandRequest, err error) {
	args := os.Args[1:]

	if helpRequested(args) {
		return nil, ErrHelpRequested
	}

	if versionRequested(args) {
		return nil, ErrVersionRequested
	}

	if len(args) < 2 {
		fmt.Print(Usage)
		return nil, ErrNotEnoughArgs
	}

	if len(args) != 2 {
		return nil, ErrTooManyArgs
	}

	dp, dir := DispatcherCmd(args[0]), DirectionArg(args[1])
	if !dp.IsValid() {
		return nil, ErrIncorrectDispatcher
	}
	if !dir.IsValid(dp) {
		return nil, ErrIncorrectDirection
	}

	return &CommandRequest{dp, dir}, nil
}
