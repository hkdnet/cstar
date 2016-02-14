package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

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
		ccCh := make(chan CommitCount)
		pwd, _ := os.Getwd()
		go search(pwd, ccCh)
		ccs := []CommitCount{}
		for {
			cc, ok := <-ccCh
			if !ok {
				break
			}
			ccs = append(ccs, cc)
		}
		printCommitCounts(ccs)
	}

	app.Run(os.Args)
}

func search(root string, ccCh chan CommitCount) {
	dirCh := make(chan string)
	go gitDirSearch(root, dirCh)
	go gitDirToLog(dirCh, ccCh)
}

type CommitCount struct {
	ProjectName string
	Count       []int
}

func gitDirSearch(root string, dirCh chan string) error {
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			if filepath.Base(rel) == ".git" {
				dirCh <- path
			}
			return nil
		})

	close(dirCh)
	return err
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

func printCommitCounts(ccs []CommitCount) {
	maxNameLength := maxProjectNameLength(ccs)
	for i := 0; i < maxNameLength+1; i++ {
		fmt.Print(" ")
	}
	fmt.Println("1 2 3 4 5 6 7")
	for _, cc := range ccs {
		pjLen := utf8.RuneCountInString(cc.ProjectName)
		padLen := maxNameLength + 1 - pjLen
		fmt.Print(cc.ProjectName)
		for i := 0; i < padLen; i++ {
			fmt.Print(" ")
		}
		for _, c := range cc.Count {
			switch {
			case c == 0:
				fmt.Printf(cSprint("red", "D"))
			case c < 9:
				fmt.Printf(cSprint("yellow", "M"))
			default:
				fmt.Printf(cSprint("green", "L"))
			}
			fmt.Print(" ")
		}
		fmt.Print("\n")
	}
}

func cSprint(color, msg string) string {
	var colNo string
	switch color {
	case "red":
		colNo = "31"
	case "green":
		colNo = "32"
	case "yellow":
		colNo = "33"
	case "blue":
		colNo = "34"
	case "magenta":
		colNo = "35"
	case "cyan":
		colNo = "36"
	case "white":
		colNo = "37"
	default:
		colNo = "30"
	}
	return fmt.Sprintf("\033[" + colNo + "m" + msg + "\033[m")
}

func maxProjectNameLength(ccs []CommitCount) int {
	max := 0
	for _, cc := range ccs {
		tmp := utf8.RuneCountInString(cc.ProjectName)
		if tmp > max {
			max = tmp
		}
	}
	return max
}
