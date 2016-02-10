package main

import (
	"os"
	"fmt"

	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "hkdnet"
	app.Email = ""
	app.Usage = ""

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound
	app.Action = func(c *cli.Context) {
			fmt.Println("hello world")
	}

	app.Run(os.Args)
}
