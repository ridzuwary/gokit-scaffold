package main

import (
	"strings"
	"testing"
)

func TestRunValidateMissingMarkerFails(t *testing.T) {
	dir := t.TempDir()

	code := run([]string{"validate", "--dir", dir})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestFormatPrintOutputContainsRequiredSections(t *testing.T) {
	out, err := formatPrintOutput(ToolName, ToolVersion)
	if err != nil {
		t.Fatalf("formatPrintOutput returned error: %v", err)
	}

	required := []string{
		"gokit-scaffold 0.1.0",
		"Template Packs",
		"- service-http",
		"Generated Outputs Tree (service-http)",
		".gokit-scaffold",
		"cmd",
		"internal",
		"Marker Schema Summary",
		"`tool`",
		"`template_pack`",
		"`spec.http_port`",
		"Example new command",
		"gokit-scaffold new --name hello-api --module github.com/acme/hello-api --http-port 8080",
	}

	for _, item := range required {
		if !strings.Contains(out, item) {
			t.Fatalf("expected output to contain %q, got:\n%s", item, out)
		}
	}
}
