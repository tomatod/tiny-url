package main

import (
	"regexp"
	"testing"
)

func TestMakeRandomStr(t *testing.T) {
	str, err := MakeRandomStr(10)
	if err != nil {
		t.Fatal(err)
		return
	}
	if l := len(str); l != 10 {
		t.Fatalf("expected: 10  real: %d\n", l)
		return
	}
	match, err := regexp.Match("^[a-z|A-Z|0-9]{10}$", []byte(str))
	if err != nil {
		t.Fatal(err)
		return
	}
	if !match {
		t.Fatalf("%s is not matched.\n", str)
		return
	}
}
