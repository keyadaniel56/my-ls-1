package ls

import (
	"fmt"
	"os"
	"strings"

	"my-ls-1/internal/models"
)

func List(path string, flags models.Flags) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var entries []models.FileEntry

	for _, f := range files {
		name := f.Name()

		if !flags.All && strings.HasPrefix(name, ".") {
			continue
		}

		info, _ := f.Info()

		entries = append(entries, models.FileEntry{
			Name:  name,
			Path:  path + "/" + name,
			Info:  info,
			IsDir: f.IsDir(),
		})
	}

	entries = Sort(entries, flags)

	if flags.Long {
		PrintLong(entries)
	} else {
		PrintSimple(entries)
	}

	if flags.Recursive {
		fmt.Println()
		Recursive(path, entries, flags)
	}

	return nil
}