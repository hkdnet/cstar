package main

import (
	"os"
	"testing"
)

func TestGitDirSearch(t *testing.T) {
	pwd, _ := os.Getwd()
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
	}
	expected1 := pwd + "/test/test_dir1/.git"
	expected2 := pwd + "/test/test_dir2/.git"
	if (dirs[0] != expected1 && dirs[1] != expected1) ||
		(dirs[0] != expected2 && dirs[1] != expected2) {
		t.Errorf("length of dirs should be \n - %s\n - %s\n but actually \n - %s\n - %s", expected1, expected2, dirs[0], dirs[1])

	}
}
