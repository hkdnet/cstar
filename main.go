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
		if raw_log, err := execGitLog(dayCount); err != nil {
			close(logCh)
			return
		} else if raw_log == "" {
			logCh <- CommitCount{pjName, make([]int, dayCount)}
			continue
		} else {
			m := sumsGroupByDate(raw_log)
			logCh <- CommitCount{pjName, sumsToArray(m, dayCount)}
		}
	}
}
func execGitLog(len int) (string, error) {
	since := time.Now().AddDate(0, 0, -len).Format(time.RFC3339)
	out, err := exec.Command("git", "reflog", "--oneline", "--date=short", "--pretty=format:%ad %an", "--since", since).Output()
	return string(out), err
}

func sumsGroupByDate(logs string) map[string]int {
	ret := map[string]int{}
	for _, log := range strings.Split(logs, "\n") {
		tmp := strings.Split(log, " ")
		// mapに存在しないときはdefault(int) -> 0
		count := ret[tmp[0]]
		ret[tmp[0]] = count + 1
	}
	return ret
}

// the tail is the most recent
func sumsToArray(m map[string]int, len int) []int {
	counts := make([]int, len)
	for i := 0; i < len; i += 1 {
		key := time.Now().AddDate(0, 0, i+1-len).Format("2006-01-02")
		counts[i] = m[key]
	}
	return counts
}
