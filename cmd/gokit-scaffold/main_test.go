package main

import (
	"testing"
)

func TestRunValidateMissingMarkerFails(t *testing.T) {
	dir := t.TempDir()

	code := run([]string{"validate", "--dir", dir})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}
