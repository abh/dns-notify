package main

import (
	"testing"
)

func TestFixupHost(t *testing.T) {

	tests := map[string]string{
		"127.0.0.1":        "127.0.0.1:53",
		"127.0.0.2:5353":   "127.0.0.2:5353",
		"[::1]":            "[::1]:53",
		"[1:2::3]:53":      "[1:2::3]:53",
		"a.example.com:53": "a.example.com:53",
		"b.example.com":    "b.example.com:53",
	}

	for in, expected := range tests {
		out, err := fixupHost(in)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if out != expected {
			t.Logf("For input '%s' expected '%s' but got '%s'\n", in, expected, out)
			t.Fail()
		}

	}

}
