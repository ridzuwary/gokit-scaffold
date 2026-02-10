package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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

	ui.PrintInfo(fmt.Sprintf("validate is planned to check %s marker in %s", spec.MarkerFileName, *dir))
	return 0
}

func runPrint(args []string) int {
	fs := flag.NewFlagSet("print", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	ui.PrintInfo("template-pack: service-http")
	ui.PrintInfo("files: go.mod, cmd/server/main.go, README.md, .gokit-scaffold")
	return 0
}
