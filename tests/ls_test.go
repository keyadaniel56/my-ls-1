package tests

import (
	"os"
	"path/filepath"
	"testing"

	"my-ls-1/internal/ls"
	"my-ls-1/internal/models"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	files := []string{"banana.txt", "apple.txt", "cherry.txt", ".hidden"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestListNoFlags(t *testing.T) {
	dir := setupTestDir(t)
	flags := models.Flags{}
	if err := ls.List(dir, flags); err != nil {
		t.Fatalf("List returned error: %v", err)
	}
}

func TestListAll(t *testing.T) {
	dir := setupTestDir(t)
	flags := models.Flags{All: true}
	entries, err := ls.ReadEntries(dir, flags)
	if err != nil {
		t.Fatal(err)
	}
	hasHidden := false
	hasDot := false
	hasDotDot := false
	for _, e := range entries {
		if e.Name == ".hidden" {
			hasHidden = true
		}
		if e.Name == "." {
			hasDot = true
		}
		if e.Name == ".." {
			hasDotDot = true
		}
	}
	if !hasHidden {
		t.Error("expected .hidden to appear with -a flag")
	}
	if !hasDot {
		t.Error("expected . to appear with -a flag")
	}
	if !hasDotDot {
		t.Error("expected .. to appear with -a flag")
	}
}

func TestListHiddenExcluded(t *testing.T) {
	dir := setupTestDir(t)
	flags := models.Flags{}
	entries, err := ls.ReadEntries(dir, flags)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name[0] == '.' {
			t.Errorf("hidden file %q should not appear without -a", e.Name)
		}
	}
}

func TestSortByName(t *testing.T) {
	dir := setupTestDir(t)
	flags := models.Flags{}
	entries, _ := ls.ReadEntries(dir, flags)
	entries = ls.Sort(entries, flags)

	for i := 1; i < len(entries); i++ {
		if entries[i-1].Name > entries[i].Name {
			t.Errorf("entries not sorted: %q > %q", entries[i-1].Name, entries[i].Name)
		}
	}
}

func TestSortReverse(t *testing.T) {
	dir := setupTestDir(t)
	flags := models.Flags{Reverse: true}
	entries, _ := ls.ReadEntries(dir, flags)
	entries = ls.Sort(entries, flags)

	for i := 1; i < len(entries); i++ {
		if entries[i-1].Name < entries[i].Name {
			t.Errorf("entries not reverse sorted: %q < %q", entries[i-1].Name, entries[i].Name)
		}
	}
}

func TestSortByTime(t *testing.T) {
	dir := setupTestDir(t)
	flags := models.Flags{SortTime: true}
	entries, _ := ls.ReadEntries(dir, flags)
	entries = ls.Sort(entries, flags)

	for i := 1; i < len(entries); i++ {
		ti := entries[i-1].Info.ModTime()
		tj := entries[i].Info.ModTime()
		if ti.Before(tj) {
			t.Errorf("entries not sorted by time (newest first): %v < %v", ti, tj)
		}
	}
}

func TestListLong(t *testing.T) {
	dir := setupTestDir(t)
	flags := models.Flags{Long: true}
	if err := ls.List(dir, flags); err != nil {
		t.Fatalf("List -l returned error: %v", err)
	}
}

func TestListRecursive(t *testing.T) {
	dir := setupTestDir(t)
	flags := models.Flags{Recursive: true}
	if err := ls.List(dir, flags); err != nil {
		t.Fatalf("List -R returned error: %v", err)
	}
}
