package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/hkdnet/cstar/command"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{
	{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "",
		Action:  command.CmdList,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file,f",
				Value: "",
				Usage: "Not implemented yet...",
			},
			cli.IntFlag{
				Name:  "day, d",
				Value: 7,
				Usage: "how many days you'd like to list up",
			},
		},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
