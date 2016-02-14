package command

import (
	"os"
	"testing"
	"time"
)

func TestGitDirSearch(t *testing.T) {
	pwd, _ := os.Getwd()
	if err := beforeTestGitDirSearch(pwd); err != nil {
		t.Fatalf("Got Error: %v", err)
		return
	}
	ch := make(chan string)
	go gitDirSearch(pwd+"/test", ch)
	dirs := []string{}
	for {
		absPath, ok := <-ch
		if !ok {
			break
		}
		dirs = append(dirs, absPath)
	}
	if len(dirs) != 2 {
		t.Errorf("length of dirs should be 2 but actually %d", len(dirs))
		for _, s := range dirs {
			t.Errorf(s)
		}
		return
	}
	expected1 := pwd + "/test/test_dir1/.git"
	expected2 := pwd + "/test/test_dir2/.git"
	if !((dirs[0] == expected1 || dirs[1] == expected1) && (dirs[0] == expected2 || dirs[1] == expected2)) && (dirs[0] != dirs[1]) {
		t.Errorf("dirs should be \n - %s\n - %s\n but actually \n - %s\n - %s", expected1, expected2, dirs[0], dirs[1])
	}
	if err := afterTestGitDirSearch(pwd); err != nil {
		t.Fatalf("Got Error: %v", err)
	}
}
func beforeTestGitDirSearch(pwd string) error {
	for _, path := range []string{"/test/test_dir1/.git", "/test/test_dir2/.git"} {
		if err := os.MkdirAll(pwd+path, 0755); err != nil {
			return err
		}
	}
	return nil
}
func afterTestGitDirSearch(pwd string) error {
	if err := os.RemoveAll(pwd + "/test"); err != nil {
		return err
	}
	return nil
}

func TestSumsGroupByDate(t *testing.T) {
	logs := "key1 hoge\nkey2 fuga^\nkey1 piyo"
	ret := sumsGroupByDate(logs)
	if val1, ok := ret["key1"]; !ok {
		t.Errorf("map should have key1: %v", ret)
	} else if val1 != 2 {
		t.Errorf("the value of key1 should be 2, but %v", val1)
	}
}
func TestSumsToArray(t *testing.T) {
	today := time.Now()
	m := map[string]int{}
	length := 3
	m[today.Format("2006-01-02")] = 1
	m[today.AddDate(0, 0, -1).Format("2006-01-02")] = 2
	m[today.AddDate(0, 0, -2).Format("2006-01-02")] = 3
	m[today.AddDate(0, 0, -3).Format("2006-01-02")] = 4
	ret := sumsToArray(m, length)
	if got, want := len(ret), length; got != want {
		t.Errorf("got: %v\nwant: %v", got, want)
	}
	if got, want := ret[0], 1; got != want {
		t.Errorf("got: %v\nwant: %v", got, want)
	}
	if got, want := ret[1], 2; got != want {
		t.Errorf("got: %v\nwant: %v", got, want)
	}
	if got, want := ret[2], 3; got != want {
		t.Errorf("got: %v\nwant: %v", got, want)
	}
}
