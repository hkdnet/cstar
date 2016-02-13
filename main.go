package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	logCh := make(chan CommitCount)
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

type CommitCount struct {
	ProjectName string
	Count       []int
}
type Commit struct {
	ProjectName string
	At          time.Time
	Author      string
}

func (c Commit) String() string {
	return fmt.Sprintf("%s %s %s", c.ProjectName, c.At.Format("2006-01-02"), c.Author)
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
func gitDirToLog(dirCh chan string, logCh chan CommitCount) {
	for {
		dir, ok := <-dirCh
		if !ok {
			close(logCh)
			return
		}
		os.Chdir(dir + "/../")
		pwd, _ := os.Getwd()
		pjName := filepath.Base(pwd)
		dayCount := 7
		since := time.Now().AddDate(0, 0, -dayCount).Format(time.RFC3339)
		out, err := exec.Command("git", "reflog", "--oneline", "--date=short", "--pretty=format:%ad %an", "--since", since).Output()
		if err != nil {
			close(logCh)
			return
		}
		raw_log := string(out)
		if raw_log == "" {
			logCh <- CommitCount{pjName, make([]int, dayCount)}
			continue
		}
		counts := make([]int, dayCount)
		/*
			logs := strings.Split(raw_log, "\n")
			for _, log := range logs {
				 arr := strings.Split(log, " ")
				 at, _ := time.Parse("2006-01-02", arr[0])
			}
		*/
		logCh <- CommitCount{pjName, counts}
	}
}
func logToCountPerDay(logs string) map[string]int {
	ret := map[string]int{}
	for _, log := range strings.Split(logs, "\n") {
		tmp := strings.Split(log, " ")
		// mapに存在しないときはdefault(int) -> 0
		count := ret[tmp[0]]
		ret[tmp[0]] = count + 1
	}
	return ret
}
