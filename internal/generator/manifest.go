package generator

import (
	"fmt"
	"sort"
	"strings"
)

const ServiceHTTPTemplatePack = "service-http"

type ManifestEntry struct {
	TemplatePath string
	OutputPath   string
}

var serviceHTTPManifest = []ManifestEntry{
	{TemplatePath: "service-http/.gokit-scaffold.tmpl", OutputPath: ".gokit-scaffold"},
	{TemplatePath: "service-http/README.md.tmpl", OutputPath: "README.md"},
	{TemplatePath: "service-http/cmd/server/main.go.tmpl", OutputPath: "cmd/server/main.go"},
	{TemplatePath: "service-http/go.mod.tmpl", OutputPath: "go.mod"},
	{TemplatePath: "service-http/internal/config/config.go.tmpl", OutputPath: "internal/config/config.go"},
	{TemplatePath: "service-http/internal/httpserver/health.go.tmpl", OutputPath: "internal/httpserver/health.go"},
	{TemplatePath: "service-http/internal/httpserver/server.go.tmpl", OutputPath: "internal/httpserver/server.go"},
	{TemplatePath: "service-http/internal/logging/logging.go.tmpl", OutputPath: "internal/logging/logging.go"},
}

func TemplatePacks() []string {
	return []string{ServiceHTTPTemplatePack}
}

func Manifest(templatePack string) ([]ManifestEntry, error) {
	var manifest []ManifestEntry
	switch templatePack {
	case ServiceHTTPTemplatePack:
		manifest = append([]ManifestEntry(nil), serviceHTTPManifest...)
	default:
		return nil, fmt.Errorf("unknown template pack: %s", templatePack)
	}

	sort.Slice(manifest, func(i, j int) bool {
		if manifest[i].OutputPath == manifest[j].OutputPath {
			return manifest[i].TemplatePath < manifest[j].TemplatePath
		}
		return manifest[i].OutputPath < manifest[j].OutputPath
	})

	return manifest, nil
}

func OutputPaths(templatePack string) ([]string, error) {
	manifest, err := Manifest(templatePack)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(manifest))
	for _, entry := range manifest {
		paths = append(paths, entry.OutputPath)
	}
	return paths, nil
}

func OutputTree(templatePack string) (string, error) {
	paths, err := OutputPaths(templatePack)
	if err != nil {
		return "", err
	}
	return BuildASCIITree(paths), nil
}

func BuildASCIITree(paths []string) string {
	type node struct {
		name     string
		children map[string]*node
	}

	root := &node{
		name:     ".",
		children: map[string]*node{},
	}

	for _, p := range paths {
		parts := strings.Split(strings.TrimSpace(p), "/")
		current := root
		for _, part := range parts {
			if part == "" {
				continue
			}
			if current.children[part] == nil {
				current.children[part] = &node{
					name:     part,
					children: map[string]*node{},
				}
			}
			current = current.children[part]
		}
	}

	var b strings.Builder
	b.WriteString(".\n")

	var render func(n *node, prefix string)
	render = func(n *node, prefix string) {
		names := make([]string, 0, len(n.children))
		for name := range n.children {
			names = append(names, name)
		}
		sort.Strings(names)

		for i, name := range names {
			child := n.children[name]
			last := i == len(names)-1

			branch := "|-- "
			nextPrefix := prefix + "|   "
			if last {
				branch = "`-- "
				nextPrefix = prefix + "    "
			}

			b.WriteString(prefix)
			b.WriteString(branch)
			b.WriteString(child.name)
			b.WriteString("\n")
			render(child, nextPrefix)
		}
	}

	render(root, "")
	return strings.TrimRight(b.String(), "\n")
}
