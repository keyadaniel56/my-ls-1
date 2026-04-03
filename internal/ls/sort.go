package ls

import (
	"sort"
	"strings"
	"unicode"

	"my-ls-1/internal/models"
)

// sortKey mimics ls/strcoll behavior:
// strips leading dots, then lowercases, then strips non-alphanumeric
// for the primary comparison key so punctuation sorts before letters/digits.
func sortKey(name string) string {
	// strip leading dots (ls ignores them for ordering)
	s := strings.TrimLeft(name, ".")
	return strings.ToLower(s)
}

// lsLess replicates the sort order of GNU ls (strcoll-based).
// Primary: compare ignoring non-alphanumeric chars (so '[' < '0' < 'a')
// Secondary: full lowercase comparison
// Tertiary: exact comparison
func lsLess(a, b string) bool {
	ka := sortKey(a)
	kb := sortKey(b)

	// build "alpha-only" keys for primary comparison
	pa := alphaKey(ka)
	pb := alphaKey(kb)

	if pa != pb {
		return pa < pb
	}
	if ka != kb {
		return ka < kb
	}
	return a < b
}

// alphaKey keeps only letters and digits, lowercased.
// Non-alphanumeric chars are replaced with a character that sorts before '0'.
func alphaKey(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			// use a char that sorts before '0' (ASCII 48) — use space (32)
			b.WriteByte(' ')
		}
	}
	return b.String()
}

func Sort(entries []models.FileEntry, flags models.Flags) []models.FileEntry {
	if flags.SortTime {
		sort.SliceStable(entries, func(i, j int) bool {
			ti := entries[i].Info.ModTime()
			tj := entries[j].Info.ModTime()
			if ti.Equal(tj) {
				return lsLess(entries[i].Name, entries[j].Name)
			}
			return ti.After(tj)
		})
	} else {
		sort.SliceStable(entries, func(i, j int) bool {
			return lsLess(entries[i].Name, entries[j].Name)
		})
	}

	if flags.Reverse {
		for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
			entries[i], entries[j] = entries[j], entries[i]
		}
	}

	return entries
}
