package command

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

var dayCount int

func CmdList(c *cli.Context) {
  fs := c.Args()
  dayCount = c.Int("d")
  pwd, _ := os.Getwd()
  if len(fs) == 0 {
    fs = []string { pwd }
  }
  ccs := []CommitCount{}
  for _, f := range fs {
    os.Chdir(pwd)
    ccs = append(ccs, search(f)...)
  }
  printCommitCounts(ccs)
}


func search(root string) []CommitCount {
	ccCh := make(chan CommitCount)
	dirCh := make(chan string)
	go gitDirSearch(root, dirCh)
	go gitDirToLog(dirCh, ccCh)
	ccs := []CommitCount{}
	for {
		cc, ok := <-ccCh
		if !ok {
			break
		}
		ccs = append(ccs, cc)
	}
  return ccs
}

type CommitCount struct {
	ProjectName string
	Count       []int
}

func (cc *CommitCount) ToStar(maxLen int) string {
	pjLen := utf8.RuneCountInString(cc.ProjectName)
	padLen := maxLen + 1 - pjLen
	ret := cc.ProjectName
	for i := 0; i < padLen; i++ {
		ret += " "
	}
	for _, c := range cc.Count {
		ret += fmt.Sprintf("%s ", countToStar(c))
	}
	return strings.Trim(ret, " ")
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
		counts[len-1-i] = m[key]
	}
	return counts
}

func printCommitCounts(ccs []CommitCount) {
	maxNameLength := maxProjectNameLength(ccs)
	printHeader(maxNameLength)
	printBody(maxNameLength, ccs)
}
func printHeader(length int) {
	for i := 0; i < length+1; i++ {
		fmt.Print(" ")
	}
  msg := ""
  for i := 0; i < dayCount; i++ {
    tmp := i
    for tmp >= 10 {
      tmp -= 10
    }
    msg += fmt.Sprintf("%d ", tmp)
  }
	fmt.Println(strings.Trim(msg, ""))
}
func printBody(length int, ccs []CommitCount) {
	for _, cc := range ccs {
		fmt.Println(cc.ToStar(length))
	}
}

func countToStar(count int) string {
	switch {
	case count == 0:
		return cSprint("red", "D")
	case count < 9:
		return cSprint("yellow", "M")
	default:
		return cSprint("green", "L")
	}
}

func cSprint(color, msg string) string {
	colNo := consoleColorNo(color)
	return fmt.Sprintf("\033[%dm%s\033[m", colNo, msg)
}

func consoleColorNo(cName string) int {
	switch cName {
	case "red":
		return 31
	case "green":
		return 32
	case "yellow":
		return 33
	case "blue":
		return 34
	case "magenta":
		return 35
	case "cyan":
		return 36
	case "white":
		return 37
	default:
		return 30
	}
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
