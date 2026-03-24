// workspace-test runs go test ./... in every module listed by go list -m (Go workspace).
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	root, err := findRepoRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "workspace-test: %v\n", err)
		os.Exit(1)
	}
	if err := os.Chdir(root); err != nil {
		fmt.Fprintf(os.Stderr, "workspace-test: chdir: %v\n", err)
		os.Exit(1)
	}
	out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "workspace-test: go list -m: %v\n", err)
		os.Exit(1)
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
		fmt.Printf("=== %s ===\n", rel)
		cmd := exec.Command("go", "test", "./...")
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil {
			failed = append(failed, rel)
		}
	}
	if len(failed) > 0 {
		fmt.Fprintf(os.Stderr, "workspace-test: failed modules: %s\n", strings.Join(failed, ", "))
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
