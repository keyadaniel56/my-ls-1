package ls

import (
	"fmt"
	"os"
	"strings"

	"my-ls-1/internal/models"
)

// ANSI color codes
const (
	colorReset = "\033[0m"
	colorBlue  = "\033[1;34m" // directory
	colorCyan  = "\033[1;36m" // symlink
	colorGreen = "\033[1;32m" // executable
)

func colorName(name string) string {
	// strip symlink arrow for mode check
	base := name
	if idx := strings.Index(name, " -> "); idx != -1 {
		base = name[:idx]
	}

	info, err := os.Lstat(base)
	if err != nil {
		return name
	}

	mode := info.Mode()
	switch {
	case mode&os.ModeSymlink != 0:
		return colorCyan + name + colorReset
	case info.IsDir():
		return colorBlue + name + colorReset
	case mode&0o111 != 0:
		return colorGreen + name + colorReset
	}
	return name
}

func PrintSimple(entries []models.FileEntry, flags models.Flags) {
	if flags.Long || flags.OnePerLine {
		for _, e := range entries {
			fmt.Println(colorName(e.Name))
		}
		return
	}
	for _, e := range entries {
		fmt.Print(colorName(e.Name) + "  ")
	}
	if len(entries) > 0 {
		fmt.Println()
	}
}
