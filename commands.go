package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var GlobalFlags = []cli.Flag{
	cli.IntFlag{
		Name:  "day, d",
		Value: 7,
		Usage: "how many days you'd like to list up",
	},
}

var Commands = []cli.Command{}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
