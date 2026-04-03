package ls

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"my-ls-1/internal/models"
)

func List(path string, flags models.Flags) error {
	entries, err := ReadEntries(path, flags)
	if err != nil {
		return err
	}

	entries = Sort(entries, flags)

	if flags.Long {
		printTotal(entries)
		PrintLong(entries)
	} else {
		PrintSimple(entries, flags)
	}

	if flags.Recursive {
		for _, e := range entries {
			if e.IsDir && e.Name != "." && e.Name != ".." {
				fmt.Printf("\n%s:\n", e.Path)
				if err := List(e.Path, flags); err != nil {
					fmt.Fprintf(os.Stderr, "my-ls: cannot access '%s': %v\n", e.Path, err)
				}
			}
		}
	}

	return nil
}

func ReadEntries(path string, flags models.Flags) ([]models.FileEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var entries []models.FileEntry

	if flags.All {
		for _, dot := range []string{".", ".."} {
			// use Lstat so we get the actual entry, not what it points to
			info, err := os.Lstat(filepath.Join(path, dot))
			if err == nil {
				entries = append(entries, models.FileEntry{
					Name:  dot,
					Path:  filepath.Join(path, dot),
					Info:  info,
					IsDir: info.IsDir(),
				})
			}
		}
	}

	for _, f := range files {
		name := f.Name()
		if !flags.All && strings.HasPrefix(name, ".") {
			continue
		}
		// use Lstat so symlinks show as symlinks, not their targets
		info, err := os.Lstat(filepath.Join(path, name))
		if err != nil {
			continue
		}
		entries = append(entries, models.FileEntry{
			Name:  name,
			Path:  filepath.Join(path, name),
			Info:  info,
			IsDir: f.IsDir(),
		})
	}

	return entries, nil
}

func printTotal(entries []models.FileEntry) {
	var total uint64
	for _, e := range entries {
		if e.Info != nil {
			if stat, ok := getStat(e.Info); ok {
				total += uint64(stat.Blocks)
			}
		}
	}
	fmt.Printf("total %d\n", total/2)
}

// FileEntries builds FileEntry list from explicit file paths (non-directory args).
// Uses Lstat so symlinks passed directly are shown as symlinks.
func FileEntries(paths []string, flags models.Flags) []models.FileEntry {
	var entries []models.FileEntry
	for _, p := range paths {
		info, err := os.Lstat(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "my-ls: cannot access '%s': %v\n", p, err)
			continue
		}
		entries = append(entries, models.FileEntry{
			Name:  filepath.Base(p),
			Path:  p,
			Info:  info,
			IsDir: info.IsDir(),
		})
	}
	return Sort(entries, flags)
}
