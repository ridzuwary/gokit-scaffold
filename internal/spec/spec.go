package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const MarkerFileName = ".gokit-scaffold"
const TemplatePackName = "service-http"

var RequiredScaffoldFiles = []string{
	"README.md",
	"cmd/server/main.go",
	"go.mod",
	"internal/config/config.go",
	"internal/httpserver/health.go",
	"internal/httpserver/server.go",
	"internal/logging/logging.go",
}

type ProjectSpec struct {
	Name     string
	Module   string
	Dir      string
	HTTPPort int
}

type Marker struct {
	Tool         string     `json:"tool"`
	Version      string     `json:"version"`
	TemplatePack string     `json:"template_pack"`
	Spec         MarkerSpec `json:"spec"`
}

type MarkerSpec struct {
	Name     string `json:"name"`
	Module   string `json:"module"`
	HTTPPort int    `json:"http_port"`
}

var (
	nameRe   = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	moduleRe = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*(/[a-zA-Z0-9._-]+)+$`)
)

func (s *ProjectSpec) Validate() error {
	if err := validateName(s.Name); err != nil {
		return err
	}
	if err := validateModule(s.Module); err != nil {
		return err
	}
	if err := validateHTTPPort(s.HTTPPort); err != nil {
		return err
	}
	if err := validateDir(s.Dir); err != nil {
		return err
	}

	return nil
}

func ValidateScaffoldDir(dir string) []error {
	var errs []error
	absDir, err := filepath.Abs(filepath.Clean(dir))
	if err != nil {
		return []error{fmt.Errorf("resolve directory: %w", err)}
	}

	info, err := os.Stat(absDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []error{fmt.Errorf("directory not found: %s", absDir)}
		}
		return []error{fmt.Errorf("stat directory: %w", err)}
	}
	if !info.IsDir() {
		return []error{fmt.Errorf("path is not a directory: %s", absDir)}
	}

	markerPath := filepath.Join(absDir, MarkerFileName)
	marker, err := readMarker(markerPath)
	if err != nil {
		errs = append(errs, err)
	}

	if err == nil {
		errs = append(errs, validateMarker(marker)...)
	}

	for _, rel := range RequiredScaffoldFiles {
		path := filepath.Join(absDir, rel)
		info, statErr := os.Stat(path)
		if statErr != nil {
			if errors.Is(statErr, os.ErrNotExist) {
				errs = append(errs, fmt.Errorf("missing required file: %s (re-run `gokit-scaffold new` into a clean directory)", rel))
				continue
			}
			errs = append(errs, fmt.Errorf("stat required file %s: %w", rel, statErr))
			continue
		}
		if info.IsDir() {
			errs = append(errs, fmt.Errorf("required file is a directory: %s", rel))
		}
	}

	return errs
}

func readMarker(path string) (Marker, error) {
	var marker Marker

	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return marker, fmt.Errorf("missing %s marker at %s (run `gokit-scaffold new` to generate a scaffold first)", MarkerFileName, path)
		}
		return marker, fmt.Errorf("read marker %s: %w", path, err)
	}

	if err := json.Unmarshal(content, &marker); err != nil {
		return marker, fmt.Errorf("parse marker %s: %w", path, err)
	}

	return marker, nil
}

func validateMarker(m Marker) []error {
	var errs []error

	if strings.TrimSpace(m.Tool) == "" {
		errs = append(errs, errors.New("marker field `tool` is required"))
	} else if m.Tool != "gokit-scaffold" {
		errs = append(errs, fmt.Errorf("marker field `tool` must be `gokit-scaffold`, got %q", m.Tool))
	}
	if strings.TrimSpace(m.Version) == "" {
		errs = append(errs, errors.New("marker field `version` is required"))
	}
	if strings.TrimSpace(m.TemplatePack) == "" {
		errs = append(errs, errors.New("marker field `template_pack` is required"))
	} else if m.TemplatePack != TemplatePackName {
		errs = append(errs, fmt.Errorf("marker field `template_pack` must be %q, got %q", TemplatePackName, m.TemplatePack))
	}
	if err := validateName(m.Spec.Name); err != nil {
		errs = append(errs, fmt.Errorf("marker field `spec.name` %v", err))
	}
	if err := validateModule(m.Spec.Module); err != nil {
		errs = append(errs, fmt.Errorf("marker field `spec.module` %v", err))
	}
	if err := validateHTTPPort(m.Spec.HTTPPort); err != nil {
		errs = append(errs, fmt.Errorf("marker field `spec.http_port` %v", err))
	}

	return errs
}

func validateName(name string) error {
	if !nameRe.MatchString(name) {
		return errors.New("must match ^[a-z][a-z0-9-]*$")
	}
	return nil
}

func validateModule(module string) error {
	if !moduleRe.MatchString(module) || strings.Contains(module, "..") {
		return errors.New("is invalid")
	}
	return nil
}

func validateHTTPPort(port int) error {
	if port <= 0 || port > 65535 {
		return errors.New("must be between 1 and 65535")
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
