package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/leonzag/hyprwin/pkg/hyprwin"
)

func main() {
	command, err := hyprwin.HandleCli()
	if errors.Is(err, hyprwin.ErrHelpRequested) {
		printHelp()
		os.Exit(0)
	} else if errors.Is(err, hyprwin.ErrVersionRequested) {
		printVersion()
		os.Exit(0)
	} else if err != nil {
		printErr(err)
		printHelp()
		os.Exit(1)
	}

	out, err := hyprwin.Dispatch(command)
	if err != nil {
		printErr(err)
		printVersion()
		os.Exit(1)
	}
	if out != "" {
		printOut(out)
	}
}

func printOut(out string) {
	out = strings.TrimSuffix(out, "\n")
	fmt.Fprintf(os.Stdout, "hyprwin: %s\n", out)
}

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "hyprwin [Error]: %s\n", err.Error())
}

func printHelp() {
	fmt.Println(hyprwin.Usage + "\n")
	printVersion()
}

func printVersion() {
	version := fmt.Sprintf("%s", hyprwin.Version)
	hyprCompat := fmt.Sprintf("tested with %s hyprland tag\n", hyprwin.HyprCompatibility)
	printOut(version)
	printOut(hyprCompat)
}
