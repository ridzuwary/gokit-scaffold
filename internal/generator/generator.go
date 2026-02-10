package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
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

func Generate(s spec.ProjectSpec, version string) error {
	data := templateData{
		Name:     s.Name,
		Module:   s.Module,
		HTTPPort: s.HTTPPort,
		Version:  version,
	}

	entries, err := Manifest(ServiceHTTPTemplatePack)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if err := renderToFile(s.Dir, entry, data); err != nil {
			return err
		}
	}

	return nil
}

func renderToFile(baseDir string, entry ManifestEntry, data templateData) error {
	body, err := templates.FS.ReadFile(entry.TemplatePath)
	if err != nil {
		return fmt.Errorf("read template %s: %w", entry.TemplatePath, err)
	}

	tpl, err := template.New(entry.TemplatePath).Parse(string(body))
	if err != nil {
		return fmt.Errorf("parse template %s: %w", entry.TemplatePath, err)
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, data); err != nil {
		return fmt.Errorf("render template %s: %w", entry.TemplatePath, err)
	}

	target := filepath.Join(baseDir, filepath.Clean(entry.OutputPath))
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
