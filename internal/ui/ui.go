package ui

import (
	"fmt"
	"os"
)

func PrintRootUsage(tool string) {
	fmt.Fprintf(os.Stderr, "Usage: %s <command> [flags]\n", tool)
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  new       Generate a new project scaffold")
	fmt.Fprintln(os.Stderr, "  validate  Validate an existing scaffold (stub)")
	fmt.Fprintln(os.Stderr, "  print     Print embedded template pack details")
}

func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
}

func PrintInfo(message string) {
	fmt.Fprintln(os.Stdout, message)
}

func PrintNewSuccess(dir string) {
	fmt.Fprintf(os.Stdout, "Scaffold generated at %s\n", dir)
}
