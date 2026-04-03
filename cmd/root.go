package cmd

import (
	"fmt"
	"os"
	"strings"

	"my-ls-1/internal/ls"
	"my-ls-1/internal/models"
)

func Run(flags models.Flags, paths []string) {
	if len(paths) == 0 {
		paths = []string{"."}
	}

	var files []string
	var dirs []string

	for _, p := range paths {
		// trailing slash forces directory treatment (follow symlink)
		hasTrailingSlash := strings.HasSuffix(p, "/")

		linfo, err := os.Lstat(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "my-ls: cannot access '%s': No such file or directory\n", p)
			continue
		}

		if linfo.Mode()&os.ModeSymlink != 0 && !hasTrailingSlash {
			// bare symlink arg: treat as file (show symlink itself)
			files = append(files, p)
			continue
		}

		// follow for stat (handles trailing slash and symlink-to-dir)
		info, err := os.Stat(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "my-ls: cannot access '%s': %v\n", p, err)
			continue
		}

		if info.IsDir() {
			dirs = append(dirs, p)
		} else {
			files = append(files, p)
		}
	}

	// print non-directory entries first
	if len(files) > 0 {
		entries := ls.FileEntries(files, flags)
		if flags.Long {
			ls.PrintLong(entries)
		} else {
			ls.PrintSimple(entries, flags)
		}
	}

	// then directories
	for i, dir := range dirs {
		if len(files) > 0 || len(dirs) > 1 {
			if i > 0 || len(files) > 0 {
				fmt.Println()
			}
			fmt.Printf("%s:\n", dir)
		}
		if err := ls.List(dir, flags); err != nil {
			fmt.Fprintf(os.Stderr, "my-ls: cannot access '%s': %v\n", dir, err)
		}
	}
}
