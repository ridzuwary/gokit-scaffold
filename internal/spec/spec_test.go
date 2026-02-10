package spec

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestProjectSpecValidate(t *testing.T) {
	baseDir := t.TempDir()
	goodDir := filepath.Join(baseDir, "new-project")

	cases := []struct {
		name      string
		spec      ProjectSpec
		wantError bool
	}{
		{
			name: "valid spec",
			spec: ProjectSpec{
				Name:     "hello-api",
				Module:   "github.com/example/hello-api",
				Dir:      goodDir,
				HTTPPort: 8080,
			},
		},
		{
			name: "invalid module",
			spec: ProjectSpec{
				Name:     "hello-api",
				Module:   "github..com/example/hello-api",
				Dir:      filepath.Join(baseDir, "bad-module"),
				HTTPPort: 8080,
			},
			wantError: true,
		},
		{
			name: "invalid port",
			spec: ProjectSpec{
				Name:     "hello-api",
				Module:   "github.com/example/hello-api",
				Dir:      filepath.Join(baseDir, "bad-port"),
				HTTPPort: 70000,
			},
			wantError: true,
		},
		{
			name: "non-empty directory",
			spec: ProjectSpec{
				Name:     "hello-api",
				Module:   "github.com/example/hello-api",
				Dir:      filepath.Join(baseDir, "occupied"),
				HTTPPort: 8080,
			},
			wantError: true,
		},
	}

	if err := os.MkdirAll(filepath.Join(baseDir, "occupied"), 0o755); err != nil {
		t.Fatalf("mkdir occupied: %v", err)
	}
	if err := os.WriteFile(filepath.Join(baseDir, "occupied", "existing.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("seed occupied dir: %v", err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.spec.Validate()
			if tc.wantError && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateScaffoldDir(t *testing.T) {
	t.Run("missing marker", func(t *testing.T) {
		dir := t.TempDir()

		errs := ValidateScaffoldDir(dir)
		if len(errs) == 0 {
			t.Fatalf("expected validation errors, got none")
		}
	})

	t.Run("valid marker and file set", func(t *testing.T) {
		dir := t.TempDir()
		marker := Marker{
			Tool:         "gokit-scaffold",
			Version:      "0.1.0",
			TemplatePack: TemplatePackName,
			Spec: MarkerSpec{
				Name:     "hello-api",
				Module:   "github.com/example/hello-api",
				HTTPPort: 8080,
			},
		}

		content, err := json.MarshalIndent(marker, "", "  ")
		if err != nil {
			t.Fatalf("marshal marker: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, MarkerFileName), content, 0o644); err != nil {
			t.Fatalf("write marker: %v", err)
		}

		for _, rel := range RequiredScaffoldFiles {
			path := filepath.Join(dir, rel)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatalf("mkdir %s: %v", path, err)
			}
			if err := os.WriteFile(path, []byte("ok"), 0o644); err != nil {
				t.Fatalf("write %s: %v", path, err)
			}
		}

		errs := ValidateScaffoldDir(dir)
		if len(errs) > 0 {
			t.Fatalf("expected no validation errors, got %d (%v)", len(errs), errs[0])
		}
	})

	t.Run("allows extra files", func(t *testing.T) {
		dir := t.TempDir()
		marker := Marker{
			Tool:         "gokit-scaffold",
			Version:      "0.1.0",
			TemplatePack: TemplatePackName,
			Spec: MarkerSpec{
				Name:     "hello-api",
				Module:   "github.com/example/hello-api",
				HTTPPort: 8080,
			},
		}

		content, err := json.MarshalIndent(marker, "", "  ")
		if err != nil {
			t.Fatalf("marshal marker: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, MarkerFileName), content, 0o644); err != nil {
			t.Fatalf("write marker: %v", err)
		}

		for _, rel := range RequiredScaffoldFiles {
			path := filepath.Join(dir, rel)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatalf("mkdir %s: %v", path, err)
			}
			if err := os.WriteFile(path, []byte("ok"), 0o644); err != nil {
				t.Fatalf("write %s: %v", path, err)
			}
		}

		extra := filepath.Join(dir, "docs", "notes.md")
		if err := os.MkdirAll(filepath.Dir(extra), 0o755); err != nil {
			t.Fatalf("mkdir extra dir: %v", err)
		}
		if err := os.WriteFile(extra, []byte("user-owned"), 0o644); err != nil {
			t.Fatalf("write extra file: %v", err)
		}

		errs := ValidateScaffoldDir(dir)
		if len(errs) > 0 {
			t.Fatalf("expected no validation errors with extra files, got %d (%v)", len(errs), errs[0])
		}
	})
}
