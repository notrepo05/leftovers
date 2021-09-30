package main

import (
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/notrepo05/leftovers/app"
	"github.com/notrepo05/leftovers/commands"
)

type command interface {
	Execute(app.Options) error
}

var Version = "Dev"

func main() {
	log.SetFlags(0)

	var o app.Options
	parser := flags.NewParser(&o, flags.HelpFlag|flags.PrintErrors)
	remaining, err := parser.ParseArgs(os.Args)
	if err != nil {
		return
	}

	if o.Version {
		log.Printf("%s\n", Version)
		return
	}

	cmd := "delete"
	if len(remaining) > 1 {
		cmd = "types"
	}
	if o.DryRun {
		cmd = "list"
	}

	logger := app.NewLogger(os.Stdout, os.Stdin, o.NoConfirm, o.Debug)

	otherEnvVars := app.NewOtherEnvVars()
	otherEnvVars.LoadConfig(&o)

	l, err := GetLeftovers(logger, o)
	if err != nil {
		log.Fatalf("\n\n%s\n", err)
	}

	commandSet := map[string]command{}
	commandSet["delete"] = commands.NewDelete(l)
	commandSet["list"] = commands.NewList(l)
	commandSet["types"] = commands.NewTypes(l)

	err = commandSet[cmd].Execute(o)
	if err != nil {
		log.Fatalf("\n\n%s\n", err)
	}
}
