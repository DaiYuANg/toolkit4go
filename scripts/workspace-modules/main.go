// workspace-modules runs go test, golangci-lint, staticcheck, or govulncheck in each go.work module.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	root, err := findRepoRoot()
	if err != nil {
		fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		fatal(err)
	}

	switch os.Args[1] {
	case "test":
		runPerModule(root, "test", func(_ string) []string {
			return []string{"go", "test", "./..."}
		})
	case "lint":
		bin, err := toolPath(root, "golangci-lint")
		if err != nil {
			fatal(fmt.Errorf("golangci-lint: %w", err))
		}
		runPerModule(root, "lint", func(_ string) []string {
			return []string{bin, "run", "./..."}
		})
	case "staticcheck":
		bin, err := toolPath(root, "staticcheck")
		if err != nil {
			fatal(fmt.Errorf("staticcheck: %w", err))
		}
		checks := strings.TrimSpace(os.Getenv("STATICCHECK_CHECKS"))
		runPerModule(root, "staticcheck", func(_ string) []string {
			if checks != "" {
				return []string{bin, "-checks", checks, "./..."}
			}
			return []string{bin, "./..."}
		})
	case "govulncheck":
		bin, err := toolPath(root, "govulncheck")
		if err != nil {
			fatal(fmt.Errorf("govulncheck: %w", err))
		}
		runPerModule(root, "govulncheck", func(_ string) []string {
			return []string{bin, "./..."}
		})
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: go run ./scripts/workspace-modules/main.go <test|lint|staticcheck|govulncheck>\n")
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "workspace-modules: %v\n", err)
	os.Exit(1)
}

func toolPath(root, tool string) (string, error) {
	cmd := exec.Command("go", "tool", "-n", tool)
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) > 0 {
			return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(string(ee.Stderr)))
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func runPerModule(root, name string, argv func(dir string) []string) {
	out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}").Output()
	if err != nil {
		fatal(fmt.Errorf("go list -m: %w", err))
	}
	var failed []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		dir := strings.TrimSpace(line)
		if dir == "" {
			continue
		}
		rel, err := filepath.Rel(root, dir)
		if err != nil {
			rel = dir
		}
		if rel == "." {
			rel = "."
		}
		fmt.Printf("=== %s (%s) ===\n", rel, name)
		args := argv(dir)
		if len(args) == 0 {
			continue
		}
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil {
			failed = append(failed, rel)
		}
	}
	if len(failed) > 0 {
		fmt.Fprintf(os.Stderr, "workspace-modules %s: failed modules: %s\n", name, strings.Join(failed, ", "))
		os.Exit(1)
	}
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.work not found from %s", dir)
		}
		dir = parent
	}
}
