package ls

import (
	"fmt"
	"os"

	"my-ls-1/internal/models"
)

func Recursive(base string, entries []models.FileEntry, flags models.Flags) {
	for _, e := range entries {
		if e.IsDir && e.Name != "." && e.Name != ".." {
			fmt.Printf("%s:\n", e.Path)

			files, err := os.ReadDir(e.Path)
			if err != nil {
				continue
			}

			var newEntries []models.FileEntry

			for _, f := range files {
				info, _ := f.Info()
				newEntries = append(newEntries, models.FileEntry{
					Name:  f.Name(),
					Path:  e.Path + "/" + f.Name(),
					Info:  info,
					IsDir: f.IsDir(),
				})
			}

			newEntries = Sort(newEntries, flags)

			if flags.Long {
				PrintLong(newEntries)
			} else {
				PrintSimple(newEntries)
			}

			fmt.Println()
			Recursive(e.Path, newEntries, flags)
		}
	}
}
