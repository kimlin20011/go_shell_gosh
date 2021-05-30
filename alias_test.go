package main

import (
	"testing"
)

// alias_test.go

func TestExpandAlias(t *testing.T) {
	// case1: 試試看 "ls -h" 有沒有被展開成 "ls -l -h"
	setAlias("ls", "ls -l")
	input, ans := "ls -h", "ls -l -h"

	if expandAlias(input) != ans {
		t.Fail()
	}

	// case2: 試試看 "gst" 有沒有被展開成 "git status"
	setAlias("gst", "git status")
	input, ans = "gst", "git status"

	if expandAlias(input) != ans {
		t.Fail()
	}
}
