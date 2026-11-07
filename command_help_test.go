package main

import (
	"strings"
	"testing"
)

func TestCommandHelpPrintsUsage(t *testing.T) {
	var err error
	out := captureStdout(t, func() { err = command_help(nil) })
	if err != nil {
		t.Fatalf("command_help: %v", err)
	}
	if !strings.Contains(out, "help") || !strings.Contains(out, "exit") {
		t.Errorf("expected usage message listing help and exit, got %q", out)
	}
}
