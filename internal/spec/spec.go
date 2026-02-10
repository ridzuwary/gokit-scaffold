package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const MarkerFileName = ".gokit-scaffold"

type ProjectSpec struct {
	Name     string
	Module   string
	Dir      string
	HTTPPort int
}

var (
	nameRe   = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	moduleRe = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*(/[a-zA-Z0-9._-]+)+$`)
)

func (s *ProjectSpec) Validate() error {
	if !nameRe.MatchString(s.Name) {
		return errors.New("name must match ^[a-z][a-z0-9-]*$")
	}
	if !moduleRe.MatchString(s.Module) || strings.Contains(s.Module, "..") {
		return errors.New("invalid module path")
	}
	if s.HTTPPort <= 0 || s.HTTPPort > 65535 {
		return errors.New("invalid http port")
	}
	if err := validateDir(s.Dir); err != nil {
		return err
	}

	return nil
}

func validateDir(dir string) error {
	if strings.TrimSpace(dir) == "" {
		return errors.New("target directory is required")
	}

	cleaned := filepath.Clean(dir)
	if hasParentTraversal(cleaned) {
		return errors.New("target directory must not contain parent traversal")
	}

	abs, err := filepath.Abs(cleaned)
	if err != nil {
		return fmt.Errorf("resolve directory: %w", err)
	}
	if isRootPath(abs) {
		return errors.New("target directory cannot be filesystem root")
	}

	info, err := os.Stat(abs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("stat target directory: %w", err)
	}

	if !info.IsDir() {
		return errors.New("target path exists and is not a directory")
	}

	entries, err := os.ReadDir(abs)
	if err != nil {
		return fmt.Errorf("read target directory: %w", err)
	}
	if len(entries) > 0 {
		return errors.New("target directory is not empty")
	}

	return nil
}

func hasParentTraversal(cleaned string) bool {
	if cleaned == ".." {
		return true
	}

	sep := string(filepath.Separator)
	for _, part := range strings.Split(cleaned, sep) {
		if part == ".." {
			return true
		}
	}
	return false
}

func isRootPath(path string) bool {
	volume := filepath.VolumeName(path)
	rest := strings.TrimPrefix(path, volume)
	rest = filepath.Clean(rest)
	return rest == string(filepath.Separator)
}
