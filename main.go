package main

import (
	"fmt"
	"os"
	"path/filepath"

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
		ch := make(chan string)
		pwd, _ := os.Getwd()
		fmt.Println(pwd)
		go getGitLists(pwd, ch)
		for {
			path, ok := <-ch
			if !ok {
				return
			}
			fmt.Println(path)
		}
	}

	app.Run(os.Args)
}

func getGitLists(root string, ch chan string) {
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				close(ch)
				return nil
			}
			if filepath.Base(rel) == ".git" {
				path, _ := filepath.Abs(rel)
				ch <- path
				return nil
			}
			return nil
		})

	if err != nil {
		fmt.Println(1, err)
	}
	close(ch)
}
