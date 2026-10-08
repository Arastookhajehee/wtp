package main

import (
	"path/filepath"
	"strings"
)

func canonicalWorktreeName(name string) string {
	// Worktree names are CLI/display identifiers, not filesystem paths. Keep them
	// Git-like on Windows so names such as feature/auth match branch names.
	return filepath.ToSlash(strings.TrimSuffix(name, "*"))
}
