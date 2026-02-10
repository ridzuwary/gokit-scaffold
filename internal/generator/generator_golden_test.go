package generator

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/ridzuwary/gokit-scaffold/internal/spec"
)

const updateGoldenEnv = "UPDATE_GOLDEN"

func TestGenerateGoldenHelloAPI(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "hello-api")
	project := spec.ProjectSpec{
		Name:     "hello-api",
		Module:   "github.com/example/hello-api",
		Dir:      outputDir,
		HTTPPort: 8080,
	}

	if err := project.Validate(); err != nil {
		t.Fatalf("validate spec: %v", err)
	}
	if err := Generate(project, "0.1.0"); err != nil {
		t.Fatalf("generate: %v", err)
	}

	goldenDir := filepath.Join("..", "..", "testdata", "golden", "hello-api")
	if shouldUpdateGolden() {
		if err := os.RemoveAll(goldenDir); err != nil {
			t.Fatalf("remove existing golden dir: %v", err)
		}
		if err := copyTree(outputDir, goldenDir); err != nil {
			t.Fatalf("update golden fixtures: %v", err)
		}
	}

	assertTreesEqual(t, outputDir, goldenDir)
}

func shouldUpdateGolden() bool {
	return os.Getenv(updateGoldenEnv) == "1"
}

func assertTreesEqual(t *testing.T, gotRoot, wantRoot string) {
	t.Helper()

	gotFiles := listFiles(t, gotRoot)
	wantFiles := listFiles(t, wantRoot)

	if len(gotFiles) != len(wantFiles) {
		t.Fatalf("file count mismatch: got %d files, want %d files", len(gotFiles), len(wantFiles))
	}

	for i := range gotFiles {
		if gotFiles[i] != wantFiles[i] {
			t.Fatalf("file list mismatch at index %d: got %s, want %s", i, gotFiles[i], wantFiles[i])
		}
	}

	for _, rel := range gotFiles {
		gotPath := filepath.Join(gotRoot, rel)
		wantPath := filepath.Join(wantRoot, rel)

		gotContent, err := os.ReadFile(gotPath)
		if err != nil {
			t.Fatalf("read generated file %s: %v", gotPath, err)
		}
		wantContent, err := os.ReadFile(wantPath)
		if err != nil {
			t.Fatalf("read golden file %s: %v", wantPath, err)
		}

		if !bytes.Equal(gotContent, wantContent) {
			t.Fatalf("content mismatch: %s", rel)
		}
	}
}

func listFiles(t *testing.T, root string) []string {
	t.Helper()

	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			t.Fatalf("golden directory missing: %s (run `%s=1 go test ./internal/generator -run TestGenerateGoldenHelloAPI` to create it)", root, updateGoldenEnv)
		}
		t.Fatalf("walk files under %s: %v", root, err)
	}

	sort.Strings(files)
	return files
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, content, 0o644)
	})
}
