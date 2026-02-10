package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ridzuwary/gokit-scaffold/internal/generator"
	"github.com/ridzuwary/gokit-scaffold/internal/spec"
	"github.com/ridzuwary/gokit-scaffold/internal/ui"
)

const (
	ToolName    = "gokit-scaffold"
	ToolVersion = "0.1.0"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		ui.PrintRootUsage(ToolName)
		return 2
	}

	switch args[0] {
	case "new":
		return runNew(args[1:])
	case "validate":
		return runValidate(args[1:])
	case "print":
		return runPrint(args[1:])
	case "-h", "--help", "help":
		ui.PrintRootUsage(ToolName)
		return 0
	default:
		ui.PrintError(fmt.Errorf("unknown command: %s", args[0]))
		ui.PrintRootUsage(ToolName)
		return 2
	}
}

func runNew(args []string) int {
	fs := flag.NewFlagSet("new", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	name := fs.String("name", "", "project name (required)")
	module := fs.String("module", "", "go module path (required)")
	dir := fs.String("dir", "", "output directory (default ./<name>)")
	httpPort := fs.Int("http-port", 8080, "HTTP listen port")

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *name == "" || *module == "" {
		ui.PrintError(fmt.Errorf("--name and --module are required"))
		fs.Usage()
		return 2
	}

	targetDir := *dir
	if targetDir == "" {
		targetDir = filepath.Join(".", *name)
	}

	project := spec.ProjectSpec{
		Name:     *name,
		Module:   *module,
		Dir:      targetDir,
		HTTPPort: *httpPort,
	}
	if err := project.Validate(); err != nil {
		ui.PrintError(err)
		return 1
	}

	if err := generator.Generate(project, ToolVersion); err != nil {
		ui.PrintError(err)
		return 1
	}

	ui.PrintNewSuccess(project.Dir)
	return 0
}

func runValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	dir := fs.String("dir", ".", "directory to validate")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	errs := spec.ValidateScaffoldDir(*dir)
	if len(errs) > 0 {
		ui.PrintError(fmt.Errorf("validation failed for %s", *dir))
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "  - %s\n", strings.TrimSpace(err.Error()))
		}
		return 1
	}

	ui.PrintInfo(fmt.Sprintf("scaffold is valid: %s", *dir))
	return 0
}

func runPrint(args []string) int {
	fs := flag.NewFlagSet("print", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	out, err := formatPrintOutput(ToolName, ToolVersion)
	if err != nil {
		ui.PrintError(err)
		return 1
	}

	fmt.Fprintln(os.Stdout, out)
	return 0
}

func formatPrintOutput(toolName, toolVersion string) (string, error) {
	packs := generator.TemplatePacks()
	tree, err := generator.OutputTree(generator.ServiceHTTPTemplatePack)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s %s\n\n", toolName, toolVersion))
	b.WriteString("Template Packs\n")
	for _, pack := range packs {
		b.WriteString("- ")
		b.WriteString(pack)
		b.WriteString("\n")
	}

	b.WriteString("\nGenerated Outputs Tree (service-http)\n")
	b.WriteString(tree)
	b.WriteString("\n\nMarker Schema Summary\n")
	b.WriteString("- `tool`: scaffold generator identifier (`gokit-scaffold`)\n")
	b.WriteString("- `version`: tool version used to generate the scaffold\n")
	b.WriteString("- `template_pack`: template pack name (`service-http`)\n")
	b.WriteString("- `spec.name`: service name (`^[a-z][a-z0-9-]*$`)\n")
	b.WriteString("- `spec.module`: Go module path\n")
	b.WriteString("- `spec.http_port`: HTTP listen port (1-65535)\n")
	b.WriteString("\nExample new command\n")
	b.WriteString("gokit-scaffold new --name hello-api --module github.com/acme/hello-api --http-port 8080\n")

	return strings.TrimRight(b.String(), "\n"), nil
}
