package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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
		go search(pwd, ch)
		for {
			output, ok := <-ch
			if !ok {
				return
			}
			fmt.Println(output)
		}
		os.Chdir(pwd)
	}

	app.Run(os.Args)
}

func search(root string, ch chan string) {
	dirCh := make(chan string)
	logCh := make(chan string)
	go gitDirSearch(root, dirCh)
	go gitDirToLog(dirCh, logCh)
	for {
		log, ok := <-logCh
		if !ok {
			break
		}
		fmt.Println(log)
	}
	close(ch)
}

func gitDirSearch(root string, dirCh chan string) {
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}
			if filepath.Base(rel) == ".git" {
				dirCh <- path
				return nil
			}
			return nil
		})

	if err != nil {
		fmt.Println(1, err)
	}
	close(dirCh)
}
func gitDirToLog(dirCh, logCh chan string) {
	for {
		dir, ok := <-dirCh
		if !ok {
			close(logCh)
			return
		}
		os.Chdir(dir + "/../")
		pwd, _ := os.Getwd()
		logCh <- fmt.Sprintf("...move to %s\n", pwd)
		since := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
		out, err := exec.Command("git", "log", "--oneline", "--since", since).Output()
		if err != nil {
			logCh <- err.Error()
			close(logCh)
			return
		}
		logCh <- string(out)
	}
}
