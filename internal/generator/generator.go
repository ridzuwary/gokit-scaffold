package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/ridzuwary/gokit-scaffold/internal/spec"
	"github.com/ridzuwary/gokit-scaffold/templates"
)

// cspell:words ridzuwary gokit tmpl

type templateData struct {
	Name     string
	Module   string
	HTTPPort int
	Version  string
}

type renderEntry struct {
	templatePath string
	outputPath   string
}

func Generate(s spec.ProjectSpec, version string) error {
	data := templateData{
		Name:     s.Name,
		Module:   s.Module,
		HTTPPort: s.HTTPPort,
		Version:  version,
	}

	entries := []renderEntry{
		{templatePath: "service-http/.gokit-scaffold.tmpl", outputPath: ".gokit-scaffold"},
		{templatePath: "service-http/README.md.tmpl", outputPath: "README.md"},
		{templatePath: "service-http/cmd/server/main.go.tmpl", outputPath: "cmd/server/main.go"},
		{templatePath: "service-http/go.mod.tmpl", outputPath: "go.mod"},
		{templatePath: "service-http/internal/config/config.go.tmpl", outputPath: "internal/config/config.go"},
		{templatePath: "service-http/internal/httpserver/health.go.tmpl", outputPath: "internal/httpserver/health.go"},
		{templatePath: "service-http/internal/httpserver/server.go.tmpl", outputPath: "internal/httpserver/server.go"},
		{templatePath: "service-http/internal/logging/logging.go.tmpl", outputPath: "internal/logging/logging.go"},
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].outputPath < entries[j].outputPath
	})

	for _, entry := range entries {
		if err := renderToFile(s.Dir, entry, data); err != nil {
			return err
		}
	}

	return nil
}

func renderToFile(baseDir string, entry renderEntry, data templateData) error {
	body, err := templates.FS.ReadFile(entry.templatePath)
	if err != nil {
		return fmt.Errorf("read template %s: %w", entry.templatePath, err)
	}

	tpl, err := template.New(entry.templatePath).Parse(string(body))
	if err != nil {
		return fmt.Errorf("parse template %s: %w", entry.templatePath, err)
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, data); err != nil {
		return fmt.Errorf("render template %s: %w", entry.templatePath, err)
	}

	target := filepath.Join(baseDir, filepath.Clean(entry.outputPath))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create output directory for %s: %w", target, err)
	}

	if err := writeFileAtomic(target, out.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write output file %s: %w", target, err)
	}

	return nil
}

func writeFileAtomic(path string, content []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".gokit-scaffold-tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()

	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if _, err := tmp.Write(content); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Chmod(mode); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	return nil
}
