package ls

import (
	"sort"

	"my-ls-1/internal/models"
)

func Sort(entries []models.FileEntry, flags models.Flags) []models.FileEntry {
	if flags.SortTime {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Info.ModTime().After(entries[j].Info.ModTime())
		})
	} else {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name < entries[j].Name
		})
	}

	if flags.Reverse {
		for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
			entries[i], entries[j] = entries[j], entries[i]
		}
	}

	return entries
}