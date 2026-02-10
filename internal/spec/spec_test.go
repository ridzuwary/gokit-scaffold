package spec

import (
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
